package server

import (
	"github.com/Barbra-GbR/barbra-backend/config"
)

//Initializes the server
func Initialize() {
	c := config.GetConfig()
	r := NewRouter()
	r.Run(c.GetString("server.port"))
}
