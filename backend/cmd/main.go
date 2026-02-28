package main

import (
	"github.com/Jayyk09/CUHackIt/config"
	"github.com/Jayyk09/CUHackIt/internal/app"
)

func main() {
	cfg := config.GetConfig()
	app.Run(cfg)
}
