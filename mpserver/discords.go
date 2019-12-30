package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	mp "github.com/matrix-profile-foundation/go-matrixprofile"
	"github.com/matrix-profile-foundation/go-matrixprofile/util"
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
		logError(RespError{Error: err}, method, endpoint, start, c)
		return
	}

	v := fetchMPCache(session)
	var p mp.MatrixProfile
	if v == nil {
		respError := RespError{
			errors.New("matrix profile is not initialized to compute discords"),
			true,
		}
		logError(respError, method, endpoint, start, c)
		return
	} else {
		p = v.(mp.MatrixProfile)
	}
	discords, err := p.DiscoverDiscords(k, p.M/2)
	if err != nil {
		respError := RespError{
			Error: errors.New("failed to compute discords"),
		}
		logError(respError, method, endpoint, start, c)
		return
	}

	var discord Discord
	discord.Groups = discords
	discord.Series = make([][]float64, len(discords))
	for i, didx := range discord.Groups {
		discord.Series[i], err = util.ZNormalize(p.A[didx : didx+p.M])
		if err != nil {
			logError(RespError{Error: err}, method, endpoint, start, c)
			return
		}
	}

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
	c.JSON(200, discord)
}
