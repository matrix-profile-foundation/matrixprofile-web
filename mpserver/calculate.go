package main

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	mp "github.com/matrix-profile-foundation/go-matrixprofile"
)

type Segment struct {
	CAC []float64 `json:"cac"`
}

func calculateMP(c *gin.Context) {
	start := time.Now()
	endpoint := "/api/v1/calculate"
	method := "POST"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	params := struct {
		M      int    `json:"m"`
		Source string `json:"source"`
	}{}
	if err := c.BindJSON(&params); err != nil {
		logError(RespError{Error: err}, method, endpoint, start, c)
		return
	}
	m := params.M
	source := params.Source

	data, err := fetchData(source)
	if err != nil {
		logError(RespError{Error: err}, method, endpoint, start, c)
		return
	}

	p, err := mp.New(data.Data, nil, m)
	if err != nil {
		logError(RespError{Error: err}, method, endpoint, start, c)
		return
	}

	if err = p.Compute(mp.NewComputeOpts()); err != nil {
		logError(RespError{Error: err}, method, endpoint, start, c)
		return
	}

	// compute the corrected arc curve based on the current index matrix profile
	_, _, cac := p.DiscoverSegments()

	// cache matrix profile for current session
	storeMPCache(session, p)

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
	c.JSON(200, Segment{cac})
}
