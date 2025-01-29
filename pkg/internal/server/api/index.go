package api

import (
	"git.solsynth.dev/hypernet/insight/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func MapAPIs(app *fiber.App, baseURL string) {
	api := app.Group(baseURL).Name("API")
	{
		api.Get("/status", func(c *fiber.Ctx) error {
			err := services.PingOllama()
			if err != nil {
				return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
			}
			return c.SendStatus(fiber.StatusOK)
		})
	}
}
