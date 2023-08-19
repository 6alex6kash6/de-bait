package user

import (
	"fmt"
	"log"

	"github.com/de-bait/internal/app/utils"
	"github.com/de-bait/internal/db/repository"
	"github.com/de-bait/internal/middlewares"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	repo     *repository.UserRepository
	app      *fiber.App
	validate *validator.Validate
}

func NewUserHandler(ur *repository.UserRepository, app *fiber.App, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		repo:     ur,
		app:      app,
		validate: validate,
	}
}

func (uh *UserHandler) FindOne(c *fiber.Ctx) error {
	user := c.Locals("user")

	return c.JSON(user)
}

func (uh *UserHandler) Create(c *fiber.Ctx) error {
	ui := new(repository.UserInput)
	c.BodyParser(ui)

	if err := uh.validate.Struct(ui); err != nil {
		return fmt.Errorf("failed creating user: %w", err)
	}

	u, err := uh.repo.Create(*ui)
	if err != nil {
		return fmt.Errorf("failed creating user: %w", err)
	}
	log.Println("user was created: ", u)
	return c.JSON(u)
}

func (uh *UserHandler) FindMany(c *fiber.Ctx) error {
	u, err := uh.repo.Find()
	if err != nil {
		return fmt.Errorf("failed getting users: %w", err)
	}
	return c.JSON(u)
}

func (uh *UserHandler) Update(c *fiber.Ctx) error {
	id := utils.ParseId(c.Params("id"))

	ui := new(repository.UserInput)
	c.BodyParser(ui)

	if err := uh.validate.Struct(ui); err != nil {
		return fmt.Errorf("failed creating user: %w", err)
	}

	u, err := uh.repo.Update(id, *ui)
	if err != nil {
		return fmt.Errorf("failed updating user: %w", err)
	}
	log.Println("user was updated: ", u)
	return c.JSON(u)
}

func (uh *UserHandler) Login(c *fiber.Ctx) error {
	ui := new(repository.UserInput)
	c.BodyParser(ui)

	if err := uh.validate.Struct(ui); err != nil {
		return fmt.Errorf("failed creating user: %w", err)
	}

	u, err := uh.repo.FindOneByNickname(ui.Nickname)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return fmt.Errorf("failed finding user: %w", err)
	}
	if u.Password != ui.Password {
		c.Status(fiber.StatusUnauthorized).SendString("wrong password")
		return fmt.Errorf("wrong password")
	}

	jwt, err := utils.CreateJwt(u.ID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		return fmt.Errorf("failed creating jwt: %w", err)
	}

	return c.JSON(jwt)
}

func (uh *UserHandler) Register() {
	auth := middlewares.Auth(*uh.repo)
	uh.app.Post("/user", uh.Create)
	uh.app.Post("/user/login", uh.Login)
	uh.app.Post("/user/:id", auth, uh.Update)
	uh.app.Get("/user", auth, uh.FindOne)
	uh.app.Get("/users/", auth, uh.FindMany)
}
