package server

import (
	"github.com/bitphinix/barbra_backend/config"
)

func Init() {
	c := config.GetConfig()
	r := NewRouter()

	r.Run(c.GetString("server.port"))
}
