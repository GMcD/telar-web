package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/GMcD/telar-web/micros/collective/database"
	models "github.com/GMcD/telar-web/micros/collective/models"
	service "github.com/GMcD/telar-web/micros/collective/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	utils "github.com/red-gold/telar-core/utils"
)

// IncreaseFollowCount a function invocation
func IncreasePostCount(c *fiber.Ctx) error {

	// params from /follow/inc/:inc/:userId
	collectiveId := c.Params("collectiveId")
	if collectiveId == "" {
		errorMessage := fmt.Sprintf("Collective Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("collectiveIdRequired", errorMessage))
	}

	collectiveUUID, uuidErr := uuid.FromString(collectiveId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("collectiveIdIsNotValid", "Post id is not valid!"))
	}

	incParam := c.Params("inc")
	inc, err := strconv.Atoi(incParam)
	if err != nil {
		log.Error("Wrong inc param %s - %s", incParam, err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidIncParam", "Wrong inc param!"))

	}

	// Create service
	collectiveService, serviceErr := service.NewCollectiveService(database.Db)
	if serviceErr != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.Error("internal/collectiveService", "Internal error happened while creating collectiveService!"))
	}

	err = collectiveService.IncreasePostCount(collectiveUUID, inc)
	if err != nil {
		errorMessage := fmt.Sprintf("Update post count %s",
			err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("updatePostCount", "Error happened while updating post count!"))

	}

	return c.SendStatus(http.StatusOK)

}

// IncreaseFollowerCount a function invocation
func IncreaseFollowerCount(c *fiber.Ctx) error {

	// params from /follower/inc/:inc/:userId
	collectiveId := c.Params("collectiveId")
	if collectiveId == "" {
		errorMessage := fmt.Sprintf("collective Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("collectiveIdRequired", errorMessage))
	}

	collectiveUUID, uuidErr := uuid.FromString(collectiveId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("collectiveIdIsNotValid", "Post id is not valid!"))
	}

	incParam := c.Params("inc")
	inc, err := strconv.Atoi(incParam)
	if err != nil {
		log.Error("Wrong inc param %s - %s", incParam, err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("invalidIncParam", "Wrong inc param!"))

	}

	// Create service
	collectiveService, serviceErr := service.NewCollectiveService(database.Db)
	if serviceErr != nil {
		return c.Status(http.StatusBadRequest).JSON(utils.Error("internal/collectiveService", "Internal error happened while creating collectiveService!"))
	}

	err = collectiveService.IncreaseFollowerCount(collectiveUUID, inc)
	if err != nil {
		errorMessage := fmt.Sprintf("Update follower count %s",
			err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("updateFollowerCount", "Error happened while updating follower count!"))

	}

	return c.SendStatus(http.StatusOK)

}
