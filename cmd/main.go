package main

import (
	"walletapp/config"
	"walletapp/internal/app"
)

func main() {
	app.Run(config.MustGet())
}
