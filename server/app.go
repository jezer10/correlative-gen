package server

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/emidiaz3/event-driven-server/database"
	"github.com/emidiaz3/event-driven-server/models"
	"github.com/emidiaz3/event-driven-server/tasks"
	"github.com/emidiaz3/event-driven-server/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

func apiKeyMiddleware(c *fiber.Ctx) error {
	requestKey := c.Get("X-API-Key")
	if requestKey != utils.ApiKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "API Key inválida",
		})
	}

	return c.Next()
}

func StartFiberServer() {
	redisClientOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	redisClient := asynq.NewClient(redisClientOpt)
	defer redisClient.Close()

	h := asynqmon.New(asynqmon.Options{
		RootPath:     "/monitor",
		RedisConnOpt: redisClientOpt,
	})

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/user/ids", apiKeyMiddleware, func(c *fiber.Ctx) error {
		var body models.RequestBody
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Formato JSON inválido"})
		}

		var checkIDs []string
		for _, raw := range body.CheckIDs {
			var strVal string
			if err := json.Unmarshal(raw, &strVal); err == nil {
				checkIDs = append(checkIDs, fmt.Sprintf("'%s'", strVal))
				continue
			}

			var intVal int
			if err := json.Unmarshal(raw, &intVal); err == nil {
				tmpConv := strconv.Itoa(intVal)
				checkIDs = append(checkIDs, fmt.Sprintf("'%s'", tmpConv))
				continue
			}

			return c.Status(400).JSON(fiber.Map{
				"check_ids": []string{},
				"err":       "Valores inválidos en check_ids"})
		}
		usrs, err := database.GetUsers(checkIDs)
		if err != nil {
			database.SaveErrorLog(err.Error())
			return c.Status(400).JSON(fiber.Map{
				"check_ids": []string{},
				"err":       "Error al consultar los usuarios",
			})
		}

		return c.JSON(fiber.Map{
			"message":   "Consulta realizada con éxito",
			"check_ids": usrs,
		})
	})

	app.Post("/", apiKeyMiddleware, func(c *fiber.Ctx) error {
		var user models.User

		if err := c.BodyParser(&user); err != nil {
			database.SaveErrorLog(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"check_id": "0",
				"err":      "Formato Erróneo de envío",
			})
		}

		insertedId, err := database.InsertUser(user)

		if err != nil {
			database.SaveErrorLog(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"check_id": "0",
				"err":      "No se insertó el usuario",
			})
		}

		if err := tasks.SendUserToQueue(redisClient, user); err != nil {
			database.SaveErrorLog(err.Error())
			if err := database.DeleteUser(insertedId); err != nil {
				database.SaveErrorLog(err.Error())
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"check_id": "0",
				"error":    "No se pudo encolar el usuario",
			})
		}
		newUser, err := database.GetUserById(insertedId)

		if err != nil {
			database.SaveErrorLog(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"check_id": "0",
				"err":      "No se insertó el usuario",
			})
		}

		return c.JSON(fiber.Map{
			"check_id": newUser.Correlative,
		})

	})
	app.All("/monitor/*", adaptor.HTTPHandler(h))

	log.Println("🚀 Servidor en ejecución")
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		database.SaveErrorLog(err.Error())
		log.Fatal("❌ Error al iniciar el servidor:", err)
	}

}
