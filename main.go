package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/puper/go-mysql-schema-sync/internal"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	configFile = flag.String("config", "config.json", "config file")
	port       = flag.String("port", "8081", "port")
)

func main() {
	flag.Parse()
	r := gin.Default()
	r.GET("/", index)
	r.Run(":" + *port)
}

type Config struct {
	AccessKey string            `json:accessKey`
	Databases map[string]string `json:databases`
}

func loadConfig(name string) *Config {
	var c Config
	data, _ := ioutil.ReadFile(name)
	json.Unmarshal(data, &c)
	return &c
}

func index(c *gin.Context) {
	config := loadConfig(*configFile)
	accessKey := c.Query("accessKey")
	if accessKey != config.AccessKey {
		c.String(http.StatusOK, "wrong accessKey")
		return
	}
	source := c.Query("source")
	target := c.Query("target")
	if source == target {
		c.String(http.StatusOK, "source can not equal target")
		return
	}
	if _, ok := config.Databases[source]; !ok {
		c.String(http.StatusOK, "source not found")
		return
	}
	if _, ok := config.Databases[target]; !ok {
		c.String(http.StatusOK, "target not found")
		return
	}
	db := c.Query("db")
	if db == "" {
		c.String(http.StatusOK, "db can not be empty")
		return
	}
	var result []string
	result = append(result, "source: "+source)
	result = append(result, "target: "+target)
	result = append(result, "db: "+db)
	result = append(result, "---- sql ----")
	result = append(result, internal.GetSql(config.Databases[source]+"/"+db, config.Databases[target]+"/"+db)...)
	c.String(http.StatusOK, strings.Join(result, "\n"))
}
