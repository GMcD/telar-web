package models

import "github.com/gofrs/uuid"

type GetCollectivesModel struct {
	CollectiveIds []uuid.UUID `json:"CollectiveIds" bson:"CollectiveIds"`
}
