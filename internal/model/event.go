package model

import (
	"time"
)


type EventType string

const (
	EventBirth      EventType = "birth"
	EventDeath      EventType = "death"
	EventMarriage   EventType = "marriage"
	EventDivorce    EventType = "divorce"
	EventMigration  EventType = "migration"
	EventGraduation EventType = "graduation"
	EventRetirement EventType = "retirement"
)


type LifeEvent struct {
	Type        EventType `json:"type"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location,omitempty"`
	Description string    `json:"description,omitempty"`
	RelatedID   string    `json:"related_id,omitempty"` 
}


func NewLifeEvent(eventType EventType, date time.Time, location string) LifeEvent {
	return LifeEvent{
		Type:     eventType,
		Date:     date,
		Location: location,
	}
}


func (e LifeEvent) WithDescription(desc string) LifeEvent {
	e.Description = desc
	return e
}


func (e LifeEvent) WithRelatedID(id string) LifeEvent {
	e.RelatedID = id
	return e
}
