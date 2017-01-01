package status

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type ReadModel struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	AggregateID string `bson:"aggregate_id" json:"aggregate_id"`
	UserID      string `bson:"user_id" json:"user_id"`
	Reason      Reason `bson:"reason" json:"reason"`
	IsAllDay    bool `bson:"is_all_day" json:"is_all_day"`
	StartsAt    time.Time `bson:"starts_at" json:"starts_at"`
	EndsAt      time.Time `bson:"ends_at" json:"ends_at"`
	Message     string `bson:"message" json:"message"`
	TimePeriod  TimePeriod `bson:"time_period" json:"time_period"`
}
