package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/aouyang1/go-matrixprofile/matrixprofile"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Data struct {
	Data []float64 `json:"data"`
}

type Motif struct {
	Groups []matrixprofile.MotifGroup
	Series [][][]float64 `json:"series"`
}

func main() {
	jsonFile, err := os.Open("./testdata.json")
	if err != nil {
		panic(err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var data Data
	var mp *matrixprofile.MatrixProfile
	if err := json.Unmarshal(byteValue, &data); err != nil {
		panic(err)
	}
	data.Data = data.Data[:50000]

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/data", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(200, data.Data)
	})

	r.GET("/calculate", func(c *gin.Context) {
		mstr := c.Query("m")
		m, err := strconv.Atoi(mstr)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
		}

		mp, err = matrixprofile.New(data.Data, nil, m)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
		}

		if err = mp.Stomp(2); err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
		}

		c.JSON(200, mp.MP)
	})

	r.GET("/topkmotifs", func(c *gin.Context) {
		mstr := c.Query("m")
		kstr := c.Query("k")
		rstr := c.Query("r")

		m, err := strconv.Atoi(mstr)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
		}

		k, err := strconv.Atoi(kstr)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
		}

		r, err := strconv.ParseFloat(rstr, 64)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
		}

		motifGroups, err := mp.TopKMotifs(k, r)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
		}

		var motif Motif
		motif.Groups = motifGroups
		motif.Series = make([][][]float64, len(motifGroups))
		for i, g := range motif.Groups {
			motif.Series[i] = make([][]float64, len(g.Idx))
			for j, midx := range g.Idx {
				motif.Series[i][j], err = matrixprofile.ZNormalize(data.Data[midx : midx+m])
				if err != nil {
					c.JSON(500, gin.H{
						"error": err,
					})
				}
			}
		}
		c.JSON(200, motif)
	})

	r.Run(":8081")
}
