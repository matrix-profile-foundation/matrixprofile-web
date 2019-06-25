package main

import (
	"time"

	"github.com/aouyang1/go-matrixprofile/matrixprofile"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

type Segment struct {
	CAC []float64 `json:"cac"`
}

func calculateMP(c *gin.Context) {
	var err error
	var mp *matrixprofile.MatrixProfile

	start := time.Now()
	endpoint := "/api/v1/calculate"
	method := "POST"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	params := struct {
		M       int    `json:"m"`
		SourceA string `json:"sourceA"`
		SourceB string `json:"sourceB"`
	}{}
	if err := c.BindJSON(&params); err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		glog.Infof("%v", err)
		c.JSON(500, RespError{Error: err.Error()})
		return
	}
	m := params.M
	sourceA := params.SourceA
	sourceB := params.SourceB

	dataA, err := fetchData(sourceA)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		glog.Infof("%v", err)
		c.JSON(500, RespError{Error: err.Error()})
		return
	}

	if sourceA != sourceB {
		dataB, err := fetchData(sourceB)
		if err != nil {
			requestTotal.WithLabelValues(method, endpoint, "500").Inc()
			serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
			glog.Infof("%v", err)
			c.JSON(500, RespError{Error: err.Error()})
			return
		}
		mp, err = matrixprofile.New(dataA.Data, dataB.Data, m)
	} else {
		// self-join case
		mp, err = matrixprofile.New(dataA.Data, nil, m)
	}
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		glog.Infof("%v", err)
		c.JSON(500, RespError{Error: err.Error()})
		return
	}

	if err = mp.Stomp(mpConcurrency); err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		glog.Infof("%v", err)
		c.JSON(500, RespError{Error: err.Error()})
		return
	}

	// compute the corrected arc curve based on the current index matrix profile
	_, _, cac := mp.Segment()

	// cache matrix profile for current session
	storeMPCache(session, mp)

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
	c.JSON(200, Segment{cac})
}
