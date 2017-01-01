package status

import (
	"github.com/hellofresh/goengine"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Projection struct {
	coll *mgo.Collection
}

func NewProjection(session *mgo.Database) *Projection {
	return &Projection{session.C("statuses")}
}

func (s *Projection) StatusCreated(event *goengine.DomainMessage) error {
	var status ReadModel
	ev := event.Payload.(*StatusCreated)
	status.ID = bson.NewObjectId()
	status.AggregateID = ev.ID
	status.UserID = ev.UserID

	return s.coll.Insert(status)
}

func (s *Projection) StatusRemoved(event *goengine.DomainMessage) error {
	ev := event.Payload.(*StatusRemoved)
	return s.coll.Remove(bson.M{"aggergate_id": ev.ID})
}

func (s *Projection) StatusAllDay(event *goengine.DomainMessage) error {
	ev := event.Payload.(*AllDay)
	status, err := s.find(ev.ID)
	if nil != err {
		return err
	}
	status.IsAllDay = true

	return s.save(status)
}

func (s *Projection) StatusReasonChanged(event *goengine.DomainMessage) error {
	ev := event.Payload.(*StatusReasonChanged)

	status, err := s.find(ev.ID)
	if nil != err {
		return err
	}
	status.Reason = ev.Reason
	status.Message = ev.Message

	err = s.save(status)
	return err
}

func (s *Projection) StatusTimePeriodChanged(event *goengine.DomainMessage) error {
	ev := event.Payload.(*StatusTimePeriodChanged)

	status, err := s.find(ev.ID)
	if nil != err {
		return err
	}
	status.TimePeriod = ev.TimePeriod

	return s.save(status)
}

func (s *Projection) StatusTime(event *goengine.DomainMessage) error {
	ev := event.Payload.(*StatusTime)
	status, err := s.find(ev.ID)
	if nil != err {
		return err
	}
	status.StartsAt = ev.StartsAt
	status.EndsAt = ev.EndsAt

	return s.save(status)
}

func (s *Projection) find(id string) (*ReadModel, error) {
	var status ReadModel
	err := s.coll.Find(bson.M{"aggregate_id": id}).One(&status)
	return &status, err
}

func (s *Projection) save(status *ReadModel) error {
	return s.coll.Update(bson.M{"aggregate_id": status.AggregateID}, status)
}
