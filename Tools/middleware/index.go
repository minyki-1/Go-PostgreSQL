package middleware

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/helmet/v2"
)

func Setting(app *fiber.App) {
	app.Use(helmet.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(etag.New(etag.Config{
		Weak: true,
	}))
	if os.Getenv("ENV_MODE") == "production" {
		app.Use(csrf.New())
		app.Use(limiter.New(limiter.Config{
			Max:               20,
			Expiration:        30 * time.Second,
			LimiterMiddleware: limiter.SlidingWindow{},
		}))
		app.Use(cache.New(cache.Config{
			Next: func(c *fiber.Ctx) bool {
				return c.Query("refresh") == "true"
			},
			Expiration:   3 * time.Minute,
			CacheControl: true,
		}))
		app.Use(cors.New(cors.Config{
			AllowOrigins: "https://html-visualize.vercel.app",
			AllowHeaders: "Origin, Content-Type, Accept",
			AllowMethods: "POST,PUT,DELETE",
		}))
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET",
			AllowHeaders: "Origin, Content-Type, Accept",
		}))
	} else {
		app.Use(cors.New())
	}
}
