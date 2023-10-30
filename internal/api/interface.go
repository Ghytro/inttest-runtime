package api

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	MountRouters(app fiber.Router)
}
