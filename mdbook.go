package mdbook

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ilaziness/mdbook/internal/booklib"
	"github.com/ilaziness/mdbook/internal/util"
)

var BookServer = &MdBook{
	Engine: gin.Default(),
}

type MdBook struct {
	Config *Config
	Engine *gin.Engine
}

func (m *MdBook) Run() error {
	m.bootstrap()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	util.Info("server listen on %s", m.Config.GetPort())
	gin.SetMode(m.Config.GetMode())
	return m.Engine.Run(m.Config.GetPort())
}

func (m *MdBook) bootstrap() {
	initConfig()
	initRoute()
	err := booklib.ScanBook(AppConfig.BookDir)
	if err != nil {
		log.Println(err)
	}
	m.Config = AppConfig
}
