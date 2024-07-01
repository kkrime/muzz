package main

import (
	"fmt"
	"muzz/internal/config"
	"muzz/internal/server"
)

func run() error {

	config, err := config.ReadConfig()
	if err != nil {
		return err
	}

	server, err := server.NewServer(config)
	if err != nil {
		return err
	}

	return server.Run()
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("err = %+v\n", err)
	}
}
