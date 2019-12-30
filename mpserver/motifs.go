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

type Motif struct {
	Groups []mp.MotifGroup `json:"groups"`
	Series [][][]float64   `json:"series"`
}

func topKMotifs(c *gin.Context) {
	start := time.Now()
	endpoint := "/api/v1/topkmotifs"
	method := "GET"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	k, err := strconv.Atoi(c.Query("k"))
	if err != nil {
		logError(RespError{Error: err}, method, endpoint, start, c)
		return
	}

	r, err := strconv.ParseFloat(c.Query("r"), 64)
	if err != nil {
		logError(RespError{Error: err}, method, endpoint, start, c)
		return
	}

	v := fetchMPCache(session)

	var p mp.MatrixProfile
	if v == nil {
		// either the cache expired or this was called directly
		respError := RespError{
			Error:        errors.New("matrix profile is not initialized to compute motifs"),
			CacheExpired: true,
		}
		logError(respError, method, endpoint, start, c)
		return
	} else {
		p = v.(mp.MatrixProfile)
	}
	motifGroups, err := p.DiscoverMotifs(k, r)
	if err != nil {
		logError(RespError{Error: err}, method, endpoint, start, c)
		return
	}

	var motif Motif
	motif.Groups = motifGroups
	motif.Series = make([][][]float64, len(motifGroups))
	for i, g := range motif.Groups {
		motif.Series[i] = make([][]float64, len(g.Idx))
		for j, midx := range g.Idx {
			motif.Series[i][j], err = util.ZNormalize(p.A[midx : midx+p.M])
			if err != nil {
				logError(RespError{Error: err}, method, endpoint, start, c)
				return
			}
		}
	}

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	serviceRequestDuration.WithLabelValues(endpoint).Observe(time.Since(start).Seconds() * 1000)
	c.JSON(200, motif)
}
