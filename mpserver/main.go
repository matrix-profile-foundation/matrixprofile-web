package main

import (
	"encoding/gob"
	"os"
	"time"

	"github.com/aouyang1/go-matrixprofile/matrixprofile"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mpConcurrency    = 2
	maxRedisBlobSize = 10 * 1024 * 1024
	retentionPeriod  = 5 * 60
	redisURL         = "localhost:6379" // override with REDIS_URL environment variable
	port             = "8081"           // override with PORT environment variable

	requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mpserver_requests_total",
			Help: "count of all HTTP requests for the mpserver",
		},
		[]string{"method", "endpoint", "code"},
	)
	serviceRequestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "mpserver_service_request_durations_ms",
			Help:       "service request duration in milliseconds.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"endpoint"},
	)
	redisClientRequestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "mpserver_client_redis_request_durations_ms",
			Help:       "redis client request duration in milliseconds.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"command", "status"},
	)
)

type RespError struct {
	Error        error `json:"error"`
	CacheExpired bool  `json:"cache_expired"`
}

func init() {
	prometheus.MustRegister(requestTotal)
	prometheus.MustRegister(serviceRequestDuration)
	prometheus.MustRegister(redisClientRequestDuration)
}

func main() {
	r := gin.Default()

	store, err := initRedis()
	if err != nil {
		panic(err)
	}

	r.Use(sessions.Sessions("mysession", store))
	r.Use(cors.Default())

	gob.RegisterName(
		"github.com/aouyang1/go-matrixprofile/matrixprofile.MatrixProfile",
		matrixprofile.MatrixProfile{},
	)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/data", getData)
		v1.POST("/calculate", calculateMP)
		v1.GET("/topkmotifs", topKMotifs)
		v1.GET("/topkdiscords", topKDiscords)
		v1.POST("/mp", getMP)
	}
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	r.Run(":" + port)
}

// initRedis initializes the connection to the redis store for caching session Matrix Profile data
func initRedis() (redis.Store, error) {
	if u := os.Getenv("REDIS_URL"); u != "" {
		// override global variable if environment variable present
		redisURL = u
	}

	store, err := redis.NewStore(10, "tcp", redisURL, "", []byte("secret"))
	if err != nil {
		return nil, err
	}

	err, rs := redis.GetRedisStore(store)
	if err != nil {
		return nil, err
	}
	rs.SetMaxLength(maxRedisBlobSize)
	rs.Options.MaxAge = retentionPeriod

	return store, nil
}

func buildCORSHeaders(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	c.Header("Access-Control-Allow-Methods", "GET, POST")
}

func fetchMPCache(session sessions.Session) interface{} {
	start := time.Now()

	v := session.Get("mp")
	if v == nil {
		redisClientRequestDuration.WithLabelValues("GET", "500").Observe(time.Since(start).Seconds() * 1000)
	} else {
		redisClientRequestDuration.WithLabelValues("GET", "200").Observe(time.Since(start).Seconds() * 1000)
	}
	return v
}

func storeMPCache(session sessions.Session, mp *matrixprofile.MatrixProfile) {
	start := time.Now()

	session.Set("mp", mp)
	err := session.Save()

	if err != nil {
		redisClientRequestDuration.WithLabelValues("SET", "500").Observe(time.Since(start).Seconds() * 1000)
	} else {
		redisClientRequestDuration.WithLabelValues("SET", "200").Observe(time.Since(start).Seconds() * 1000)
	}
}
