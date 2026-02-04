package model

import (
	"time"
)

// Gender represents biological sex
type Gender string

const (
	Male   Gender = "M"
	Female Gender = "F"
)

// EducationLevel represents highest education attained
type EducationLevel string

const (
	NoEducation EducationLevel = "none"
	Primary     EducationLevel = "primary"
	Secondary   EducationLevel = "secondary"
	Tertiary    EducationLevel = "tertiary"
)

// EmploymentStatus represents current employment state
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

// HealthProfile contains health-related attributes
type HealthProfile struct {
	AlcoholConsumption float64 `json:"alcohol_consumption"` // liters per year
	TobaccoUse         bool    `json:"tobacco_use"`
}

// MaritalStatus represents current marital state
type MaritalStatus string

const (
	Single       MaritalStatus = "single"
	Married      MaritalStatus = "married"
	Divorced     MaritalStatus = "divorced"
	Widowed      MaritalStatus = "widowed"
	Remarried    MaritalStatus = "remarried"
)

// Person represents an individual in the family tree
type Person struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    Gender `json:"gender"`

	BirthDate      time.Time  `json:"birth_date"`
	DeathDate      *time.Time `json:"death_date,omitempty"`
	BirthCountry   string     `json:"birth_country"`
	CurrentCountry string     `json:"current_country"`

	// Family relationships (store IDs for serialization)
	FatherID    *string  `json:"father_id,omitempty"`
	MotherID    *string  `json:"mother_id,omitempty"`
	SpouseIDs   []string `json:"spouse_ids,omitempty"`
	ChildrenIDs []string `json:"children_ids,omitempty"`

	// Attributes influenced by statistics
	Education  EducationLevel   `json:"education"`
	Employment EmploymentStatus `json:"employment"`
	Health     HealthProfile    `json:"health"`
	Underweight bool            `json:"underweight,omitempty"`
	Residence  ResidenceType    `json:"residence,omitempty"`
	GDPPerCapita float64        `json:"gdp_per_capita,omitempty"`
	WealthIndex  float64        `json:"wealth_index,omitempty"`
	FamilyWealth float64        `json:"family_wealth,omitempty"`
	IsRich       bool           `json:"is_rich,omitempty"`

	// Marriage details
	MaritalStatus    MaritalStatus `json:"marital_status"`
	MarriageAge      int           `json:"marriage_age,omitempty"`
	NumberOfChildren int           `json:"number_of_children"`
	IsSingleParent   bool          `json:"is_single_parent,omitempty"`
	BornOutsideMarriage bool       `json:"born_outside_marriage,omitempty"`

	// Life events
	Events []LifeEvent `json:"events,omitempty"`

	// Generation level (0 = root, negative = ancestors, positive = descendants)
	Generation int `json:"generation"`
}

// IsAlive returns true if the person has no death date
func (p *Person) IsAlive() bool {
	return p.DeathDate == nil
}

// Age returns the person's age at the given time, or at death if deceased
func (p *Person) Age(at time.Time) int {
	endDate := at
	if p.DeathDate != nil && p.DeathDate.Before(at) {
		endDate = *p.DeathDate
	}
	return yearsBetween(p.BirthDate, endDate)
}

// AgeAtDeath returns the age at death, or -1 if still alive
func (p *Person) AgeAtDeath() int {
	if p.DeathDate == nil {
		return -1
	}
	return yearsBetween(p.BirthDate, *p.DeathDate)
}

// FullName returns the person's full name
func (p *Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

// yearsBetween calculates the number of complete years between two dates
func yearsBetween(start, end time.Time) int {
	years := end.Year() - start.Year()
	if end.YearDay() < start.YearDay() {
		years--
	}
	return years
}

// NewPerson creates a new person with a unique ID
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
