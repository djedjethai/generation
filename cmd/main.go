package main

import (
	"github.com/djedjethai/generation/internal/agent"
	"log"
	"os"
)

func main() {

	cfg, err := setupSrv()
	if err != nil {
		os.Exit(1)
	}

	_, err = agent.New(cfg)
	if err != nil {
		log.Println("the err from setting the agent: ", err)
	}
}
