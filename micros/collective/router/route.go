// Copyright (c) 2021 Amirhossein Movahedi (@qolzam)
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package router

import (
	"github.com/GMcD/telar-web/micros/collective/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/middleware/authcookie"
	"github.com/red-gold/telar-core/middleware/authhmac"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {

	// Middleware
	authHMACMiddleware := func(hmacWithCookie bool) func(*fiber.Ctx) error {
		var Next func(c *fiber.Ctx) bool = nil
		if hmacWithCookie {
			Next = func(c *fiber.Ctx) bool {
				if c.Get(types.HeaderHMACAuthenticate) != "" {
					return false
				}
				return true
			}
		}
		return authhmac.New(authhmac.Config{
			Next:          Next,
			PayloadSecret: *config.AppConfig.PayloadSecret,
		})
	}

	authCookieMiddleware := func(hmacWithCookie bool) func(*fiber.Ctx) error {
		var Next func(c *fiber.Ctx) bool = nil
		if hmacWithCookie {
			Next = func(c *fiber.Ctx) bool {
				if c.Get(types.HeaderHMACAuthenticate) != "" {
					return true
				}
				return false
			}
		}
		return authcookie.New(authcookie.Config{
			Next:         Next,
			JWTSecretKey: []byte(*config.AppConfig.PublicKey),
			Authorizer:   utils.VerifyJWT,
		})
	}

	hmacCookieHandlers := []func(*fiber.Ctx) error{authHMACMiddleware(true), authCookieMiddleware(true)}

	// Routers
	app.Get("/id/:collectiveId", append(hmacCookieHandlers, handlers.ReadCollectiveHandle)...)
	app.Get("/", append(hmacCookieHandlers, handlers.QueryCollectiveHandle)...)
	app.Post("/index", authHMACMiddleware(false), handlers.InitCollectiveIndexHandle)

	// // Invoke between functions and protected by HMAC
	app.Get("/dto/id/:collectiveId", authHMACMiddleware(false), handlers.ReadDtoCollectiveHandle)
	app.Post("/dto", authHMACMiddleware(false), handlers.CreateDtoCollectiveHandle)
	app.Post("/dispatch", authHMACMiddleware(false), handlers.DispatchCollectiveHandle)
	app.Post("/dto/ids", authHMACMiddleware(false), handlers.GetCollectiveByIds)
	app.Put("/follow/inc/:inc/:collectiveId", authHMACMiddleware(false), handlers.IncreasePostCount)
	app.Put("/follower/inc/:inc/:collectiveId", authHMACMiddleware(false), handlers.IncreaseFollowerCount)

}
