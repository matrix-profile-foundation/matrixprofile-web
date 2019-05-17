package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aouyang1/go-matrixprofile/matrixprofile"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Discord struct {
	Groups []int       `json:"groups"`
	Series [][]float64 `json:"series"`
}

func topKDiscords(c *gin.Context) {
	start := time.Now()
	endpoint := "/api/v1/topkdiscords"
	method := "GET"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	kstr := c.Query("k")

	k, err := strconv.Atoi(kstr)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		c.JSON(500, RespError{Error: err})
		return
	}

	v := fetchMPCache(session)
	var mp matrixprofile.MatrixProfile
	if v == nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		c.JSON(500, RespError{
			errors.New("matrix profile is not initialized to compute discords"),
			true,
		})
		return
	} else {
		mp = v.(matrixprofile.MatrixProfile)
	}
	discords, err := mp.TopKDiscords(k, mp.M/2)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		c.JSON(500, RespError{
			Error: errors.New("failed to compute discords"),
		})
		return
	}

	var discord Discord
	discord.Groups = discords
	discord.Series = make([][]float64, len(discords))
	for i, didx := range discord.Groups {
		fmt.Println(didx, len(mp.A))
		discord.Series[i], err = matrixprofile.ZNormalize(mp.A[didx : didx+mp.M])
		if err != nil {
			requestTotal.WithLabelValues(method, endpoint, "500").Inc()
			serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
			c.JSON(500, RespError{Error: err})
			return
		}
	}

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
	c.JSON(200, discord)
}
