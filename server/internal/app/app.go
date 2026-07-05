package app

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/zain-23/local-vault/server/internal/common/apperror"
	"github.com/zain-23/local-vault/server/internal/config"
)

const serverVersion = "3.0.0"

// App holds the Fiber app and database — the single object main.go runs
type App struct {
	Fiber *fiber.App
	DB    *mongo.Database
}

// New connects to MongoDB, configures Fiber with middleware, and registers health endpoint
func New(cfg config.Config) (*App, error) {
	// --- MongoDB ---
	// ApplyURI parses connection string and configures connection pool, timeouts, etc.
	clientOpts := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel() // defer runs when function returns — frees timeout resources
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	log.Printf("✅ Connected to MongoDB: %s", cfg.MongoDB)

	// Select database — MongoDB creates it automatically on first write
	db := client.Database(cfg.MongoDB)

	// --- Fiber ---
	app := fiber.New(fiber.Config{
		// ErrorHandler runs when any handler returns an error — centralizes error responses
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Check if error is our custom apperror type — if yes, use its status code
			if appErr, ok := errors.AsType[*apperror.Error](err); ok {
				return c.Status(appErr.Status).JSON(fiber.Map{"error": appErr.Message})
			}
			// Unknown error — send 500, don't leak internal details
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		},
	})

	// CORS — allows frontend (different port/domain) to call API, without this browser blocks requests
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: true,
	}))

	// Logger — logs every request: method, path, status, duration
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${method} ${path} → ${status} (${latency})\n",
	}))

	// Health endpoint — no auth needed, used by monitoring/Docker to check if server is alive
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": serverVersion,
			"time":    time.Now(),
		})
	})

	return &App{Fiber: app, DB: db}, nil
}