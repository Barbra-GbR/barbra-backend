package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/bitphinix/barbra_backend/db"
	"github.com/bitphinix/barbra_backend/server"
	"github.com/bitphinix/barbra_backend/config"
	"github.com/bitphinix/barbra_backend/auth"
	"github.com/bitphinix/barbra_backend/helpers"
)

func main() {
	environment := flag.String("e", "development", "")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()

	//Initialization
	helpers.InitValidator()
	config.Init(*environment)
	db.Init()
	auth.InitJWT()
	auth.InitUserManager()

	//Initialize and start server
	server.Init()
}
