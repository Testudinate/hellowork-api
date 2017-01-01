package status

import (
	"github.com/hellofresh/goengine"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type WriteRepository interface {
	Add(aggregate goengine.AggregateRoot) error
	Reconstitute(id string, aggregate goengine.AggregateRoot) error
	NextIdentity() string
}

type MongoDBWriteRepository struct {
	streamName    goengine.StreamName
	aggregateRepo goengine.AggregateRepository
}

// NewMongoDBWriteRepository creates a new instance of MongoDBWriteRepository
func NewMongoDBWriteRepository(repo goengine.AggregateRepository) (*MongoDBWriteRepository, error) {
	return &MongoDBWriteRepository{goengine.StreamName("statuses"), repo}, nil
}

func (r *MongoDBWriteRepository) Add(aggregate goengine.AggregateRoot) error {
	return r.aggregateRepo.Save(aggregate, r.streamName)
}

func (r *MongoDBWriteRepository) Reconstitute(id string, aggregate goengine.AggregateRoot) error {
	return r.aggregateRepo.Reconstitute(id, aggregate, r.streamName)
}

func (r *MongoDBWriteRepository) NextIdentity() string {
	return uuid.NewV4().String()
}

type ReadStatusRepository interface {
	Find(id string) (*ReadModel, error)
	FindAllByUserID(userID string) ([]*ReadModel, error)
	FindActiveStatus(userID string) (*ReadModel, error)
}

type MongoDBReadRepository struct {
	coll *mgo.Collection
}

// NewMongoDBReadRepository creates a new instance of MongoDBReadRepository
func NewMongoDBReadRepository(db *mgo.Database) (*MongoDBReadRepository, error) {
	coll := db.C("statuses")

	return &MongoDBReadRepository{coll}, nil
}

func (r *MongoDBReadRepository) FindAllByUserID(userID string) ([]*ReadModel, error) {
	var statuses []*ReadModel
	err := r.coll.Find(bson.M{"user_id": userID}).All(&statuses)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func (r *MongoDBReadRepository) FindActiveStatus(userID string) (*ReadModel, error) {
	var status *ReadModel
	err := r.coll.Find(bson.M{"user_id": userID, "is_active": true}).One(&status)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (r *MongoDBReadRepository) Find(id string) (*ReadModel, error) {
	var status *ReadModel
	err := r.coll.Find(bson.M{"aggregate_id": id}).One(&status)
	if err != nil {
		return nil, err
	}

	return status, nil
}
