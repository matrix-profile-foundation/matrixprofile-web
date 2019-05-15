package main

import (
	"errors"
	"strconv"

	"github.com/aouyang1/go-matrixprofile/matrixprofile"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Motif struct {
	Groups []matrixprofile.MotifGroup `json:"groups"`
	Series [][][]float64              `json:"series"`
}

func topKMotifs(c *gin.Context) {
	endpoint := "/api/v1/topkmotifs"
	method := "GET"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	k, err := strconv.Atoi(c.Query("k"))
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, RespError{Error: err})
		return
	}

	r, err := strconv.ParseFloat(c.Query("r"), 64)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, RespError{Error: err})
		return
	}

	v := session.Get("mp")

	var mp matrixprofile.MatrixProfile
	if v == nil {
		// either the cache expired or this was called directly
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, RespError{
			Error:        errors.New("matrix profile is not initialized to compute motifs"),
			CacheExpired: true,
		})
		return
	} else {
		mp = v.(matrixprofile.MatrixProfile)
	}
	motifGroups, err := mp.TopKMotifs(k, r)
	if err != nil {
		requestTotal.WithLabelValues(method, endpoint, "500").Inc()
		c.JSON(500, RespError{Error: err})
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
				c.JSON(500, RespError{Error: err})
				return
			}
		}
	}

	requestTotal.WithLabelValues(method, endpoint, "200").Inc()
	c.JSON(200, motif)
}
