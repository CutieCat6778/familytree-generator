package model

import (
	"time"
)

type Family struct {
	ID          string     `json:"id"`
	HusbandID   *string    `json:"husband_id,omitempty"`
	WifeID      *string    `json:"wife_id,omitempty"`
	ChildrenIDs []string   `json:"children_ids"`
	MarriedDate time.Time  `json:"married_date"`
	DivorceDate *time.Time `json:"divorce_date,omitempty"`
}

func NewFamily(id string, marriedDate time.Time) *Family {
	return &Family{
		ID:          id,
		ChildrenIDs: make([]string, 0),
		MarriedDate: marriedDate,
	}
}

func (f *Family) SetHusband(id string) {
	f.HusbandID = &id
}

func (f *Family) SetWife(id string) {
	f.WifeID = &id
}

func (f *Family) AddChild(id string) {
	f.ChildrenIDs = append(f.ChildrenIDs, id)
}

func (f *Family) IsDivorced() bool {
	return f.DivorceDate != nil
}

func (f *Family) ChildCount() int {
	return len(f.ChildrenIDs)
}
