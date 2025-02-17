package server

import (
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
	fmt.Print(requestKey)

	if requestKey != utils.ApiKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "API Key inválida",
		})
	}

	return c.Next()
}

func getUserByID(userID int) (*models.User, string, error) {
	user, err := database.GetUserByIdMysql(userID)
	fmt.Print(err)

	if err == nil {
		return user, "Main", nil
	}

	user, err = database.GetUserById(userID)
	if err == nil {
		return user, "Replica", nil
	}

	return nil, "", err
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

	app.Get("/user/:id", apiKeyMiddleware, func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		userID, err := strconv.Atoi(idParam)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
		}

		user, source, err := getUserByID(userID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Usuario no encontrado"})
		}

		return c.JSON(fiber.Map{"source": source, "user": user})
	})

	app.Post("/", apiKeyMiddleware, func(c *fiber.Ctx) error {
		var user models.User

		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No se pudo parsear el body",
			})
		}

		correlative, err := database.InsertUser(user)
		if err != nil {
			return fmt.Errorf("error insertando usuario en la BD: %w", err)
		}
		user.Correlative = correlative

		if err := tasks.SendUserToQueue(redisClient, user); err != nil {
			log.Printf("Error encolando usuario: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No se pudo encolar el usuario"})
		}

		return c.JSON(fiber.Map{
			"status": "OK",
			"data":   user,
		})

	})
	app.All("/monitor/*", adaptor.HTTPHandler(h))

	log.Println("🚀 Servidor en ejecución")
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatal("❌ Error al iniciar el servidor:", err)
	}

}
