package http

import (
	"biatosh/contract"
	"biatosh/http/router"
	"embed"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/sqlite3"
	"github.com/gofiber/template/html/v2"
)

//go:embed templates/*
var templatesFS embed.FS

func New(store contract.Store, log contract.Logger) *fiber.App {
	engine := html.NewFileSystem(http.FS(templatesFS), ".html")

	storage := sqlite3.New()
	sessStore := newSessionStore(storage)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Use(csrf.New(csrf.Config{
		Storage: storage,
		Session: sessStore,

		KeyLookup:         "form:csrf_",
		CookieName:        "csrf_",
		CookieSameSite:    "Lax",
		Expiration:        1 * time.Hour,
		KeyGenerator:      utils.UUIDv4,
		Extractor:         csrf.CsrfFromForm("csrf_"),
		SessionKey:        "fiber.csrf.token",
		HandlerContextKey: "fiber.csrf.handler",

		CookieSecure:   true,
		CookieHTTPOnly: true,
	}))

	router.Setup(store, log, app, sessStore)
	return app
}

func newSessionStore(storage *sqlite3.Storage) *session.Store {
	store := session.New(session.Config{
		Storage: storage,

		CookieSessionOnly: true,
		CookieHTTPOnly:    true,
		CookieSecure:      true,
	})
	return store
}
