package service

import (
	dto "github.com/GMcD/telar-web/micros/notifications/dto"
	uuid "github.com/gofrs/uuid"
	coreData "github.com/red-gold/telar-core/data"
)

type NotificationService interface {
	SaveNotification(notification *dto.Notification) error
	FindOneNotification(filter interface{}) (*dto.Notification, error)
	FindNotificationList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Notification, error)
	FindNotificationsReceiver(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Notification, error)
	FindById(objectId uuid.UUID) (*dto.Notification, error)
	FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.Notification, error)
	UpdateNotification(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error
	UpdateNotificationById(data *dto.Notification) error
	UpdateBulkNotificationList(userNotification []dto.Notification) error
	UpdateEmailSent(notifyIds []uuid.UUID) error
	DeleteNotification(filter interface{}) error
	DeleteNotificationByOwner(notificationReceiverId uuid.UUID, notificationId uuid.UUID) error
	DeleteManyNotifications(filter interface{}) error
	CreateNotificationIndex(indexes map[string]interface{}) error
	GetNotificationByUserId(userId *uuid.UUID, sortBy string, page int64, limit int64) ([]dto.Notification, error)
	GetLastNotifications() ([]dto.Notification, error)
	SeenNotification(objectId uuid.UUID, userId uuid.UUID) error
	SeenAllNotifications(userId uuid.UUID) error
	DeleteNotificationsByUserId(userId uuid.UUID) error
}
