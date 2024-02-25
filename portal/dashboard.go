package portal

import "github.com/gofiber/fiber/v2"

func dashboard(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}
