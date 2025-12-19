package main

import (
	"os"

	"github.com/4udiwe/big-bob-pizza/gateway/internal/app"
)

func main() {
	application := app.New(os.Getenv("CONFIG_PATH"))
	application.Start()
}
