package xsession_test

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	xsession "xsession/pkg"
	xhttp "xsession/pkg/http"
)

func Test_server(t *testing.T) {
	s := xhttp.Default()
	s.Engine.GET("/test/:name", func(c *gin.Context) {
		session, _ := xhttp.GetSession(c)
		if session.Get("k1") == nil {
			session.Set("k1", "v1")
			data := make(map[string]interface{})
			data["k2"] = "v2"
			data["k3"] = "v3"
			session.Sets(data)
		}

		name := c.Param("name")
		c.JSON(200, gin.H{
			name: session.Get(name),
		})
	})
	s.Run(":8000")

	t.Run("Defaultserver", func(t *testing.T) {
		resp, err := http.NewRequest("GET", "http://localhost:8000/test/v1", nil)
		assert.NoError(t, err)
		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, string(body), "")
	})
}

func Test_RedisServer(t *testing.T) {
	redis := redis.NewClient(&redis.Options{
		Addr:     "43.138.13.24:6379",
		Password: "",
		DB:       0,
	})
	s := xhttp.Default()
	s.SetSessionStorage(xsession.NewStorageRedis(redis))
	s.Engine.GET("/test/:name", func(c *gin.Context) {
		session, _ := xhttp.GetSession(c)
		if session.Get("k1") == nil {
			session.Set("k1", "v1")
			data := make(map[string]interface{})
			data["k2"] = "v2"
			data["k3"] = "v3"
			session.Sets(data)
		}

		name := c.Param("name")
		c.JSON(200, gin.H{
			name: session.Get(name),
		})
	})
	s.Run(":8000")

	t.Run("Defaultserver", func(t *testing.T) {
		resp, err := http.NewRequest("GET", "http://localhost:8000/test/v1", nil)
		assert.NoError(t, err)
		body, _ := ioutil.ReadAll(resp.Body)
		assert.Equal(t, string(body), "")
	})
}
