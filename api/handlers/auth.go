package handlers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/TDiblik/project-template/api/database"
	"github.com/TDiblik/project-template/api/models"
	"github.com/TDiblik/project-template/api/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type LoginHandlerRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthHandlerResponse struct {
	AuthToken    string `json:"auth_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func LoginHandler(c fiber.Ctx) error {
	var req LoginHandlerRequestBody
	if err := utils.GetValidRequestBody(&req, c); err != nil {
		return utils.InvalidRequestResponse(c, err)
	}

	db, err := database.CreateConnection()
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("unable to create db connection: %w", err))
	}

	var userUUID uuid.UUID
	var userPasswordHash models.SQLNullString
	if err := db.QueryRow(utils.SelectIdAndPasswordHashByEmailQuery(), req.Email).Scan(&userUUID, &userPasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return utils.ConflictResponse(c, "be.error.login.username_or_password_incorrect")
		}
		return utils.InternalServerErrorResponse(c, err)
	}

	// Handle accounts created by oauth without a password
	if !userPasswordHash.Valid || userPasswordHash.String == "" {
		return utils.ConflictResponse(c, "be.error.login.username_or_password_incorrect")
	}

	if !utils.CompareHashAndPassword(userPasswordHash.String, req.Password) {
		return utils.ConflictResponse(c, "be.error.login.username_or_password_incorrect")
	}

	authToken, refreshToken, err := GetJwtPostLogin(userUUID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("unable to execute GetJwtPostLogin: %w", err))
	}

	return utils.OkResponse(c, AuthHandlerResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	})
}

type SignUpHandlerRequestBody struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	UseUsername bool   `json:"useUsername"`
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Username    string `json:"username,omitempty"`
}

func (req *SignUpHandlerRequestBody) validateSignUpRequestBody() error {
	if len([]byte(req.Password)) > 72 {
		return fmt.Errorf("password cannot exceed 72 bytes")
	}
	if req.UseUsername {
		if strings.TrimSpace(req.Username) == "" {
			return fmt.Errorf("username is required")
		}
		if len(req.Username) < 3 {
			return fmt.Errorf("username must have at least 3 characters")
		}
	} else {
		if strings.TrimSpace(req.FirstName) == "" {
			return fmt.Errorf("first name is required")
		}
		if strings.TrimSpace(req.LastName) == "" {
			return fmt.Errorf("last name is required")
		}
	}
	return nil
}

func SignUpHandler(c fiber.Ctx) error {
	var req SignUpHandlerRequestBody
	if err := utils.GetValidRequestBody(&req, c); err != nil {
		return utils.InvalidRequestResponse(c, err)
	}
	if err := req.validateSignUpRequestBody(); err != nil {
		return utils.InvalidRequestResponse(c, err)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("failed to hash password: %w", err))
	}

	db, err := database.CreateConnection()
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("unable to create db connection: %w", err))
	}

	var emailExists bool
	if err := db.QueryRow("select "+utils.UserEmailExistsQuery(), req.Email).Scan(&emailExists); err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("failed to check existing user: %w", err))
	}
	if emailExists {
		return utils.ConflictResponse(c, "be.error.login.email_already_in_use")
	}

	var newUser = models.UserModelDB{
		Email:         req.Email,
		EmailVerified: false,
		FirstName:     utils.SQLNullStringFromString(req.FirstName),
		LastName:      utils.SQLNullStringFromString(req.LastName),
		PasswordHash:  utils.SQLNullStringFromString(hashedPassword),
	}
	if req.UseUsername {
		var handleExists bool
		if err := db.QueryRow(`select exists(select 1 from users where handle = $1)`, req.Username).Scan(&handleExists); err != nil {
			return utils.InternalServerErrorResponse(c, fmt.Errorf("failed to check existing user: %w", err))
		}
		if handleExists {
			return utils.ConflictResponse(c, "be.error.login.handle_already_in_use")
		}
		newUser.Handle = utils.SQLNullStringFromString(req.Username)
	}

	userUUID, err := CreateOrUpdateUser(c, newUser)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}
	authToken, refreshToken, err := GetJwtPostLogin(userUUID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("unable to execute GetJwtPostLogin: %w", err))
	}

	return utils.OkResponse(c, AuthHandlerResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	})
}

type RefreshTokenRequestBody struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func RefreshTokenHandler(c fiber.Ctx) error {
	var req RefreshTokenRequestBody
	if err := utils.GetValidRequestBody(&req, c); err != nil {
		return utils.InvalidRequestResponse(c, err)
	}

	userId, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return utils.UnauthorizedResponse(c, fmt.Errorf("be.error.auth.invalid_refresh_token"))
	}

	authToken, refreshToken, err := GetJwtPostLogin(userId)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("unable to execute GetJwtPostLogin: %w", err))
	}

	return utils.OkResponse(c, AuthHandlerResponse{
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	})
}
