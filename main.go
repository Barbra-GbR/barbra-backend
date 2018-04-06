package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/bitphinix/babra_backend/db"
	"github.com/bitphinix/babra_backend/server"
	"github.com/bitphinix/babra_backend/config"
	"github.com/bitphinix/babra_backend/auth"
)

func main() {
	environment := flag.String("e", "development", "")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()

	//Initialization
	config.Init(*environment)
	db.Init()
	auth.InitJWT()
	auth.InitUserManager()

	//Initialize and start server
	server.Init()
}