package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/TDiblik/project-template/api/database"
	"github.com/TDiblik/project-template/api/models"
	"github.com/TDiblik/project-template/api/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	githubProviderName   = "GitHub"
	googleProviderName   = "Google"
	facebookProviderName = "Facebook"
	spotifyProviderName  = "Spotify"
)

type OauthRedirectHandlerResponse struct {
	OAuthState  string `json:"oauth_state"`
	RedirectURL string `json:"redirect_url"`
}

func GithubRedirectHandler(c fiber.Ctx) error {
	redirectParam := c.Query("redirect_back_to_after_oauth", string(utils.RedirectAfterOauthIndex))
	redirectBackTo := utils.ValidateRedirectAfterOauth(redirectParam)
	isMobile := c.Query("isMobile") == "true"

	state, err := utils.GenerateOauthState(githubProviderName, redirectBackTo)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("failed to generate OAuth state: %w", err))
	}

	config := utils.EnvData.OAUTH_GITHUB_CONFIG
	if isMobile {
		config = utils.EnvData.OAUTH_GITHUB_CONFIG_MOBILE
	}

	return utils.OkResponse(c, OauthRedirectHandlerResponse{
		OAuthState:  state,
		RedirectURL: config.AuthCodeURL(state),
	})
}

func GoogleRedirectHandler(c fiber.Ctx) error {
	redirectParam := c.Query("redirect_back_to_after_oauth", string(utils.RedirectAfterOauthIndex))
	redirectBackTo := utils.ValidateRedirectAfterOauth(redirectParam)
	isMobile := c.Query("isMobile") == "true"

	state, err := utils.GenerateOauthState(googleProviderName, redirectBackTo)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("failed to generate OAuth state: %w", err))
	}

	config := utils.EnvData.OAUTH_GOOGLE_CONFIG
	if isMobile {
		config = utils.EnvData.OAUTH_GOOGLE_CONFIG_MOBILE
	}

	return utils.OkResponse(c, OauthRedirectHandlerResponse{
		OAuthState:  state,
		RedirectURL: config.AuthCodeURL(state),
	})
}

func FacebookRedirectHandler(c fiber.Ctx) error {
	redirectParam := c.Query("redirect_back_to_after_oauth", string(utils.RedirectAfterOauthIndex))
	redirectBackTo := utils.ValidateRedirectAfterOauth(redirectParam)
	isMobile := c.Query("isMobile") == "true"

	state, err := utils.GenerateOauthState(facebookProviderName, redirectBackTo)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("failed to generate OAuth state: %w", err))
	}

	config := utils.EnvData.OAUTH_FACEBOOK_CONFIG
	if isMobile {
		config = utils.EnvData.OAUTH_FACEBOOK_CONFIG_MOBILE
	}

	return utils.OkResponse(c, OauthRedirectHandlerResponse{
		OAuthState:  state,
		RedirectURL: config.AuthCodeURL(state),
	})
}

func SpotifyRedirectHandler(c fiber.Ctx) error {
	redirectParam := c.Query("redirect_back_to_after_oauth", string(utils.RedirectAfterOauthIndex))
	redirectBackTo := utils.ValidateRedirectAfterOauth(redirectParam)
	isMobile := c.Query("isMobile") == "true"

	state, err := utils.GenerateOauthState(spotifyProviderName, redirectBackTo)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("failed to generate OAuth state: %w", err))
	}

	config := utils.EnvData.OAUTH_SPOTIFY_CONFIG
	if isMobile {
		config = utils.EnvData.OAUTH_SPOTIFY_CONFIG_MOBILE
	}

	return utils.OkResponse(c, OauthRedirectHandlerResponse{
		OAuthState:  state,
		RedirectURL: config.AuthCodeURL(state),
	})
}

type OAuthPostReturnHandlerQuery struct {
	State string `query:"state" validate:"required"`
	Code  string `query:"code" validate:"required"`
}

type OAuthPostReturnHandlerResponse struct {
	AuthToken                string                   `json:"auth_token" validate:"required"`
	RefreshToken             string                   `json:"refresh_token" validate:"required"`
	RedirectBackToAfterOauth utils.RedirectAfterOauth `json:"redirect_back_to_after_oauth" validate:"required"`
}

func OAuthPostReturnHandler(c fiber.Ctx) error {
	var query OAuthPostReturnHandlerQuery
	if err := utils.GetValidQuery(&query, c); err != nil {
		return utils.InvalidRequestResponse(c, err)
	}

	if query.State == "" || !utils.IsValidOauthState(query.State) {
		return utils.UnauthorizedResponse(c, fmt.Errorf("invalid OAuth state"))
	}

	provider, redirect, err := utils.GetOauthProviderAndRedirectFromOauthState(query.State)
	if err != nil {
		return utils.InvalidRequestResponse(c, fmt.Errorf("invalid provider or redirect info inside the state: %w", err))
	}

	var userUUID uuid.UUID
	// when adding a new oauth provider and user table fields, add new "case" here:
	switch provider {
	case githubProviderName:
		if userUUID, err = githubReturn(c, query.Code); err != nil {
			return utils.InternalServerErrorResponse(c, err)
		}
	case googleProviderName:
		if userUUID, err = googleReturn(c, query.Code); err != nil {
			return utils.InternalServerErrorResponse(c, err)
		}
	case facebookProviderName:
		if userUUID, err = facebookReturn(c, query.Code); err != nil {
			return utils.InternalServerErrorResponse(c, err)
		}
	case spotifyProviderName:
		if userUUID, err = spotifyReturn(c, query.Code); err != nil {
			return utils.InternalServerErrorResponse(c, err)
		}
	default:
		return utils.InvalidRequestResponse(c, fmt.Errorf("invalid provider name (not implemented) inside the state: %w", err))
	}

	authToken, refreshToken, err := GetJwtPostLogin(userUUID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, fmt.Errorf("unable to execute GetJwtPostLogin: %w", err))
	}

	return utils.OkResponse(c, OAuthPostReturnHandlerResponse{
		AuthToken:                authToken,
		RefreshToken:             refreshToken,
		RedirectBackToAfterOauth: redirect,
	})
}

type githubUserApiResponse struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	HtmlUrl   string `json:"html_url"`
}

func githubReturn(c fiber.Ctx, authCode string) (uuid.UUID, error) {
	token, err := utils.EnvData.OAUTH_GITHUB_CONFIG.Exchange(context.Background(), authCode)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := utils.EnvData.OAUTH_GITHUB_CONFIG.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return uuid.Nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var ghUserResponse githubUserApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghUserResponse); err != nil {
		return uuid.Nil, fmt.Errorf("failed to decode user info: %w", err)
	}
	var firstName, lastName string
	if ghUserResponse.Name != "" {
		parts := strings.SplitN(ghUserResponse.Name, " ", 2)
		firstName = parts[0]
		if len(parts) > 1 {
			lastName = parts[1]
		}
	}

	return CreateOrUpdateUser(c, models.UserModelDB{
		Email:         ghUserResponse.Email,
		EmailVerified: true,
		FirstName:     utils.SQLNullStringFromString(firstName),
		LastName:      utils.SQLNullStringFromString(lastName),
		Handle:        utils.SQLNullStringFromString(ghUserResponse.Login),
		GithubId:      utils.SQLNullStringFromString(strconv.FormatInt(ghUserResponse.ID, 10)),
		GithubHandle:  utils.SQLNullStringFromString(ghUserResponse.Login),
		GithubEmail:   utils.SQLNullStringFromString(ghUserResponse.Email),
		GithubUrl:     utils.SQLNullStringFromString(ghUserResponse.HtmlUrl),
		AvatarUrl:     utils.SQLNullStringFromString(ghUserResponse.AvatarURL),
	})
}

type googleUserApiResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func googleReturn(c fiber.Ctx, authCode string) (uuid.UUID, error) {
	token, err := utils.EnvData.OAUTH_GOOGLE_CONFIG.Exchange(context.Background(), authCode)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := utils.EnvData.OAUTH_GOOGLE_CONFIG.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return uuid.Nil, fmt.Errorf("google API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var gUser googleUserApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		return uuid.Nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return CreateOrUpdateUser(c, models.UserModelDB{
		Email:         gUser.Email,
		EmailVerified: gUser.VerifiedEmail,
		FirstName:     utils.SQLNullStringFromString(gUser.GivenName),
		LastName:      utils.SQLNullStringFromString(gUser.FamilyName),
		GoogleId:      utils.SQLNullStringFromString(gUser.ID),
		GoogleEmail:   utils.SQLNullStringFromString(gUser.Email),
		AvatarUrl:     utils.SQLNullStringFromString(gUser.Picture),
	})
}

type facebookUserApiResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Picture   struct {
		Data struct {
			Height       int    `json:"height"`
			IsSilhouette bool   `json:"is_silhouette"`
			URL          string `json:"url"`
			Width        int    `json:"width"`
		} `json:"data"`
	} `json:"picture"`
	Link string `json:"link"`
}

func facebookReturn(c fiber.Ctx, authCode string) (uuid.UUID, error) {
	token, err := utils.EnvData.OAUTH_FACEBOOK_CONFIG.Exchange(context.Background(), authCode)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := utils.EnvData.OAUTH_FACEBOOK_CONFIG.Client(context.Background(), token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,first_name,last_name,email,picture,link")
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return uuid.Nil, fmt.Errorf("facebook API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var fbUser facebookUserApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&fbUser); err != nil {
		return uuid.Nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return CreateOrUpdateUser(c, models.UserModelDB{
		Email:         fbUser.Email,
		EmailVerified: true,
		FirstName:     utils.SQLNullStringFromString(fbUser.FirstName),
		LastName:      utils.SQLNullStringFromString(fbUser.LastName),
		FacebookId:    utils.SQLNullStringFromString(fbUser.ID),
		FacebookEmail: utils.SQLNullStringFromString(fbUser.Email),
		FacebookUrl:   utils.SQLNullStringFromString(fbUser.Link),
		AvatarUrl:     utils.SQLNullStringFromString(fbUser.Picture.Data.URL),
	})
}

type spotifyUserApiResponse struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Images      []struct {
		URL string `json:"url"`
	} `json:"images"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

func spotifyReturn(c fiber.Ctx, authCode string) (uuid.UUID, error) {
	token, err := utils.EnvData.OAUTH_SPOTIFY_CONFIG.Exchange(context.Background(), authCode)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := utils.EnvData.OAUTH_SPOTIFY_CONFIG.Client(context.Background(), token)
	resp, err := client.Get("https://api.spotify.com/v1/me")
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return uuid.Nil, fmt.Errorf("spotify API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var spUser spotifyUserApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&spUser); err != nil {
		return uuid.Nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	var firstName, lastName string
	if spUser.DisplayName != "" {
		parts := strings.SplitN(spUser.DisplayName, " ", 2)
		firstName = parts[0]
		if len(parts) > 1 {
			lastName = parts[1]
		}
	}

	var avatarURL string
	if len(spUser.Images) > 0 {
		avatarURL = spUser.Images[0].URL
	}

	return CreateOrUpdateUser(c, models.UserModelDB{
		Email:         spUser.Email,
		EmailVerified: true,
		FirstName:     utils.SQLNullStringFromString(firstName),
		LastName:      utils.SQLNullStringFromString(lastName),
		SpotifyId:     utils.SQLNullStringFromString(spUser.ID),
		SpotifyEmail:  utils.SQLNullStringFromString(spUser.Email),
		SpotifyUrl:    utils.SQLNullStringFromString(spUser.ExternalURLs.Spotify),
		AvatarUrl:     utils.SQLNullStringFromString(avatarURL),
	})
}

func CreateOrUpdateUser(c fiber.Ctx, possiblyNewUser models.UserModelDB) (uuid.UUID, error) {
	db, err := database.CreateConnection()
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to create connection to db inside CreateOrUpdateUser: %w", err)
	}

	possiblyUserEmail := possiblyNewUser.Email
	loggedInUserInfo, err := utils.GetUserInfoFromJWT(c)
	processingLoggedInUser := loggedInUserInfo != nil
	if err == nil && processingLoggedInUser {
		possiblyUserEmail = loggedInUserInfo.UserEmail
	}

	var emailExists, handleExists bool
	err = db.QueryRow(`select `+utils.UserEmailExistsQuery()+` as email_exists, exists(select 1 from users where handle = $2) as handle_exists`,
		possiblyUserEmail, possiblyNewUser.Handle).Scan(&emailExists, &handleExists)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to query exists staments inside CreateOrUpdateUser: %w", err)
	}

	processAvatar := func(avatarUrl string) string {
		if strings.HasPrefix(avatarUrl, "http") {
			newImageId, err := utils.DownloadAndSaveImage(avatarUrl, utils.EnvData.IMAGES_PATH_AVATAR, 512, 512)
			if err == nil {
				localUrl, _ := utils.GetAvatarImageUrl(newImageId)
				return localUrl
			}
			utils.Logger.Errorw("Failed to download avatar", "url", avatarUrl, "error", err)
		}
		return avatarUrl
	}

	var existingUser models.UserModelDB
	if !emailExists {
		if !possiblyNewUser.Handle.Valid || handleExists {
			newHandle, genErr := utils.GenerateUniqueUserHandle(db, possiblyNewUser.FirstName, possiblyNewUser.LastName)
			if genErr != nil {
				return uuid.Nil, fmt.Errorf("unable to generate unique handle: %w", genErr)
			}
			possiblyNewUser.Handle = utils.SQLNullStringFromString(newHandle)
		}

		if possiblyNewUser.AvatarUrl.Valid {
			possiblyNewUser.AvatarUrl.String = processAvatar(possiblyNewUser.AvatarUrl.String)
		}

		// when adding a new oauth provider and user table fields, add the checks here:
		rows, err := db.NamedQuery(`
			insert into users (
				email, email_verified, password_hash, handle, first_name, last_name, avatar_url, github_id, github_email, github_handle, github_url, google_id, google_email, facebook_id, facebook_email, facebook_url, spotify_id, spotify_email, spotify_url
			)
			values (
				:email, :email_verified, :password_hash, :handle, :first_name, :last_name, :avatar_url, :github_id, :github_email, :github_handle, :github_url, :google_id, :google_email, :facebook_id, :facebook_email, :facebook_url, :spotify_id, :spotify_email, :spotify_url
			)
			returning *
		`, possiblyNewUser)
		if err != nil {
			return uuid.Nil, fmt.Errorf("unable to insert new user: %w", err)
		}
		defer rows.Close()
		if !rows.Next() {
			return uuid.Nil, fmt.Errorf("no user inserted, but no error provided ? this should never occured, user object: %+v", possiblyNewUser)
		}
		if err := rows.StructScan(&existingUser); err != nil {
			return uuid.Nil, fmt.Errorf("unable to scan inserted user: %w", err)
		}
	}
	if emailExists {
		if processingLoggedInUser {
			err := db.Get(&existingUser, utils.SelectUserByEmailQuery(), possiblyNewUser.Email)
			if err != nil {
				err := db.Get(&existingUser, utils.SelectUserByEmailQuery(), possiblyUserEmail)
				if err != nil {
					return uuid.Nil, fmt.Errorf("unable to select existing user: %w", err)
				}
			}
			if loggedInUserInfo.UserId != existingUser.Id {
				return uuid.Nil, fmt.Errorf("email is already linked to an account: %v (logged in user) VS %v (existing user)", loggedInUserInfo, existingUser)
			}
		} else {
			err := db.Get(&existingUser, utils.SelectUserByEmailQuery(), possiblyUserEmail)
			if err != nil {
				return uuid.Nil, fmt.Errorf("unable to select existing user: %w", err)
			}
		}

		if !existingUser.PasswordHash.Valid {
			existingUser.PasswordHash = possiblyNewUser.PasswordHash
		}
		if !existingUser.Handle.Valid && !handleExists {
			existingUser.Handle = possiblyNewUser.Handle
		}
		if !existingUser.FirstName.Valid {
			existingUser.FirstName = possiblyNewUser.FirstName
		}
		if !existingUser.LastName.Valid {
			existingUser.LastName = possiblyNewUser.LastName
		}
		if !existingUser.AvatarUrl.Valid && possiblyNewUser.AvatarUrl.Valid {
			existingUser.AvatarUrl.String = processAvatar(possiblyNewUser.AvatarUrl.String)
			existingUser.AvatarUrl.Valid = true
		}
		if !existingUser.EmailVerified {
			existingUser.EmailVerified = possiblyNewUser.EmailVerified
		}

		// when adding a new oauth provider and user table fields, add the checks here:
		if !existingUser.GithubId.Valid {
			existingUser.GithubId = possiblyNewUser.GithubId
		}
		if !existingUser.GithubEmail.Valid {
			existingUser.GithubEmail = possiblyNewUser.GithubEmail
		}
		if !existingUser.GithubHandle.Valid {
			existingUser.GithubHandle = possiblyNewUser.GithubHandle
		}
		if !existingUser.GithubUrl.Valid {
			existingUser.GithubUrl = possiblyNewUser.GithubUrl
		}
		if !existingUser.GoogleId.Valid {
			existingUser.GoogleId = possiblyNewUser.GoogleId
		}
		if !existingUser.GoogleEmail.Valid {
			existingUser.GoogleEmail = possiblyNewUser.GoogleEmail
		}
		if !existingUser.FacebookId.Valid {
			existingUser.FacebookId = possiblyNewUser.FacebookId
		}
		if !existingUser.FacebookEmail.Valid {
			existingUser.FacebookEmail = possiblyNewUser.FacebookEmail
		}
		if !existingUser.FacebookUrl.Valid {
			existingUser.FacebookUrl = possiblyNewUser.FacebookUrl
		}
		if !existingUser.SpotifyId.Valid {
			existingUser.SpotifyId = possiblyNewUser.SpotifyId
		}
		if !existingUser.SpotifyEmail.Valid {
			existingUser.SpotifyEmail = possiblyNewUser.SpotifyEmail
		}
		if !existingUser.SpotifyUrl.Valid {
			existingUser.SpotifyUrl = possiblyNewUser.SpotifyUrl
		}

		// when adding a new oauth provider and user table fields, change the query here:
		if _, err := db.NamedExec(`
			update users set
				password_hash = :password_hash,
				handle = :handle,
				first_name = :first_name,
				last_name = :last_name,
				avatar_url = :avatar_url,
				email_verified = :email_verified,
				github_id = :github_id,
				github_email = :github_email,
				github_handle = :github_handle,
				github_url = :github_url,
				google_id = :google_id,
				google_email = :google_email,
				facebook_id = :facebook_id,
				facebook_email = :facebook_email,
				facebook_url = :facebook_url,
				spotify_id = :spotify_id,
				spotify_email = :spotify_email,
				spotify_url = :spotify_url
			where id = :id
		`, existingUser); err != nil {
			return uuid.Nil, fmt.Errorf("unable to update existing user: %w", err)
		}
	}

	return existingUser.Id, nil
}

// This MUST always be called ONLY IF THE USER WAS 100% AUTHENTICATED USING OAUTH OR PASSWORD
func GetJwtPostLogin(userUUID uuid.UUID) (string, string, error) {
	db, err := database.CreateConnection()
	if err != nil {
		return "", "", fmt.Errorf("unable to create connection to db inside GetJwtPostLogin: %w", err)
	}
	var existingUser models.UserModelDB
	if err := db.Get(&existingUser, `update users set last_login_at = now() where id = $1 returning *`, userUUID); err != nil {
		return "", "", fmt.Errorf("unable to get existingUser inside GetJwtPostLogin: %w", err)
	}

	authToken, refreshToken, err := utils.GenerateTokens(existingUser)
	if err != nil {
		return "", "", fmt.Errorf("unable to GenerateTokens inside GetJwtPostLogin: %w", err)
	}

	return authToken, refreshToken, nil
}
