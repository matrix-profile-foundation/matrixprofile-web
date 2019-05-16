package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

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
	start := time.Now()
	endpoint := "/api/v1/data"
	method := "GET"
	data, err := fetchData()
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		c.JSON(500, RespError{Error: err})
		return
	}

	c.Header("Content-Type", "application/json")
	buildCORSHeaders(c)

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
	c.JSON(200, data.Data)
}
