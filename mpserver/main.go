package main

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/aouyang1/go-matrixprofile/matrixprofile"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mpConcurrency    = 2
	maxRedisBlobSize = 1024 * 1024
	retentionPeriod  = 5 * 60
	redisURL         = "localhost:6379" // override with REDIS_URL environment variable
	port             = "8081"           // override with PORT environment variable

	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mpserver_requests_total",
			Help: "count of all HTTP requests for the mpserver",
		},
		[]string{"method", "endpoint", "code"},
	)
)

func init() {
	prometheus.MustRegister(requestTotal)
}

func main() {
	r := gin.Default()

	store, err := initRedis()
	if err != nil {
		panic(err)
	}

	r.Use(sessions.Sessions("mysession", store))
	r.Use(cors.Default())

	gob.RegisterName(
		"github.com/aouyang1/go-matrixprofile/matrixprofile.MatrixProfile",
		matrixprofile.MatrixProfile{},
	)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/data", getData)
		v1.POST("/calculate", calculateMP)
		v1.GET("/topkmotifs", topKMotifs)
		v1.GET("/topkdiscords", topKDiscords)
		v1.GET("/segment", segment)
		v1.POST("/anvector", setAnnotationVector)
	}
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	r.Run(":" + port)
}

// initRedis initializes the connection to the redis store for caching session Matrix Profile data
func initRedis() (redis.Store, error) {
	if u := os.Getenv("REDIS_URL"); u != "" {
		// override global variable if environment variable present
		redisURL = u
	}

	store, err := redis.NewStore(10, "tcp", redisURL, "", []byte("secret"))
	if err != nil {
		return nil, err
	}

	err, rs := redis.GetRedisStore(store)
	if err != nil {
		return nil, err
	}
	rs.SetMaxLength(maxRedisBlobSize)
	rs.Options.MaxAge = retentionPeriod

	return store, nil
}

type Data struct {
	Data []float64 `json:"data"`
}

func fetchData() (Data, error) {
	jsonFile, err := os.Open("./penguin_data.json")
	if err != nil {
		return Data{}, err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Data{}, err
	}

	var data Data
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return Data{}, err
	}

	data.Data = smooth(data.Data, 21)[:24*60*7]

	return data, nil
}

// smooth performs a non causal averaging of neighboring data points
func smooth(data []float64, m int) []float64 {
	leftSpan := m / 2
	rightSpan := m / 2

	var sum float64
	var s, e int
	sdata := make([]float64, len(data))

	for i := range data {
		s = i - leftSpan
		if s < 0 {
			s = 0
		}

		e = i + rightSpan + 1
		if e > len(data) {
			e = len(data)
		}

		sum = 0
		for _, d := range data[s:e] {
			sum += d
		}

		sdata[i] = sum / float64(e-s)
	}
	return sdata
}

func getData(c *gin.Context) {
	endpoint := "/api/v1/data"
	method := "GET"
	data, err := fetchData()
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")
	buildCORSHeaders(c)

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	c.JSON(200, data.Data)
}

func calculateMP(c *gin.Context) {
	endpoint := "/api/v1/calculate"
	method := "POST"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	params := struct {
		M int `json:"m"`
	}{}
	if err := c.BindJSON(&params); err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": err})
		return
	}
	m := params.M

	data, err := fetchData()
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": err})
		return
	}

	mp, err := matrixprofile.New(data.Data, nil, m)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": err})
		return
	}
	if err = mp.Stomp(mpConcurrency); err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": err})
		return
	}

	// cache matrix profile for current session
	session.Set("mp", &mp)
	session.Save()

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	c.JSON(200, gin.H{})
}

type Motif struct {
	Groups []matrixprofile.MotifGroup `json:"groups"`
	Series [][][]float64              `json:"series"`
}

func topKMotifs(c *gin.Context) {
	endpoint := "/api/v1/topkmotifs"
	method := "GET"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	k, err := strconv.Atoi(c.Query("k"))
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": err})
		return
	}

	r, err := strconv.ParseFloat(c.Query("r"), 64)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": err})
		return
	}

	v := session.Get("mp")

	var mp matrixprofile.MatrixProfile
	if v == nil {
		// either the cache expired or this was called directly
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{
			"error": "matrix profile is not initialized to compute motifs",
		})
		return
	} else {
		mp = v.(matrixprofile.MatrixProfile)
	}
	motifGroups, err := mp.TopKMotifs(k, r)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": err})
		return
	}

	var motif Motif
	motif.Groups = motifGroups
	motif.Series = make([][][]float64, len(motifGroups))
	for i, g := range motif.Groups {
		motif.Series[i] = make([][]float64, len(g.Idx))
		for j, midx := range g.Idx {
			motif.Series[i][j], err = matrixprofile.ZNormalize(mp.A[midx : midx+mp.M])
			if err != nil {
				requestTotal.WithLabelValues(method, endpoint, "500").Inc()
				c.JSON(500, gin.H{"error": err})
				return
			}
		}
	}

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	c.JSON(200, motif)
}

type Discord struct {
	Groups []int       `json:"groups"`
	Series [][]float64 `json:"series"`
}

func topKDiscords(c *gin.Context) {
	endpoint := "/api/v1/topkdiscords"
	method := "GET"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	kstr := c.Query("k")

	k, err := strconv.Atoi(kstr)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": err})
		return
	}

	v := session.Get("mp")
	var mp matrixprofile.MatrixProfile
	if v == nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{
			"error": "matrix profile is not initialized to compute discords",
		})
		return
	} else {
		mp = v.(matrixprofile.MatrixProfile)
	}
	discords, err := mp.TopKDiscords(k, mp.M/2)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{
			"error": "failed to compute discords",
		})
		return
	}

	var discord Discord
	discord.Groups = discords
	discord.Series = make([][]float64, len(discords))
	for i, didx := range discord.Groups {
		discord.Series[i], err = matrixprofile.ZNormalize(mp.A[didx : didx+mp.M])
		if err != nil {
			requestTotal.WithLabelValues(method, endpoint, "500").Inc()
			c.JSON(500, gin.H{"error": err})
			return
		}
	}

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	c.JSON(200, discord)
}

type Segment struct {
	CAC []float64 `json:"cac"`
}

func segment(c *gin.Context) {
	endpoint := "/api/v1/segment"
	method := "GET"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	v := session.Get("mp")
	var mp matrixprofile.MatrixProfile
	if v == nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, gin.H{
			"error": "matrix profile is not initialized to compute discords",
		})
		return
	} else {
		mp = v.(matrixprofile.MatrixProfile)
	}
	_, _, cac := mp.Segment()

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	c.JSON(200, Segment{cac})
}

type AnnotationVector struct {
	Values []float64 `json:"values"`
	NewMP  []float64 `json:"newmp"`
}

func setAnnotationVector(c *gin.Context) {
	endpoint := "/api/v1/anvector"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	params := struct {
		Name string `json:"name"`
	}{}
	if err := c.BindJSON(&params); err != nil {
		requestTotal.WithLabelValues("POST", endpoint, "500").Inc()
		c.JSON(500, AnnotationVector{})
		return
	}
	avname := params.Name

	v := session.Get("mp")
	var mp matrixprofile.MatrixProfile
	if v == nil {
		// matrix profile is not initialized so don't return any data back for the
		// annotation vector
		requestTotal.WithLabelValues("POST", endpoint, "200").Inc()
		c.JSON(200, AnnotationVector{})
		return
	} else {
		mp = v.(matrixprofile.MatrixProfile)
	}

	switch avname {
	case "default":
		mp.AV = matrixprofile.DefaultAV
	case "complexity":
		mp.AV = matrixprofile.ComplexityAV
	case "meanstd":
		mp.AV = matrixprofile.MeanStdAV
	case "clipping":
		mp.AV = matrixprofile.ClippingAV
	default:
		requestTotal.WithLabelValues("POST", endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": "invalid annotation vector name " + avname})
		return
	}

	// cache matrix profile for current session
	session.Set("mp", &mp)
	session.Save()

	av, err := mp.GetAV()
	if err != nil {
		requestTotal.WithLabelValues("POST", endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": err})
		return
	}

	newMP, err := mp.ApplyAV(av)
	if err != nil {
		requestTotal.WithLabelValues("POST", endpoint, "500").Inc()
		c.JSON(500, gin.H{"error": err})
		return
	}

	requestTotal.WithLabelValues("POST", endpoint, "200").Inc()
	c.JSON(200, AnnotationVector{Values: av, NewMP: newMP})
}

func buildCORSHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	c.Header("Access-Control-Allow-Methods", "GET, POST")
}
