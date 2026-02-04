package model

import (
	"time"
)


type Gender string

const (
	Male   Gender = "M"
	Female Gender = "F"
)


type EducationLevel string

const (
	NoEducation EducationLevel = "none"
	Primary     EducationLevel = "primary"
	Secondary   EducationLevel = "secondary"
	Tertiary    EducationLevel = "tertiary"
)


type EmploymentStatus string

const (
	Employed   EmploymentStatus = "employed"
	Unemployed EmploymentStatus = "unemployed"
	Retired    EmploymentStatus = "retired"
	Student    EmploymentStatus = "student"
	Child      EmploymentStatus = "child"
)

type ResidenceType string

const (
	Urban ResidenceType = "urban"
	Rural ResidenceType = "rural"
)


type HealthProfile struct {
	AlcoholConsumption float64 `json:"alcohol_consumption"` 
	TobaccoUse         bool    `json:"tobacco_use"`
}


type MaritalStatus string

const (
	Single       MaritalStatus = "single"
	Married      MaritalStatus = "married"
	Divorced     MaritalStatus = "divorced"
	Widowed      MaritalStatus = "widowed"
	Remarried    MaritalStatus = "remarried"
)


type Person struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    Gender `json:"gender"`

	BirthDate      time.Time  `json:"birth_date"`
	DeathDate      *time.Time `json:"death_date,omitempty"`
	BirthCountry   string     `json:"birth_country"`
	CurrentCountry string     `json:"current_country"`

	
	FatherID    *string  `json:"father_id,omitempty"`
	MotherID    *string  `json:"mother_id,omitempty"`
	SpouseIDs   []string `json:"spouse_ids,omitempty"`
	ChildrenIDs []string `json:"children_ids,omitempty"`

	
	Education  EducationLevel   `json:"education"`
	Employment EmploymentStatus `json:"employment"`
	Health     HealthProfile    `json:"health"`
	Underweight bool            `json:"underweight,omitempty"`
	Residence  ResidenceType    `json:"residence,omitempty"`
	GDPPerCapita float64        `json:"gdp_per_capita,omitempty"`
	WealthIndex  float64        `json:"wealth_index,omitempty"`
	FamilyWealth float64        `json:"family_wealth,omitempty"`
	IsRich       bool           `json:"is_rich,omitempty"`

	
	MaritalStatus    MaritalStatus `json:"marital_status"`
	MarriageAge      int           `json:"marriage_age,omitempty"`
	NumberOfChildren int           `json:"number_of_children"`
	IsSingleParent   bool          `json:"is_single_parent,omitempty"`
	BornOutsideMarriage bool       `json:"born_outside_marriage,omitempty"`

	
	Events []LifeEvent `json:"events,omitempty"`

	
	Generation int `json:"generation"`
}


func (p *Person) IsAlive() bool {
	return p.DeathDate == nil
}


func (p *Person) Age(at time.Time) int {
	endDate := at
	if p.DeathDate != nil && p.DeathDate.Before(at) {
		endDate = *p.DeathDate
	}
	return yearsBetween(p.BirthDate, endDate)
}


func (p *Person) AgeAtDeath() int {
	if p.DeathDate == nil {
		return -1
	}
	return yearsBetween(p.BirthDate, *p.DeathDate)
}


func (p *Person) FullName() string {
	return p.FirstName + " " + p.LastName
}


func yearsBetween(start, end time.Time) int {
	years := end.Year() - start.Year()
	if end.YearDay() < start.YearDay() {
		years--
	}
	return years
}


func NewPerson(id, firstName, lastName string, gender Gender, birthDate time.Time, country string, generation int) *Person {
	return &Person{
		ID:             id,
		FirstName:      firstName,
		LastName:       lastName,
		Gender:         gender,
		BirthDate:      birthDate,
		BirthCountry:   country,
		CurrentCountry: country,
		SpouseIDs:      make([]string, 0),
		ChildrenIDs:    make([]string, 0),
		Events:         make([]LifeEvent, 0),
		Generation:     generation,
	}
}
