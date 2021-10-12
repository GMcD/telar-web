package function

import (
	"context"
	"net/http"

	micros "github.com/GMcD/telar-web/micros"
	authConfig "github.com/GMcD/telar-web/micros/auth/config"
	"github.com/GMcD/telar-web/micros/auth/database"
	"github.com/GMcD/telar-web/micros/auth/router"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
)

// Cache state
var app *fiber.App

// init
func init() {

	// Init config
	micros.InitConfig()
	authConfig.InitConfig()

	// Initialize app
	app = fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(
		logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path} - ${header:}\n",
		},
	))
	log.Info("CORS origin: %s", *config.AppConfig.Origin)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     *config.AppConfig.Origin,
		AllowCredentials: true,
		AllowHeaders:     "Authorization, uid, email, avatar, displayName, role, tagLine, x-cloud-signature, Origin, Content-Type, Accept, Access-Control-Allow-Headers, X-Requested-With, X-HTTP-Method-Override, access-control-allow-origin, access-control-allow-headers",
	}))
	router.SetupRoutes(app)
}

// Handler function
func Handle(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	// Connect
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
