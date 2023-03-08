package mdbook

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

var AppConfig *Config

func initConfig() {
	AppConfig = &Config{}
}

type Config struct {
	// 运行模式
	Mode    string
	Port    int
	BookDir string
}

func (c *Config) GetMode() string {
	switch c.Mode {
	case gin.DebugMode:
		return gin.DebugMode
	case gin.ReleaseMode:
		return gin.ReleaseMode
	default:
		return gin.TestMode
	}
}

func (c *Config) GetPort() string {
	if c.Port != 0 && c.Port > 1000 && c.Port <= 65534 {
		return ":" + strconv.Itoa(c.Port)
	}

	switch c.Mode {
	case gin.ReleaseMode:
		return ":80"
	default:
		return ":8080"
	}
}
