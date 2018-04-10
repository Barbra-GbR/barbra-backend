package server

import (
	"github.com/Barbra-GbR/barbra-backend/config"
)

func Init() {
	c := config.GetConfig()
	r := NewRouter()

	r.Run(c.GetString("server.port"))
}
