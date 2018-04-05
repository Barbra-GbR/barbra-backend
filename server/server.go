package server

import (
	"../config"
)

func Init() {
	c := config.GetConfig()
	r := NewRouter()

	r.Run(c.GetString("server.port"))
}
