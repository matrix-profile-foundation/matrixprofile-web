package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/aouyang1/go-matrixprofile/matrixprofile"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type Data struct {
	Data []float64 `json:"data"`
}

type Motif struct {
	Groups []matrixprofile.MotifGroup `json:"groups"`
	Series [][][]float64              `json:"series"`
}

type Discord struct {
	Groups []int       `json:"groups"`
	Series [][]float64 `json:"series"`
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

func main() {
	r := gin.Default()

	var redis_url string
	if redis_url = os.Getenv("REDIS_URL"); redis_url == "" {
		redis_url = "localhost:6379"
	}

	store, err := redis.NewStore(10, "tcp", redis_url, "", []byte("secret"))
	if err != nil {
		panic(err)
	}

	err, rs := redis.GetRedisStore(store)
	if err != nil {
		panic(err)
	}
	rs.SetMaxLength(1024 * 1024)

	r.Use(sessions.Sessions("mysession", store))
	r.Use(cors.Default())

	gob.RegisterName("github.com/aouyang1/go-matrixprofile/matrixprofile.MatrixProfile", matrixprofile.MatrixProfile{})

	r.GET("/data", func(c *gin.Context) {
		data, err := fetchData()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.JSON(200, data.Data)
	})

	r.GET("/calculate", func(c *gin.Context) {
		session := sessions.Default(c)
		c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

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
		if err = mp.Stomp(2); err != nil {
			c.JSON(500, gin.H{"error": err})
			return
		}

		session.Set("mp", &mp)
		session.Set("user", "aouyang")
		session.Save()

		c.JSON(200, mp.MP)
	})

	r.GET("/topkmotifs", func(c *gin.Context) {
		session := sessions.Default(c)
		c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

		kstr := c.Query("k")
		rstr := c.Query("r")

		k, err := strconv.Atoi(kstr)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{"error": err})
			return
		}

		r, err := strconv.ParseFloat(rstr, 64)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{"error": err})
			return
		}

		v := session.Get("mp")
		fmt.Println(session.Get("user").(string))

		var mp matrixprofile.MatrixProfile
		if v == nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": "matrix profile was never initialized to compute motifs",
			})
			return
		} else {
			mp = v.(matrixprofile.MatrixProfile)
		}
		out, _ := json.MarshalIndent(mp, "", "  ")
		fmt.Println(string(out))
		fmt.Printf("%d, %d\n", len(mp.A), len(mp.B))
		motifGroups, err := mp.TopKMotifs(k, r)
		if err != nil {
			fmt.Println(err)
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
					fmt.Println(err)
					c.JSON(500, gin.H{"error": err})
					return
				}
			}
		}

		c.JSON(200, motif)
	})

	r.GET("/topkdiscords", func(c *gin.Context) {
		session := sessions.Default(c)
		c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

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
				"error": "matrix profile was never initialized to compute discords",
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
	})

	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8081"
	}
	r.Run(":" + port)
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
