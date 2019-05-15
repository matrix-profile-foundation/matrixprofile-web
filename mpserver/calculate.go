package main

import (
	"github.com/aouyang1/go-matrixprofile/matrixprofile"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Segment struct {
	CAC []float64 `json:"cac"`
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
		c.JSON(500, RespError{Error: err})
		return
	}
	m := params.M

	data, err := fetchData()
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, RespError{Error: err})
		return
	}

	mp, err := matrixprofile.New(data.Data, nil, m)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, RespError{Error: err})
		return
	}

	if err = mp.Stomp(mpConcurrency); err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, RespError{Error: err})
		return
	}

	// compute the corrected arc curve based on the current index matrix profile
	_, _, cac := mp.Segment()

	// cache matrix profile for current session
	session.Set("mp", &mp)
	session.Save()

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	c.JSON(200, Segment{cac})
}
