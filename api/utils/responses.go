package utils

import (
	"fmt"

	"github.com/TDiblik/gofiber-swagger/gofiberswagger"
	"github.com/gofiber/fiber/v3"
)

type ErrorResponseType struct {
	Status string `json:"status" validate:"required"`
	Msg    string `json:"msg" validate:"required"`
}

var DefaultErrorResponses = []gofiberswagger.ResponseInfo{
	gofiberswagger.NewResponseInfo[ErrorResponseType]("400", "invalid request"),
	gofiberswagger.NewResponseInfo[ErrorResponseType]("401", "user unauthenticated"),
	gofiberswagger.NewResponseInfo[ErrorResponseType]("403", "user unauthorized"),
	gofiberswagger.NewResponseInfo[ErrorResponseType]("404", "not found"),
	gofiberswagger.NewResponseInfo[ErrorResponseType]("409", "conflicting request"),
	gofiberswagger.NewResponseInfo[ErrorResponseType]("429", "too many requests"),
	gofiberswagger.NewResponseInfo[ErrorResponseType]("500", "internal server error"),
}

func NewSwaggerResponsesWithErrors(responses ...gofiberswagger.ResponseInfo) *gofiberswagger.Responses {
	return gofiberswagger.NewResponses(append(responses, DefaultErrorResponses...)...)
}

func InvalidRequestResponse(c fiber.Ctx, e error) error {
	Logger.Errorw("ERROR", "error", e)
	if EnvData.Debug {
		return c.Status(400).JSON(fiber.Map{"status": "error", "msg": "be.error.invalid_request", "message": "Review your input", "error_info": fmt.Sprint(e)})
	}
	return c.Status(400).JSON(fiber.Map{"status": "error", "msg": "be.error.invalid_request"})
}

func InternalServerErrorResponse(c fiber.Ctx, e error) error {
	Logger.Errorw("ERROR", "error", e)
	if EnvData.Debug {
		return c.Status(500).JSON(fiber.Map{"status": "error", "msg": "be.error.internal_server_error", "message": "Internal Server Error", "error_info": fmt.Sprint(e)})
	}
	return c.Status(500).JSON(fiber.Map{"status": "error", "msg": "be.error.internal_server_error", "message": "Internal Server Error"})
}

func UnauthentizatedResponse(c fiber.Ctx, e error) error {
	Logger.Errorw("ERROR", "error", e)
	if EnvData.Debug {
		return c.Status(401).JSON(fiber.Map{"status": "unauthenticated", "msg": "be.error.invalid_token", "message": "valid token is required", "error_info": fmt.Sprint(e)})
	}
	return c.Status(401).JSON(fiber.Map{"status": "unauthenticated", "msg": "be.error.invalid_token"})
}

func UnauthorizedResponse(c fiber.Ctx, e error) error {
	Logger.Errorw("ERROR", "error", e)
	if EnvData.Debug {
		return c.Status(403).JSON(fiber.Map{"status": "unauthorized", "msg": "be.error.unauthorized", "message": "you cannot access this endpoint", "error_info": fmt.Sprint(e)})
	}
	return c.Status(403).JSON(fiber.Map{"status": "unauthorized", "msg": "be.error.unauthorized"})
}

func ConflictResponse(c fiber.Ctx, msg_reason string) error {
	Logger.Info("conflict: " + msg_reason)
	return c.Status(409).JSON(fiber.Map{"status": "conflict", "msg": msg_reason})
}

func NotFoundResponse(c fiber.Ctx, msg_reason string) error {
	return c.Status(404).JSON(fiber.Map{"status": "not found", "msg": msg_reason})
}

func TooManyRequestsResponse(c fiber.Ctx, e error) error {
	Logger.Errorw("ERROR", "error", e)
	if EnvData.Debug {
		return c.Status(429).JSON(fiber.Map{"status": "too_many_requests", "msg": "be.error.rate_limit_exceeded", "message": "rate limit exceeded", "error_info": fmt.Sprint(e)})
	}
	return c.Status(429).JSON(fiber.Map{"status": "too_many_requests", "msg": "be.error.rate_limit_exceeded"})
}

func OkResponse(c fiber.Ctx, data any) error {
	return c.Status(200).JSON(data)
}
