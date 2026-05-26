package utils

import (
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"

	"github.com/TDiblik/project-template/api/constants"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func SaveImage(c fiber.Ctx, file *multipart.FileHeader, newFolderPath string, maxW int, maxH int) (string, error) {
	tempFilePath := GetTempImagePath(uuid.New().String())
	if err := c.SaveFile(file, tempFilePath); err != nil {
		return "", err
	}
	defer func() {
		if err := os.Remove(tempFilePath); err != nil {
			LogErr(err)
		}
	}()

	tempSavedSrc, err := imaging.Open(tempFilePath, imaging.AutoOrientation(true))
	if err != nil {
		return "", err
	}
	currentBounds := tempSavedSrc.Bounds()
	shouldResize := (maxW > 0 && currentBounds.Dx() > maxW) || (maxH > 0 && currentBounds.Dy() > maxH)
	if shouldResize {
		tempSavedSrc = imaging.Fit(tempSavedSrc, maxW, maxH, imaging.Lanczos)
	}

	newImageId := uuid.New().String()
	newFilePath := filepath.Join(newFolderPath, AddImageExtensionIfNeeded(newImageId))
	if err := imaging.Save(tempSavedSrc, newFilePath, imaging.JPEGQuality(75)); err != nil {
		return "", err
	}

	return newImageId, nil
}

func AddImageExtensionIfNeeded(fileName string) string {
	if filepath.Ext(fileName) == "" {
		fileName += ".jpg"
	}
	return fileName
}

func GetAvatarImageFolder() string {
	return filepath.Join(EnvData.IMAGES_PATH, "avatar/")
}
func GetAvatarImagePath(fileName string) string {
	return filepath.Join(GetAvatarImageFolder(), AddImageExtensionIfNeeded(fileName))
}
func GetAvatarImageUrlBase() (string, error) {
	return url.JoinPath(EnvData.API_PROD_URL, constants.IMAGES_PATH_PREFIX_FULL, "avatar/")
}
func GetAvatarImageUrl(fileName string) (string, error) {
	base, err := GetAvatarImageUrlBase()
	if err != nil {
		return "", err
	}
	return url.JoinPath(base, AddImageExtensionIfNeeded(fileName))
}

func GetTempImageFolder() string {
	return filepath.Join(EnvData.IMAGES_PATH, "temp/")
}
func GetTempImagePath(fileName string) string {
	return filepath.Join(GetTempImageFolder(), AddImageExtensionIfNeeded(fileName))
}
