package model

import (
	"time"
)

// Family represents a nuclear family unit
type Family struct {
	ID          string    `json:"id"`
	HusbandID   *string   `json:"husband_id,omitempty"`
	WifeID      *string   `json:"wife_id,omitempty"`
	ChildrenIDs []string  `json:"children_ids"`
	MarriedDate time.Time `json:"married_date"`
	DivorceDate *time.Time `json:"divorce_date,omitempty"`
}

// NewFamily creates a new family
func NewFamily(id string, marriedDate time.Time) *Family {
	return &Family{
		ID:          id,
		ChildrenIDs: make([]string, 0),
		MarriedDate: marriedDate,
	}
}

// SetHusband sets the husband of the family
func (f *Family) SetHusband(id string) {
	f.HusbandID = &id
}

// SetWife sets the wife of the family
func (f *Family) SetWife(id string) {
	f.WifeID = &id
}

// AddChild adds a child to the family
func (f *Family) AddChild(id string) {
	f.ChildrenIDs = append(f.ChildrenIDs, id)
}

// IsDivorced returns true if the family is divorced
func (f *Family) IsDivorced() bool {
	return f.DivorceDate != nil
}

// ChildCount returns the number of children
func (f *Family) ChildCount() int {
	return len(f.ChildrenIDs)
}
