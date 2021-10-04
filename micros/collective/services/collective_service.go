package service

import (
	"fmt"

	dto "github.com/GMcD/telar-web/micros/collective/dto"
	uuid "github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/config"
	coreData "github.com/red-gold/telar-core/data"
	"github.com/red-gold/telar-core/data/mongodb"
	mongoRepo "github.com/red-gold/telar-core/data/mongodb"
	"github.com/red-gold/telar-core/utils"
)

// CollectiveService handlers with injected dependencies
type CollectiveServiceImpl struct {
	CollectiveRepo coreData.Repository
}

// NewCollectiveService initializes CollectiveService's dependencies and create new CollectiveService struct
func NewCollectiveService(db interface{}) (CollectiveService, error) {

	collectiveService := &CollectiveServiceImpl{}

	switch *config.AppConfig.DBType {
	case config.DB_MONGO:

		mongodb := db.(mongodb.MongoDatabase)
		collectiveService.CollectiveRepo = mongoRepo.NewDataRepositoryMongo(mongodb)

	}
	if collectiveService.CollectiveRepo == nil {
		fmt.Printf("collectiveService.CollectiveRepo is nil! \n")
	}
	return collectiveService, nil
}

// SaveCollective save ollective informaition
func (s CollectiveServiceImpl) SaveCollective(collective *dto.Collective) error {

	if collective.ObjectId == uuid.Nil {
		var uuidErr error
		collective.ObjectId, uuidErr = uuid.NewV4()
		if uuidErr != nil {
			return uuidErr
		}
	}

	if collective.CreatedDate == 0 {
		collective.CreatedDate = utils.UTCNowUnix()
	}

	result := <-s.CollectiveRepo.Save(collectiveCollectionName, collective)

	return result.Error
}

// UpdateCollective update collective information
func (s CollectiveServiceImpl) UpdateCollective(filter interface{}, data interface{}) error {

	result := <-s.CollectiveRepo.Update(collectiveCollectionName, filter, data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindOneCollective get one collective informaition
func (s CollectiveServiceImpl) FindOneCollective(filter interface{}) (*dto.Collective, error) {

	result := <-s.CollectiveRepo.FindOne(collectiveCollectionName, filter)
	if result.Error() != nil {
		if result.Error() == coreData.ErrNoDocuments {
			return nil, nil
		}
		return nil, result.Error()
	}

	var collectiveResult dto.Collective
	errDecode := result.Decode(&collectiveResult)
	if errDecode != nil {
		return nil, fmt.Errorf("Error docoding on dto.Collective")
	}
	return &collectiveResult, nil
}

// FindCollectiveList get all Collective informaition
func (s CollectiveServiceImpl) FindCollectiveList(filter interface{}, limit int64, skip int64, sort map[string]int) ([]dto.Collective, error) {

	result := <-s.CollectiveRepo.Find("collective", filter, limit, skip, sort)
	defer result.Close()
	if result.Error() != nil {
		return nil, result.Error()
	}
	var collectiveList []dto.Collective
	for result.Next() {
		var collective dto.Collective
		errDecode := result.Decode(&collective)
		if errDecode != nil {
			return nil, fmt.Errorf("Error docoding on dto.Collective")
		}
		collectiveList = append(collectiveList, collective)
	}

	return collectiveList, nil
}

// QueryPost get all user profile by query
func (s CollectiveServiceImpl) QueryCollective(search string, sortBy string, page int64, notIncludeCollectiveIDList []uuid.UUID) ([]dto.Collective, error) {
	sortMap := make(map[string]int)
	sortMap[sortBy] = -1
	skip := numberOfItems * (page - 1)
	limit := numberOfItems
	filter := make(map[string]interface{})
	if search != "" {
		filter["$text"] = coreData.SearchOperator{Search: search}
	}
	if notIncludeCollectiveIDList != nil && len(notIncludeCollectiveIDList) > 0 {
		nin := make(map[string]interface{})
		nin["$nin"] = notIncludeCollectiveIDList
		filter["objectId"] = nin
	}
	fmt.Println(filter)
	result, err := s.FindCollectiveList(filter, limit, skip, sortMap)

	return result, err
}

// FindCollectiveByCollectiveIds Find profile by Collective IDs
func (s CollectiveServiceImpl) FindCollectiveByCollectiveIds(collectiveIds []uuid.UUID) ([]dto.Collective, error) {
	sortMap := make(map[string]int)
	sortMap["createdDate"] = -1

	include := make(map[string]interface{})
	include["$in"] = collectiveIds

	filter := make(map[string]interface{})
	filter["objectId"] = include

	result, err := s.FindCollectiveList(filter, 0, 0, sortMap)

	return result, err
}

// FindByCollectivename find user Collective by name
func (s CollectiveServiceImpl) FindByCollectiveName(collectiveName string) (*dto.Collective, error) {

	filter := struct {
		Name string `json:"name"`
	}{
		Name: collectiveName,
	}
	return s.FindOneCollective(filter)
}

// FindByCollectiveId find Collective by collectiveId
func (s CollectiveServiceImpl) FindByCollectiveId(collectiveId uuid.UUID) (*dto.Collective, error) {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: collectiveId,
	}
	return s.FindOneCollective(filter)
}

// DeleteCollective get all Collective information.
func (s CollectiveServiceImpl) DeleteCollective(filter interface{}) error {

	result := <-s.CollectiveRepo.Delete(collectiveCollectionName, filter, true)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteManyCollective get all Collective informaition.
func (s CollectiveServiceImpl) DeleteManyCollective(filter interface{}) error {

	result := <-s.CollectiveRepo.Delete(collectiveCollectionName, filter, false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateCollectiveIndex create index for Collective search.
func (s CollectiveServiceImpl) CreateCollectiveIndex(indexes map[string]interface{}) error {
	result := <-s.CollectiveRepo.CreateIndex(collectiveCollectionName, indexes)
	return result
}

// Increment increment a profile field
func (s CollectiveServiceImpl) Increment(objectId uuid.UUID, field string, value int) error {

	filter := struct {
		ObjectId uuid.UUID `json:"objectId" bson:"objectId"`
	}{
		ObjectId: objectId,
	}

	data := make(map[string]interface{})
	data[field] = value

	incOperator := coreData.IncrementOperator{
		Inc: data,
	}
	return s.UpdateCollective(filter, incOperator)
}

// IncreaseFollowCount increment follow count of post
func (s CollectiveServiceImpl) IncreaseFollowCount(objectId uuid.UUID, inc int) error {
	return s.Increment(objectId, "followCount", inc)
}

// IncreaseFollowerCount increment follower count of post
func (s CollectiveServiceImpl) IncreaseFollowerCount(objectId uuid.UUID, inc int) error {
	return s.Increment(objectId, "followerCount", inc)
}

// IncreasePostCount increment post count of collective
func (s CollectiveServiceImpl) IncreasePostCount(objectId uuid.UUID, inc int) error {
	return s.Increment(objectId, "postCount", inc)
}
