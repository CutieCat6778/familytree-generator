package model

import (
	"time"
)

// EventType represents the type of life event
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

// LifeEvent represents a significant event in a person's life
type LifeEvent struct {
	Type        EventType `json:"type"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location,omitempty"`
	Description string    `json:"description,omitempty"`
	RelatedID   string    `json:"related_id,omitempty"` // e.g., spouse ID for marriage
}

// NewLifeEvent creates a new life event
func NewLifeEvent(eventType EventType, date time.Time, location string) LifeEvent {
	return LifeEvent{
		Type:     eventType,
		Date:     date,
		Location: location,
	}
}

// WithDescription adds a description to the event
func (e LifeEvent) WithDescription(desc string) LifeEvent {
	e.Description = desc
	return e
}

// WithRelatedID adds a related person ID to the event
func (e LifeEvent) WithRelatedID(id string) LifeEvent {
	e.RelatedID = id
	return e
}
