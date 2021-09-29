// Copyright (c) 2021 Amirhossein Movahedi (@qolzam)
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package router

import (
	"github.com/GMcD/telar-web/micros/admin/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/middleware/authcookie"
	"github.com/red-gold/telar-core/middleware/authrole"
	"github.com/red-gold/telar-core/utils"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {

	// Middleware
	authCookieMiddleware := authcookie.New(authcookie.Config{
		JWTSecretKey: []byte(*config.AppConfig.PublicKey),
		Authorizer:   utils.VerifyJWT,
	})
	authRoleMiddleware := authrole.New(authrole.ConfigDefault)

	// Router
	app.Post("/setup", authCookieMiddleware, authRoleMiddleware, handlers.SetupHandler)
	app.Get("/setup", authCookieMiddleware, authRoleMiddleware, handlers.SetupPageHandler)
	app.Get("/login", handlers.LoginPageHandler)
	app.Post("/login", handlers.LoginAdminHandler)
}
