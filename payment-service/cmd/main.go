package main

import (
	"os"

	"github.com/4udiwe/big-bob-pizza/payment-service/internal/app"
)

// @title Big Bob Pizza - Payment Service API
// @version 1.0
// @description API для обработки платежей в системе Big Bob Pizza
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@bigbobpizza.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
// @schemes http https
func main() {
	app := app.New(os.Getenv("CONFIG_PATH"))
	app.Start()
}
