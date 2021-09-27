package models

import "github.com/gofrs/uuid"

type GetCollectiveModel struct {
	CollectiveIds []uuid.UUID `json:"CollectiveIds" bson:"CollectiveIds"`
}
