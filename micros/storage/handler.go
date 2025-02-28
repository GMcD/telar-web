package function

import (
	"net/http"

	micros "github.com/GMcD/telar-web/micros"
	appConfig "github.com/GMcD/telar-web/micros/storage/config"
	"github.com/GMcD/telar-web/micros/storage/router"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/red-gold/telar-core/config"
)

// Cache state
var app *fiber.App

// Extra Headers
var helmetHeaders = helmet.Config{
	ContentSecurityPolicy: "upgrade-insecure-requests; script-src 'self' *.prod.monitalks.io; img-src 'self' https://prod-monitalks-media.s3.eu-west-2.amazonaws.com;",
	CSPReportOnly:         false,
	HSTSPreloadEnabled:    true,
	ReferrerPolicy:        "origin",
	HSTSMaxAge:            31536000,
	HSTSExcludeSubdomains: true,
}

func init() {
	appConfig.InitConfig()
	micros.InitConfig()

	// Initialize app
	app = fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // this is the default limit of 10MB
	})
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(helmet.New(helmetHeaders))
	app.Use(logger.New(
		logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path} - ${header:}\n​",
		},
	))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     *config.AppConfig.Origin,
		AllowCredentials: true,
		AllowHeaders:     "Authorization, uid, email, avatar, displayName, role, tagLine, x-cloud-signature, Origin, Content-Type, Accept, Access-Control-Allow-Headers, X-Requested-With, X-HTTP-Method-Override, access-control-allow-origin, access-control-allow-headers",
	}))
	router.SetupRoutes(app)
}

// Handler function
func Handle(w http.ResponseWriter, r *http.Request) {

	adaptor.FiberApp(app)(w, r)
}
