package function

import (
	"context"
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
	micros "github.com/red-gold/ts-serverless/micros"
	"github.com/red-gold/ts-serverless/micros/user-rels/database"
	"github.com/red-gold/ts-serverless/micros/user-rels/router"
)

// Cache state
var app *fiber.App

func init() {

	micros.InitConfig()

	// Initialize app
	app = fiber.New()
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(
		logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path} - ${header:}\n​",
		},
	))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     *config.AppConfig.Origin,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Headers, X-Requested-With, X-HTTP-Method-Override, access-control-allow-origin, access-control-allow-headers",
	}))
	router.SetupRoutes(app)
}

// Handler function
func Handle(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	if database.Db == nil {
		var startErr error
		startErr = database.Connect(ctx)
		if startErr != nil {
			log.Error("Error startup: %s", startErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(startErr.Error()))
		}
	}

	adaptor.FiberApp(app)(w, r)

}
