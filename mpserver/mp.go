package main

import (
	"errors"

	"github.com/aouyang1/go-matrixprofile/matrixprofile"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type MP struct {
	AV         []float64 `json:"annotation_vector"`
	AdjustedMP []float64 `json:"adjusted_mp"`
}

func getMP(c *gin.Context) {
	endpoint := "/api/v1/mp"
	session := sessions.Default(c)
	buildCORSHeaders(c)

	params := struct {
		Name string `json:"name"`
	}{}
	if err := c.BindJSON(&params); err != nil {
		requestTotal.WithLabelValues("POST", endpoint, "500").Inc()
		c.JSON(500, RespError{
			Error: errors.New("failed to unmarshall POST parameters with field `name`"),
		})
		return
	}
	avname := params.Name

	v := session.Get("mp")
	var mp matrixprofile.MatrixProfile
	if v == nil {
		// matrix profile is not initialized so don't return any data back for the
		// annotation vector
		requestTotal.WithLabelValues("POST", endpoint, "500").Inc()
		c.JSON(500, RespError{
			Error:        errors.New("matrix profile is not initialized"),
			CacheExpired: true,
		})
		return
	} else {
		mp = v.(matrixprofile.MatrixProfile)
	}

	switch avname {
	case "default", "":
		mp.AV = matrixprofile.DefaultAV
	case "complexity":
		mp.AV = matrixprofile.ComplexityAV
	case "meanstd":
		mp.AV = matrixprofile.MeanStdAV
	case "clipping":
		mp.AV = matrixprofile.ClippingAV
	default:
		requestTotal.WithLabelValues("POST", endpoint, "500").Inc()
		c.JSON(500, RespError{
			Error: errors.New("invalid annotation vector name " + avname),
		})
		return
	}

	// cache matrix profile for current session
	session.Set("mp", &mp)
	session.Save()

	av, err := mp.GetAV()
	if err != nil {
		requestTotal.WithLabelValues("POST", endpoint, "500").Inc()
		c.JSON(500, RespError{Error: err})
		return
	}

	adjustedMP, err := mp.ApplyAV(av)
	if err != nil {
		requestTotal.WithLabelValues("POST", endpoint, "500").Inc()
		c.JSON(500, RespError{Error: err})
		return
	}

	requestTotal.WithLabelValues("POST", endpoint, "200").Inc()
	c.JSON(200, MP{AV: av, AdjustedMP: adjustedMP})
}
