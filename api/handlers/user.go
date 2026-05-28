package handlers

import (
	"fmt"
	"strings"

	"github.com/TDiblik/project-template/api/database"
	"github.com/TDiblik/project-template/api/models"
	"github.com/TDiblik/project-template/api/utils"
	"github.com/gofiber/fiber/v3"
)

type GetUserMeHandlerResponse struct {
	UserInfo models.UserModelDB `json:"user_info"`
}

func GetUserMeHandler(c fiber.Ctx) error {
	userJWTInfo, err := utils.GetJWTFromLocals(c)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}

	db, err := database.CreateConnection()
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}

	var user models.UserModelDB
	if err := db.Get(&user, utils.SelectUserById(), userJWTInfo.UserId); err != nil {
		return utils.NotFoundResponse(c, "be.error.user.not_found")
	}

	return utils.OkResponse(c, GetUserMeHandlerResponse{UserInfo: user})
}

type PatchUserMeHandlerRequest struct {
	FirstName        string                          `json:"first_name,omitempty" validate:"omitempty,min=1,max=50"`
	LastName         string                          `json:"last_name,omitempty" validate:"omitempty,min=1,max=50"`
	PreferedTheme    utils.ThemePosibilities         `json:"prefered_theme,omitempty" validate:"omitempty,oneof=dark light"`
	PreferedLanguage utils.TranslationsPossibilities `json:"prefered_language,omitempty" validate:"omitempty,oneof=cs en"`
	Password         string                          `json:"password,omitempty" validate:"omitempty,min=6"`
}
type PatchUserMeHandlerResponse struct{}

func PatchUserMeHandler(c fiber.Ctx) error {
	userJWTInfo, err := utils.GetJWTFromLocals(c)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}

	var req PatchUserMeHandlerRequest
	if err := utils.GetValidRequestBody(&req, c); err != nil {
		return utils.InvalidRequestResponse(c, err)
	}

	db, err := database.CreateConnection()
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}

	var user models.UserModelDB
	if err := db.Get(&user, utils.SelectUserById(), userJWTInfo.UserId); err != nil {
		return utils.NotFoundResponse(c, "be.error.user.not_found")
	}

	if req.FirstName != "" {
		user.FirstName = utils.SQLNullStringFromString(req.FirstName)
	}
	if req.LastName != "" {
		user.LastName = utils.SQLNullStringFromString(req.LastName)
	}
	if req.PreferedTheme != "" {
		user.PreferedTheme = utils.SQLNullStringFromString(string(req.PreferedTheme))
	}
	if req.PreferedLanguage != "" {
		user.PreferedLanguage = utils.SQLNullStringFromString(string(req.PreferedLanguage))
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return utils.InternalServerErrorResponse(c, err)
		}
		user.PasswordHash = utils.SQLNullStringFromString(hashedPassword)
	}

	_, err = db.NamedExec(`
		update users set
			first_name = :first_name,
			last_name = :last_name,
			prefered_theme = :prefered_theme,
			prefered_language = :prefered_language,
			password_hash = :password_hash
		where id = :id
	`, user)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}

	return utils.OkResponse(c, PatchUserMeHandlerResponse{})
}

type PostUserMeAvatarHandlerResponse struct{}

func PostUserMeAvatarHandler(c fiber.Ctx) error {
	file, err := c.FormFile("avatar")
	if err != nil {
		return utils.InvalidRequestResponse(c, err)
	}

	userJWTInfo, err := utils.GetJWTFromLocals(c)
	if err != nil {
		return utils.InvalidRequestResponse(c, err)
	}

	db, err := database.CreateConnection()
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}

	var user models.UserModelDB
	if err := db.Get(&user, utils.SelectUserById(), userJWTInfo.UserId); err != nil {
		return utils.NotFoundResponse(c, "be.error.user.not_found")
	}

	newAvatarImageUUID, err := utils.SaveImage(c, file, utils.GetAvatarImageFolder(), 450, 450)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}
	newAvatarImageUrlPath, err := utils.GetAvatarImageUrl(newAvatarImageUUID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}
	user.AvatarUrl = utils.SQLNullStringFromString(newAvatarImageUrlPath)

	_, err = db.NamedExec(`update users set avatar_url = :avatar_url where id = :id`, user)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}

	return utils.OkResponse(c, PostUserMeAvatarHandlerResponse{})
}

type DeleteUserOauthHandlerResponse struct{}

func DeleteUserOauthHandler(c fiber.Ctx) error {
	userJWTInfo, err := utils.GetJWTFromLocals(c)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}

	provider := c.Params("provider")
	if provider == "" {
		return utils.InvalidRequestResponse(c, fmt.Errorf("provider is required"))
	}

	db, err := database.CreateConnection()
	if err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}

	var user models.UserModelDB
	if err := db.Get(&user, utils.SelectUserById(), userJWTInfo.UserId); err != nil {
		return utils.NotFoundResponse(c, "user not found")
	}

	authMethodsCount := 0
	if user.PasswordHash.Valid {
		authMethodsCount++
	}
	if user.GithubId.Valid {
		authMethodsCount++
	}
	if user.GoogleId.Valid {
		authMethodsCount++
	}
	if user.FacebookId.Valid {
		authMethodsCount++
	}
	if user.SpotifyId.Valid {
		authMethodsCount++
	}

	if authMethodsCount <= 1 {
		return utils.ConflictResponse(c, "be.error.user.cannot_remove_last_auth_method")
	}

	var updateQuery string
	switch strings.ToLower(provider) {
	case "github":
		if !user.GithubId.Valid {
			return utils.InvalidRequestResponse(c, fmt.Errorf("provider not connected"))
		}
		updateQuery = `update users set github_id = null, github_email = null, github_handle = null, github_url = null where id = $1`
	case "google":
		if !user.GoogleId.Valid {
			return utils.InvalidRequestResponse(c, fmt.Errorf("provider not connected"))
		}
		updateQuery = `update users set google_id = null, google_email = null where id = $1`
	case "facebook":
		if !user.FacebookId.Valid {
			return utils.InvalidRequestResponse(c, fmt.Errorf("provider not connected"))
		}
		updateQuery = `update users set facebook_id = null, facebook_email = null, facebook_url = null where id = $1`
	case "spotify":
		if !user.SpotifyId.Valid {
			return utils.InvalidRequestResponse(c, fmt.Errorf("provider not connected"))
		}
		updateQuery = `update users set spotify_id = null, spotify_email = null, spotify_url = null where id = $1`
	default:
		return utils.InvalidRequestResponse(c, fmt.Errorf("invalid provider %s", provider))
	}

	if _, err := db.Exec(updateQuery, userJWTInfo.UserId); err != nil {
		return utils.InternalServerErrorResponse(c, err)
	}

	return utils.OkResponse(c, DeleteUserOauthHandlerResponse{})
}
