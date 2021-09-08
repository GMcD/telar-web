package service

import (
	dto "github.com/GMcD/telar-web/micros/actions/dto"
	uuid "github.com/gofrs/uuid"
	coreData "github.com/red-gold/telar-core/data"
)

type ActionRoomService interface {
	SaveActionRoom(actionRoom *dto.ActionRoom) error
	FindOneActionRoom(filter interface{}) (*dto.ActionRoom, error)
	FindActionRoomList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.ActionRoom, error)
	FindById(objectId uuid.UUID) (*dto.ActionRoom, error)
	FindByOwnerUserId(ownerUserId uuid.UUID) ([]dto.ActionRoom, error)
	UpdateActionRoom(filter interface{}, data interface{}, opts ...*coreData.UpdateOptions) error
	UpdateActionRoomById(data *dto.ActionRoom) error
	DeleteActionRoom(filter interface{}) error
	DeleteActionRoomByOwner(ownerUserId uuid.UUID, actionRoomId uuid.UUID) error
	DeleteManyActionRooms(filter interface{}) error
	CreateActionRoomIndex(indexes map[string]interface{}) error
	SetAccessKey(ownerUserId uuid.UUID) (string, error)
	VerifyAccessKey(ownerUserId uuid.UUID, accessKey string) (bool, error)
	GetAccessKey(ownerUserId uuid.UUID) (string, error)
}
