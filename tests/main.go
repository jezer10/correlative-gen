package tests

import (
	"encoding/json"

	"github.com/emidiaz3/event-driven-server/database"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Person struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Identity      string `json:"identity"`
	Birthday      string `json:"birthday"`
	NativeCountry string `json:"native_country"`
	Country       string `json:"country"`
}

func InitClient() {
	app := fiber.New()
	err := database.InitDBMSSQL()
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "TuContraseñaSegura123!",
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/", func(c *fiber.Ctx) error {
		var person Person

		if err := c.BodyParser(&person); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No se pudo parsear el body",
			})
		}
		data, err := json.Marshal(person)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al serializar a JSON",
			})
		}
		key := "person: " + person.Identity

		if err := rdb.Set(ctx, key, data, 0).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No se pudo guardar en Redis",
			})
		}

		return c.JSON(fiber.Map{
			"status": "OK",
			"key":    key,
			"data":   person,
		})

	})

	app.Listen(":3000")
}
