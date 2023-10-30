package monitoring

import "github.com/gofiber/fiber/v2"

type API struct {
}

func (a API) MountRouters(app fiber.Router) {
	r := fiber.New()

	app.Mount("/monitoring", r)
}
