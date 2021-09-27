package handlers

import (
	"fmt"
	"net/http"

	"github.com/GMcD/telar-web/micros/collective/database"
	models "github.com/GMcD/telar-web/micros/collective/models"
	service "github.com/GMcD/telar-web/micros/collective/services"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	utils "github.com/red-gold/telar-core/utils"
)

type MembersPayload struct {
	Users map[string]interface{} `json:"users"`
}

// ReadDtoCollectiveHandle a function invocation
func ReadDtoCollectiveHandle(c *fiber.Ctx) error {

	collectiveId := c.Params("collectiveId")
	log.Info("Read dto collective by collectiveId %s", userId)
	userUUID, uuidErr := uuid.FromString(userId)
	if uuidErr != nil {
		log.Error("Parse UUID %s ", uuidErr.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("parseUUIDError", "Can not parse user id!"))
	}
	// Create service
	collectiveService, serviceErr := service.NewCollectiveService(database.Db)
	if serviceErr != nil {
		log.Error("NewCollectiveService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/collectiveService", "Error happened while creating collectiveService!"))
	}
	//
	foundCollective, err := collectiveService.FindByCollectiveId(userUUID)
	if err != nil {
		log.Error("FindByCollectiveId %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/findByCollectiveId", "Error happened while finding collective!"))
	}

	if foundCollective == nil {
		log.Error("Could not find collecitive " + collectiveUUID.String())
		return c.Status(http.StatusNotFound).JSON(utils.Error("notFoundCollective", "Error happened while finding collective!"))
	}

	return c.JSON(foundCollective)

}

// ReadProfileHandle a function invocation
func ReadCollectiveHandle(c *fiber.Ctx) error {

	collectiveId := c.Params("collectiveId")
	if collectiveId == "" {
		errorMessage := fmt.Sprintf("collective Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("collectiveIdRequired", errorMessage))
	}
	collectiveUUID, uuidErr := uuid.FromString(collectiveId)
	if uuidErr != nil {
		log.Error("Parse UUID %s ", uuidErr.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("parseUUIDError", "Can not parse collective id!"))
	}

	// Create service
	collectiveService, serviceErr := service.NewCollectiveService(database.Db)
	if serviceErr != nil {
		log.Error("NewCollectiveService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/collectiveService", "Error happened while creating collectiveService!"))
	}

	foundCollective, err := collectiveService.FindByCollectiveId(collectiveUUID)
	if err != nil {
		log.Error("FindBycCllectiveId %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/findByCollectiveId", "Error happened while finding collective!"))
	}

	if foundCollective == nil {
		log.Error("Could not find collective " + collectiveUUID.String())
		return c.Status(http.StatusNotFound).JSON(utils.Error("notFoundCollective", "Error happened while finding collective!"))
	}

	collectiveModel := models.MyCollectiveModel{
		ObjectId:      foundCollective.ObjectId,
		Name:          foundCollective.Name,
		Avatar:        foundCollective.Avatar,
		Banner:        foundCollective.Banner,
		FollowerCount: foundCollective.FollowerCount,
		PostCound:     foundCollective.PostCount,
		CreatedDate:   foundCollective.CreatedDate,
	}

	return c.JSON(collectiveModel)

}

// DispatchCollectiveHandle a function invocation to read authed user profile
func DispatchCollectiveHandle(c *fiber.Ctx) error {

	// Parse model object
	model := new(models.DispatchCollectiveModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Unmarshal  models.DispatchCollectiveModel array %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/dispatchCollectiveModelParser", "Error happened while parsing model!"))
	}

	if len(model.CollectiveIds) == 0 {
		errorMessage := "CollectiveIds are required!"
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("collectiveIdsRequired", errorMessage))
	}

	// Create service
	collectiveService, serviceErr := service.NewCollectiveService(database.Db)
	if serviceErr != nil {
		log.Error("NewCollectiveService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/collectiveService", "Error happened while creating collectiveService!"))
	}

	foundCollectives, err := collectiveService.FindCollectiveByUserIds(model.CollectiveIds)
	if err != nil {
		log.Error("FindCollectiveByCollectiveIds %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/findCollectiveByCollectiveIds", "Error happened while finding collectives!"))
	}

	mappedCollectives := make(map[string]interface{})
	for _, v := range foundCollectives {
		mappedCollective := make(map[string]interface{})
		mappedCollective["collectiveId"] = v.ObjectId
		mappedCollective["name"] = v.Name
		mappedCollective["avatar"] = v.Avatar
		mappedCollective["banner"] = v.Banner
		mappedCollective["lastSeen"] = v.LastSeen
		mappedCollective["createdDate"] = v.CreatedDate
		mappedCollective[v.ObjectId.String()] = mappedUser
	}

	actionRoomPayload := &MembersPayload{
		Collectives: mappedCollectives,
	}

	activeRoomAction := Action{
		Type:    "SET_USER_ENTITIES",
		Payload: actionRoomPayload,
	}

	currentUser, ok := c.Locals("user").(types.UserContext)
	if !ok {
		log.Warn("[DispatchCollectiveHandle] Can not get current user")
		currentUser = types.UserContext{}
	}
	log.Info("Current USER %v", currentUser)

	userInfoReq := &UserInfoInReq{
		UserId:      currentUser.UserID,
		Username:    currentUser.Username,
		Avatar:      currentUser.Avatar,
		DisplayName: currentUser.DisplayName,
		SystemRole:  currentUser.SystemRole,
	}

	go dispatchAction(activeRoomAction, userInfoReq)

	return c.SendStatus(http.StatusOK)

}

// GetCollectiveByIds a function invocation to profiles by ids
func GetCollectiveByIds(c *fiber.Ctx) error {

	// Parse model object
	model := new(models.GetCollectiveModel)
	if err := c.BodyParser(model); err != nil {
		errorMessage := fmt.Sprintf("Unmarshal  models.GetCollectiveModel array %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/getCollectiveModelParser", "Error happened while parsing model!"))
	}

	// Create service
	collectiveService, serviceErr := service.NewCollectiveService(database.Db)
	if serviceErr != nil {
		log.Error("NewCollectiveService %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/collectiveService", "Error happened while creating collectiveService!"))
	}

	foundCollective, err := collectiveService.FindCollectiveIds(model.CollectiveIds)
	if err != nil {
		log.Error("FindByCollectiveId %s", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/findByCollectiveId", "Error happened while finding collective!"))
	}

	return c.JSON(foundCollective)

}
