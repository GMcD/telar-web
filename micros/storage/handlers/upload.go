package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	appConfig "github.com/GMcD/telar-web/micros/storage/config"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"

	coreSetting "github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
)

// UploadeHandle a function invocation
func UploadeHandle(c *fiber.Ctx) error {

	currentUser, ok := c.Locals("user").(types.UserContext)
	if !ok {
		log.Error("[UploadeHandle] Can not get current user")
		return c.Status(http.StatusUnauthorized).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	storageConfig := &appConfig.StorageConfig
	log.Info("Hit upload endpoint by userId : %v", currentUser.UserID)

	// params from /storage/:uid/:dir
	dirName := c.Params("dir")
	if dirName == "" {
		errorMessage := fmt.Sprintf("Directory name is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("directoryNameRequired", "Directory name is required!"))
	}

	log.Info("Directory name: %s", dirName)

	// FormFile returns the first file for the given key `file`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file

	file, err := c.FormFile("file")
	if err != nil {
		log.Error("Error Retrieving the File %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/upload", "Error Retrieving the File!"))

	}

	log.Info("Uploaded File: %+v", file.Filename)
	log.Info("File Size: %+v", file.Size)
	log.Info("MIME Header: %+v", file.Header)

	extension := filepath.Ext(file.Filename)
	fileNameUUID := uuid.Must(uuid.NewV4())

	fileName := fileNameUUID.String()
	fileNameWithExtension := fmt.Sprintf("%s%s", fileName, extension)

	objectName := fmt.Sprintf("%s/%s/%s", currentUser.UserID.String(), dirName, fileNameWithExtension)
	log.Info("Object Name: %s", objectName)

	s3Path, uploadErr := UploadImage(c, file, objectName)
	if uploadErr != nil {
		log.Error("Copy file to storage error %s", uploadErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/upload", "Copy file to storage error!"))
	}
	log.Info("Object Path: %s", s3Path)

	prettyURL := utils.GetPrettyURLf(storageConfig.BaseRoute)
	downloadURL := fmt.Sprintf("%s/%s/%s/%s", *coreSetting.AppConfig.Gateway+prettyURL,
		currentUser.UserID, dirName, fileNameWithExtension)
	log.Info("Object URL: %s", downloadURL)

	return c.JSON(fiber.Map{
		"payload": s3Path,
		"path":    downloadURL,
	})

}

// config := &firebase.Config{
// 	StorageBucket: storageConfig.BucketName,
// }

// opt := option.WithCredentialsJSON([]byte(storageConfig.StorageSecret))
// app, err := firebase.NewApp(ctx, config, opt)
// if err != nil {
// 	log.Error("Credential parse %s", err.Error())
// 	return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/upload", "Credential parse error!"))
// }

// client, err := app.Storage(ctx)
// if err != nil {
// 	log.Error("Get storage client %s", err.Error())
// 	return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/upload", "Get storage client!"))
// }

// bucket, err := client.DefaultBucket()
// if err != nil {
// 	log.Error("Get default bucket %s", err.Error())
// 	return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/upload", "Get default bucket!"))
// }

// wc := bucket.Object(objectName).NewWriter(ctx)

// multiFile, openFileErr := file.Open()
// if openFileErr != nil {
// 	log.Error("Open file error %s", openFileErr.Error())
// 	return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/upload", "Open file error!"))
// }

// if _, err = io.Copy(wc, multiFile); err != nil {
// 	log.Error("Copy file to storage error %s", err.Error())
// 	return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/upload", "Copy file to storage error!"))

// }
// if err := wc.Close(); err != nil {
// 	log.Error("Close storage writer error %s", err.Error())
// 	return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/upload", "Close storage writer error!"))
// }
