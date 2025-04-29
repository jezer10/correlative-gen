package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/emidiaz3/event-driven-server/config"
	"github.com/emidiaz3/event-driven-server/controllers"
	"github.com/emidiaz3/event-driven-server/database"
	"github.com/emidiaz3/event-driven-server/models"
	"github.com/emidiaz3/event-driven-server/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"

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

var store *session.Store

func StartFiberServer() {
	// Configuración de Asynq y Redis
	config.SetupDatabaseConnection()
	store = session.New()

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

	// Inicializa el motor de plantillas. Aquí se buscan los templates en "./views" con extensión ".html"
	engine := html.New("./views", ".html")

	// Crea una única instancia de Fiber configurada con el motor de plantillas
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Ruta para renderizar la página de inicio
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Página de Inicio",
			"Msg":   "Bienvenido a mi aplicación en Go con Fiber XDD!",
		})
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		users, err := GetAllUsers()
		if err != nil {
			database.SaveErrorLog(err.Error())
			return c.Status(500).SendString("Error al obtener usuarios")
		}

		return c.Render("users", fiber.Map{
			"Title": "Lista de Usuarios",
			"Users": users,
		})
	})

	// Rutas de API
	app.Post("/user/ids", apiKeyMiddleware, controllers.GetUsers(config.DB))

	app.Post("/", apiKeyMiddleware, controllers.CreateUser(config.DB, redisClient))

	app.Post("/login", func(c *fiber.Ctx) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		db := database.GetDB()

		var user models.Users // 👈 asegurate de usar `models.Users` con "s"

		// Buscar usuario por username
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Credenciales inválidas")
		}

		// Verificar la contraseña con bcrypt
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Credenciales inválidas")
		}

		// Crear sesión
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		sess.Set("user", user.Username)
		if err := sess.Save(); err != nil {
			return err
		}

		return c.Redirect("/dashboard")
	})

	app.Get("/dashboard", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		// Verifica que la sesión tenga el dato de usuario
		user := sess.Get("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Debes iniciar sesión")
		}
		postulantes, err := GetAllPostulantes()
		if err != nil {
			database.SaveErrorLog(err.Error())
			return c.Status(500).SendString("Error al obtener usuarios")
		}
		log.Println(postulantes)
		return c.Render("dashboard", fiber.Map{
			"Title":       "Página de Inicio",
			"Postulantes": postulantes,
		})

	})
	app.Post("/api/users", CreateUser)

	app.Get("/postulantes/edit/:id", EditPostulante)
	app.Post("/postulantes/update", UpdatePostulante)

	// Ruta para el panel de monitoreo
	app.All("/monitor/*", adaptor.HTTPHandler(h))

	log.Println("🚀 Servidor en ejecución")
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		database.SaveErrorLog(err.Error())
		log.Fatal("❌ Error al iniciar el servidor:", err)
	}
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := database.GetDB().Find(&users).Error; err != nil {
		return nil, err
	}
	log.Print(users)
	return users, nil
}

func GetAllPostulantes() ([]models.Postulantes, error) {
	var postulantes []models.Postulantes
	if err := database.GetDB().Find(&postulantes).Error; err != nil {
		return nil, err
	}
	log.Print(postulantes)
	return postulantes, nil
}
func EditPostulante(c *fiber.Ctx) error {
	// Obtiene el ID desde la URL
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ID inválido")
	}

	// Busca el postulante en la base de datos
	var postulante models.Postulantes
	if err := database.GetDB().First(&postulante, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Postulante no encontrado")
	}

	// Renderiza la vista de edición
	return c.Render("edit_postulante", fiber.Map{
		"Title":      "Editar Postulante",
		"Postulante": postulante,
	})
}
func UpdatePostulante(c *fiber.Ctx) error {
	idParam := c.FormValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("ID inválido")
	}

	var postulante models.Postulantes
	if err := database.GetDB().First(&postulante, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Postulante no encontrado")
	}

	postulante.FirstName = c.FormValue("first_name")
	postulante.LastName = c.FormValue("last_name")
	postulante.Identity = c.FormValue("identity")
	postulante.Birthday = c.FormValue("birthday")
	postulante.NativeCountry = c.FormValue("native_country")
	postulante.Country = c.FormValue("country")
	postulante.Correlative = c.FormValue("correlative")
	// --- NUEVOS CAMPOS ---
	// Status (sql.NullString)
	statusVal := c.FormValue("status")
	postulante.Status = sql.NullString{
		String: statusVal,
		Valid:  statusVal != "",
	}
	// Score (pointer a string)
	if s := c.FormValue("score"); s != "" {
		postulante.Score = &s
	} else {
		postulante.Score = nil
	}
	// Score Description
	if sd := c.FormValue("score_description"); sd != "" {
		postulante.ScoreDescription = &sd
	} else {
		postulante.ScoreDescription = nil
	}
	// Score Note
	if sn := c.FormValue("score_note"); sn != "" {
		postulante.ScoreNote = &sn
	} else {
		postulante.ScoreNote = nil
	}
	// --- FIN NUEVOS CAMPOS ---

	if err := database.GetDB().Save(&postulante).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error al guardar")
	}

	return c.Redirect("/dashboard")
}

type UserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateUser(c *fiber.Ctx) error {
	var input UserInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "JSON inválido",
		})
	}

	if input.Username == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username y Password son obligatorios",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No se pudo encriptar la contraseña",
		})
	}

	user := models.Users{
		Id:        uuid.New(),
		Username:  input.Username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al guardar usuario",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":       user.Id,
		"username": user.Username,
		"created":  user.CreatedAt,
	})
}
