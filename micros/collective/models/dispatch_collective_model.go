package models

import "github.com/gofrs/uuid"

type DispatchCollectiveModel struct {
	CollectiveIds []uuid.UUID `json:"collectiveIds" bson:"collectiveIds"`
	ReqUserId     uuid.UUID   `json:"reqUserId" bson:"reqUserId"`
}
