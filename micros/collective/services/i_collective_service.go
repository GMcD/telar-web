package service

import (
	dto "github.com/GMcD/telar-web/micros/collective/dto"
	uuid "github.com/gofrs/uuid"
)

type CollectiveService interface {
	SaveCollective(collective *dto.Collective) error
	FindOneCollective(filter interface{}) (*dto.Collective, error)
	FindCollectiveList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Collective, error)
	QueryCollective(search string, sortBy string, page int64, notIncludeCollectiveIDList []uuid.UUID) ([]dto.Collective, error)
	FindCollectiveByCollectiveIds(collectiveIds []uuid.UUID) ([]dto.Collective, error)
	FindByCollectiveId(collectiveId uuid.UUID) (*dto.Collective, error)
	DeleteCollective(filter interface{}) error
	FindByCollectiveName(name string) (*dto.Collective, error)
	CreateCollectiveIndex(indexes map[string]interface{}) error
	IncreaseFollowCount(objectId uuid.UUID, inc int) error
	IncreaseFollowerCount(objectId uuid.UUID, inc int) error
	IncreasePostCount(collectiveUUID uuid.UUID, inc int) error
}
