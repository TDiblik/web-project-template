package utils

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type XValidator struct {
	validate *validator.Validate
}

func (v XValidator) Validate(data interface{}) error {
	return v.validate.Struct(data)
}

var Validator XValidator

func SetupValidator() {
	log.Println("Setting up validator: start")
	Validator = XValidator{
		validate: validator.New(),
	}
	log.Println("Setting up validator: done")
}

func GetValidRequestBody(request_body any, c fiber.Ctx) error {
	if err := c.Bind().Body(request_body); err != nil {
		return err
	}
	if err := Validator.Validate(request_body); err != nil {
		return err
	}
	return nil
}

func GetValidQuery(query any, c fiber.Ctx) error {
	if err := c.Bind().Query(query); err != nil {
		return err
	}
	if err := Validator.Validate(query); err != nil {
		return err
	}
	return nil
}
