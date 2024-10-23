package controller

import (
	"biatosh/contract"
	"biatosh/entity"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type FuncHttp func(c *fiber.Ctx) error

type controller struct {
	wsStore   *wsStore
	sessStore *session.Store

	dbStore contract.Store
	log     contract.Logger
}

func New(store contract.Store, log contract.Logger, sessStore *session.Store) *controller {
	wsStore := newWsStore(log)

	return &controller{
		wsStore:   wsStore,
		sessStore: sessStore,
		dbStore:   store,
		log:       log,
	}
}

func (c *controller) getCsrfToken(ctx *fiber.Ctx) (string, error) {
	sess, err := c.sessStore.Get(ctx)
	if err != nil {
		return "", err
	}

	csrfToken, ok := sess.Get("fiber.csrf.token").(csrf.Token)
	if !ok {
		return "", err
	}

	return csrfToken.Key, nil
}

func (c *controller) LoginPage() FuncHttp {
	return func(ctx *fiber.Ctx) error {
		csrfToken, err := c.getCsrfToken(ctx)
		if err != nil {
			c.log.Error("Login page get error csrf token:", err)
			return ctx.Render("templates/login", fiber.Map{
				"CSRFToken": "",
				"error":     "Internal server error",
			})
		}

		return ctx.Render("templates/login", fiber.Map{
			"CSRFToken": csrfToken,
		})
	}
}

func (c *controller) LoginUser() FuncHttp {
	return func(ctx *fiber.Ctx) error {
		username := ctx.FormValue("username")
		password := ctx.FormValue("password")

		user, err := c.dbStore.LoginUser(ctx.Context(), &entity.User{
			Username: username,
			Password: password,
		})

		if err != nil {
			c.log.Error("Login user error (not found in db):", err)
			csrfToken, err := c.getCsrfToken(ctx)
			if err != nil {
				c.log.Error("Login page get error csrf token:", err)
				return ctx.Render("templates/login", fiber.Map{
					"CSRFToken": "",
					"Error":     "Internal server error",
				})
			}

			return ctx.Render("templates/login", fiber.Map{
				"CSRFToken": csrfToken,
				"Error":     "wrong username or password",
			})
		}

		session, err := c.sessStore.Get(ctx)
		if err != nil {
			c.log.Error("Login user error (sesstion error):", err)
			csrfToken, err := c.getCsrfToken(ctx)
			if err != nil {
				c.log.Error("Login page get error csrf token:", err)
				return ctx.Render("templates/login", fiber.Map{
					"CSRFToken": "",
					"Error":     "Internal server error",
				})
			}

			return ctx.Render("templates/login", fiber.Map{
				"CSRFToken": csrfToken,
				"Error":     "wrong username or password",
			})
		}

		session.Set("userId", user.ID)
		session.Set("name", user.Name)
		session.Set("email", user.Email)
		session.Set("username", user.Username)
		session.Set("authenticated", true)

		err = session.Save()
		if err != nil {
			c.log.Error("Login user error:", err)
			csrfToken, err := c.getCsrfToken(ctx)
			if err != nil {
				c.log.Error("Login page get error csrf token:", err)
				return ctx.Render("templates/login", fiber.Map{
					"CSRFToken": "",
					"Error":     "Internal server error",
				})
			}

			return ctx.Render("templates/login", fiber.Map{
				"CSRFToken": csrfToken,
				"Error":     "wrong username or password",
			})
		}

		return ctx.Redirect("/")
	}
}

func (c *controller) Dashboard() FuncHttp {
	return func(ctx *fiber.Ctx) error {
		sess, err := c.sessStore.Get(ctx)
		if err != nil {
			c.log.Error("Dashboard error:", err)
			return ctx.Redirect("/login")
		}

		authenticated := sess.Get("authenticated")
		if authenticated == nil || authenticated == false {
			return ctx.Redirect("/login")
		}

		return ctx.Render("templates/index", fiber.Map{
			"UserId": sess.Get("userId"),
		})
	}
}

func (c *controller) Logout() FuncHttp {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (c *controller) WS() FuncHttp {
	return websocket.New(func(con *websocket.Conn) {
		defer con.Close()

		// Extract session ID from query parameters
		sessionID := con.Cookies("session_id")
		if sessionID == "" {
			c.log.Error("session_id is empty")
			return
		}

		sess, ok := con.Locals("session").(*session.Session)
		if !ok {
			c.log.Error("Failed to get session")
			return
		}

		authenticated := sess.Get("authenticated").(bool)
		if !authenticated {
			return
		}

		userId := sess.Get("userId").(int)
		name := sess.Get("name").(string)

		strUserId := strconv.Itoa(userId)
		c.wsStore.addClient(strUserId, name, con)

		defer c.wsStore.sendUsersForAllClient()
		defer c.wsStore.removeClient(strUserId)

		// Send all clients to the new client
		c.wsStore.sendUsersForAllClient()

		var (
			mt     int // Message type
			msg    []byte
			err    error
			action struct {
				Notify struct {
					UserId string `json:"userId"`
				} `json:"notify"`
			}
		)

		for {
			if mt, msg, err = con.ReadMessage(); err != nil {
				if errors.Is(err, websocket.ErrCloseSent) ||
					mt == websocket.CloseGoingAway {
					return
				}

				c.log.Error("Failed to read message:", err)
				break
			}

			switch mt {
			case websocket.CloseGoingAway:
				return

			case websocket.CloseInternalServerErr:
				return

			case websocket.TextMessage, websocket.BinaryMessage:
				if err := json.Unmarshal(msg, &action); err != nil {
					c.log.Error("Failed to read JSON:", err)
					break
				}

				if action.Notify.UserId != "" {
					c.wsStore.notifyClient(action.Notify.UserId, name)
				}
			}
		}

	})
}
