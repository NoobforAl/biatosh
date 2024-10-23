package router

import (
	"biatosh/contract"
	"biatosh/http/controller"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func Setup(store contract.Store, log contract.Logger, app *fiber.App, sessStore *session.Store) {
	controller := controller.New(store, log, sessStore)

	app.Get("/", controller.Dashboard())

	app.Get("/login", controller.LoginPage())
	app.Post("/login", controller.LoginUser())

	app.Get("/logout", controller.Logout())

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			sess, err := sessStore.Get(c)
			if err != nil {
				return err
			}

			c.Locals("session", sess)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", controller.WS())
}
