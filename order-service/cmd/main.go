package main

import (
	"os"

	"github.com/4udiwe/big-bob-pizza/order-service/app"
)

func main() {
	app := app.New(os.Getenv("CONFIG_PATH"))
	app.Start()
}
