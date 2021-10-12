package service

import (
	dto "github.com/GMcD/telar-web/micros/profile/dto"
	uuid "github.com/gofrs/uuid"
)

type UserProfileService interface {
	SaveUserProfile(userProfile *dto.UserPublicProfile) error
	FindOneUserProfile(filter interface{}) (*dto.UserPublicProfile, error)
	FindUserProfileList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.UserPublicProfile, error)
	QueryUserProfile(search string, sortBy string, page int64, notIncludeUserIDList []uuid.UUID) ([]dto.UserPublicProfile, error)
	FindProfileByUserIds(userIds []uuid.UUID) ([]dto.UserPublicProfile, error)
	FindByUserId(userId uuid.UUID) (*dto.UserPublicProfile, error)
	FindBySocialName(socialName string) (*dto.UserPublicProfile, error)
	UpdateUserProfile(filter interface{}, data interface{}) error
	UpdateLastSeenNow(userId uuid.UUID) error
	UpdateUserProfileById(userId uuid.UUID, data interface{}) error
	DeleteUserProfile(filter interface{}) error
	DeleteManyUserProfile(filter interface{}) error
	FindByUsername(username string) (*dto.UserPublicProfile, error)
	CreateUserProfileIndex(indexes map[string]interface{}) error
	IncreaseFollowCount(objectId uuid.UUID, inc int) error
	IncreaseFollowerCount(objectId uuid.UUID, inc int) error
}
