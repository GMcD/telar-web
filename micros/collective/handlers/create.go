package handlers

import (
	"fmt"
	"net/http"

	"github.com/GMcD/telar-web/micros/collective/database"
	"github.com/GMcD/telar-web/micros/collective/dto"
	service "github.com/GMcD/telar-web/micros/collective/services"
	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/pkg/log"
	utils "github.com/red-gold/telar-core/utils"
)

// InitPCollectiveIndexHandle handle create a new index
func InitCollectiveIndexHandle(c *fiber.Ctx) error {

	// Create service
	collectiveService, serviceErr := service.NewCollectiveService(database.Db)
	if serviceErr != nil {
		errorMessage := fmt.Sprintf("Collective service Error %s", serviceErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/collectiveService", "Error happened while creating collectiveService!"))
	}

	postIndexMap := make(map[string]interface{})
	postIndexMap["Name"] = "text"
	postIndexMap["objectId"] = 1
	if err := collectiveService.CreateCollectiveIndex(postIndexMap); err != nil {
		errorMessage := fmt.Sprintf("Create post index Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("createPostIndex", "Error happened while creating post index!"))
	}

	return c.SendStatus(http.StatusOK)

}

// CreateCollectiveHandle handle create a new collective
func CreateDtoCollectiveHandle(c *fiber.Ctx) error {

	// Create service
	collectiveService, serviceErr := service.NewCollectiveService(database.Db)
	if serviceErr != nil {
		errorMessage := fmt.Sprintf("Collective service Error %s", serviceErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/collectiveService", "Error happened while creating collectiveService!"))
	}

	model := new(dto.Collective)
	err := c.BodyParser(model)
	if err != nil {
		errorMessage := fmt.Sprintf("parse collective model %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("parseCollectiveModel", "Error happened while parsing model!"))

	}
	if err = collectiveService.SaveCollective(model); err != nil {
		errorMessage := fmt.Sprintf("Create collective error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("createCollectiveError", "Error happened while saving collective!"))
	}

	return c.SendStatus(http.StatusOK)
}

func CreateCollectiveHandle(c *fiber.Ctx) error {

	// Create service
	collectiveService, serviceErr := service.NewCollectiveService(database.Db)
	if serviceErr != nil {
		errorMessage := fmt.Sprintf("Collective service Error %s", serviceErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/collectiveService", "Error happened while creating collectiveService!"))
	}

	model := new(dto.Collective)
	err := c.BodyParser(model)
	if err != nil {
		errorMessage := fmt.Sprintf("parse collective model %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("parseCollectiveModel", "Error happened while parsing model!"))
	}

	if err = collectiveService.SaveCollective(model); err != nil {
		errorMessage := fmt.Sprintf("Create collective error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("createCollectiveError", "Error happened while saving collective!"))
	}
	return c.SendStatus(http.StatusOK)

}
