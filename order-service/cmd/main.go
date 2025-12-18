package main

import (
	"os"

	"github.com/4udiwe/big-bob-pizza/order-service/internal/app"
)

// @title Big Bob Pizza - Order Service API
// @version 1.0
// @description API для управления заказами в системе Big Bob Pizza
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@bigbobpizza.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	app := app.New(os.Getenv("CONFIG_PATH"))
	app.Start()
}
