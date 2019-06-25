package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/aouyang1/go-matrixprofile/matrixprofile"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

type Motif struct {
	Groups []matrixprofile.MotifGroup `json:"groups"`
	Series [][][]float64              `json:"series"`
}

func topKMotifs(c *gin.Context) {
	start := time.Now()
	endpoint := "/api/v1/topkmotifs"
	method := "GET"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	k, err := strconv.Atoi(c.Query("k"))
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		glog.Infof("%v", err)
		c.JSON(500, RespError{Error: err.Error()})
		return
	}

	r, err := strconv.ParseFloat(c.Query("r"), 64)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		glog.Infof("%v", err)
		c.JSON(500, RespError{Error: err.Error()})
		return
	}

	v := fetchMPCache(session)

	var mp matrixprofile.MatrixProfile
	if v == nil {
		// either the cache expired or this was called directly
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		err = errors.New("matrix profile is not initialized to compute motifs")
		glog.Infof("%v", err)
		c.JSON(500, RespError{
			Error:        err.Error(),
			CacheExpired: true,
		})
		return
	} else {
		mp = v.(matrixprofile.MatrixProfile)
	}
	motifGroups, err := mp.TopKMotifs(k, r)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
		glog.Infof("%v", err)
		c.JSON(500, RespError{Error: err.Error()})
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
				requestTotal.WithLabelValues(method, endpoint, "500").Inc()
				serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
				glog.Infof("%v", err)
				c.JSON(500, RespError{Error: err.Error()})
				return
			}
		}
	}

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
	c.JSON(200, motif)
}
