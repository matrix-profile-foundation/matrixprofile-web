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
)

var (
	mpConcurrency    = 2
	maxRedisBlobSize = 1024 * 1024
	retentionPeriod  = 5 * 60
	redisURL         = "localhost:6379"
)

func main() {
	r := gin.Default()

	store, err := initRedis()
	if err != nil {
		panic(err)
	}

	r.Use(sessions.Sessions("mysession", store))
	r.Use(cors.Default())

	gob.RegisterName("github.com/aouyang1/go-matrixprofile/matrixprofile.MatrixProfile", matrixprofile.MatrixProfile{})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/data", getData)
		v1.GET("/calculate", calculateMP)
		v1.GET("/topkmotifs", topKMotifs)
		v1.GET("/topkdiscords", topKDiscords)
	}

	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}

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
	data, err := fetchData()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")
	buildCORSHeaders(c)
	c.JSON(200, data.Data)
}

func calculateMP(c *gin.Context) {
	session := sessions.Default(c)
	buildCORSHeaders(c)

	mstr := c.Query("m")
	m, err := strconv.Atoi(mstr)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	data, err := fetchData()
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	mp, err := matrixprofile.New(data.Data, nil, m)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	if err = mp.Stomp(mpConcurrency); err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	session.Set("mp", &mp)
	session.Save()

	c.JSON(200, mp.MP)
}

type Motif struct {
	Groups []matrixprofile.MotifGroup `json:"groups"`
	Series [][][]float64              `json:"series"`
}

func topKMotifs(c *gin.Context) {
	session := sessions.Default(c)
	buildCORSHeaders(c)

	kstr := c.Query("k")
	rstr := c.Query("r")

	k, err := strconv.Atoi(kstr)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	r, err := strconv.ParseFloat(rstr, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	v := session.Get("mp")

	var mp matrixprofile.MatrixProfile
	if v == nil {
		// either the cache expired or this was called directly
		c.JSON(500, gin.H{
			"error": "matrix profile is not initialized to compute motifs",
		})
		return
	} else {
		mp = v.(matrixprofile.MatrixProfile)
	}
	motifGroups, err := mp.TopKMotifs(k, r)
	if err != nil {
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
				c.JSON(500, gin.H{"error": err})
				return
			}
		}
	}

	c.JSON(200, motif)
}

type Discord struct {
	Groups []int       `json:"groups"`
	Series [][]float64 `json:"series"`
}

func topKDiscords(c *gin.Context) {
	session := sessions.Default(c)
	buildCORSHeaders(c)

	kstr := c.Query("k")

	k, err := strconv.Atoi(kstr)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	v := session.Get("mp")
	var mp matrixprofile.MatrixProfile
	if v == nil {
		c.JSON(500, gin.H{
			"error": "matrix profile is not initialized to compute discords",
		})
		return
	} else {
		mp = v.(matrixprofile.MatrixProfile)
	}
	discords := mp.TopKDiscords(k, mp.M/2)

	var discord Discord
	discord.Groups = discords
	discord.Series = make([][]float64, len(discords))
	for i, didx := range discord.Groups {
		discord.Series[i], err = matrixprofile.ZNormalize(mp.A[didx : didx+mp.M])
		if err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}
	}

	c.JSON(200, discord)
}

func buildCORSHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}
