package handlers

import (
	"fmt"
	"net/http"

	"github.com/GMcD/telar-web/micros/notifications/database"
	domain "github.com/GMcD/telar-web/micros/notifications/dto"
	models "github.com/GMcD/telar-web/micros/notifications/models"
	service "github.com/GMcD/telar-web/micros/notifications/services"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
)

// UpdateNotificationHandle handle update a notification
func UpdateNotificationHandle(c *fiber.Ctx) error {

	// Create the model object
	model := new(models.NotificationModel)
	if err := c.BodyParser(model); err != nil {
		log.Error("[UpdateNotificationHandle.parse.NotificationModel] %s", err.Error())
		return c.Status(http.StatusBadRequest).JSON(utils.Error("parseModel", "Error happened while parsing model!"))
	}

	// Create service
	notificationService, serviceErr := service.NewNotificationService(database.Db)
	if serviceErr != nil {
		log.Error("[UpdateNotificationHandle.NewNotificationService] %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/notificationService", "Error happened while creating notificationService!"))
	}

	currentUser, ok := c.Locals("user").(types.UserContext)
	if !ok {
		log.Error("[UpdateNotificationHandle] Can not get current user")
		return c.Status(http.StatusUnauthorized).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	updatedNotification := &domain.Notification{
		ObjectId:             model.ObjectId,
		OwnerUserId:          currentUser.UserID,
		OwnerDisplayName:     currentUser.DisplayName,
		OwnerAvatar:          currentUser.Avatar,
		Title:                model.Title,
		Description:          model.Description,
		URL:                  model.URL,
		NotifyRecieverUserId: model.NotifyRecieverUserId,
		TargetId:             model.TargetId,
		IsSeen:               model.IsSeen,
		Type:                 model.Type,
		EmailNotification:    model.EmailNotification,
	}

	if err := notificationService.UpdateNotificationById(updatedNotification); err != nil {
		errorMessage := fmt.Sprintf("Update Notification Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updateNotificationById", "Error happened while updating notification!"))
	}

	return c.SendStatus(http.StatusOK)

}

// SeenNotificationHandle handle set notification seen
func SeenNotificationHandle(c *fiber.Ctx) error {

	// params from /notifications/seen/:notificationId
	notificationId := c.Params("notificationId")
	if notificationId == "" {
		errorMessage := fmt.Sprintf("Notification Id is required!")
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("notificationIdRequired", "Notification id is required!"))
	}

	notificationUUID, uuidErr := uuid.FromString(notificationId)
	if uuidErr != nil {
		errorMessage := fmt.Sprintf("UUID Error %s", uuidErr.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusBadRequest).JSON(utils.Error("parseUUID", "Can not parse UUID!"))
	}
	// Create service
	notificationService, serviceErr := service.NewNotificationService(database.Db)
	if serviceErr != nil {
		log.Error("[SeenNotificationHandle.NewNotificationService] %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/notificationService", "Error happened while creating notificationService!"))
	}

	currentUser, ok := c.Locals("user").(types.UserContext)
	if !ok {
		log.Error("[SeenNotificationHandle] Can not get current user")
		return c.Status(http.StatusUnauthorized).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := notificationService.SeenNotification(notificationUUID, currentUser.UserID); err != nil {
		errorMessage := fmt.Sprintf("Update Notification Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updateNotification", "Can not update notification!"))
	}

	return c.SendStatus(http.StatusOK)

}

// SeenAllNotificationsHandle handle set all notifications seen
func SeenAllNotificationsHandle(c *fiber.Ctx) error {

	// Create service
	notificationService, serviceErr := service.NewNotificationService(database.Db)
	if serviceErr != nil {
		log.Error("[SeenAllNotificationHandle.NewNotificationService] %s", serviceErr.Error())
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/notificationService", "Error happened while creating notificationService!"))
	}

	currentUser, ok := c.Locals("user").(types.UserContext)
	if !ok {
		log.Error("[SeenAllNotificationHandle] Can not get current user")
		return c.Status(http.StatusUnauthorized).JSON(utils.Error("invalidCurrentUser",
			"Can not get current user"))
	}

	if err := notificationService.SeenAllNotifications(currentUser.UserID); err != nil {
		errorMessage := fmt.Sprintf("Update Notification Error %s", err.Error())
		log.Error(errorMessage)
		return c.Status(http.StatusInternalServerError).JSON(utils.Error("internal/updateNotification", "Can not update notification!"))
	}

	return c.SendStatus(http.StatusOK)

}
