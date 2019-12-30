package main

import (
	"errors"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	mp "github.com/matrix-profile-foundation/go-matrixprofile"
	"github.com/matrix-profile-foundation/go-matrixprofile/av"
)

type MP struct {
	AV         []float64 `json:"annotation_vector"`
	AdjustedMP []float64 `json:"adjusted_mp"`
}

func getMP(c *gin.Context) {
	start := time.Now()
	endpoint := "/api/v1/mp"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	params := struct {
		Name string `json:"name"`
	}{}
	if err := c.BindJSON(&params); err != nil {
		respError := RespError{
			Error: errors.New("failed to unmarshall POST parameters with field `name`"),
		}
		logError(respError, "POST", endpoint, start, c)
		return
	}
	avname := params.Name

	v := fetchMPCache(session)
	var p mp.MatrixProfile
	if v == nil {
		// matrix profile is not initialized so don't return any data back for the
		// annotation vector
		respErr := RespError{
			Error:        errors.New("matrix profile is not initialized"),
			CacheExpired: true,
		}
		logError(respErr, "POST", endpoint, start, c)
		return
	} else {
		p = v.(mp.MatrixProfile)
	}

	switch avname {
	case "default", "":
		p.AV = av.Default
	case "complexity":
		p.AV = av.Complexity
	case "meanstd":
		p.AV = av.MeanStd
	case "clipping":
		p.AV = av.Clipping
	default:
		respErr := RespError{
			Error: errors.New("invalid annotation vector name " + avname),
		}
		logError(respErr, "POST", endpoint, start, c)
		return
	}

	// cache matrix profile for current session
	storeMPCache(session, &p)

	adjustedMP, err := p.ApplyAV()
	if err != nil {
		logError(RespError{Error: err}, "POST", endpoint, start, c)
		return
	}

	avec, err := av.Create(p.AV, p.A, p.M)
	if err != nil {
		logError(RespError{Error: err}, "POST", endpoint, start, c)
		return
	}
	requestTotal.WithLabelValues("POST", endpoint, "200").Inc()
	serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
	c.JSON(200, MP{AV: avec, AdjustedMP: adjustedMP})
}
