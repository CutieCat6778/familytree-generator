package generator

import (
	"fmt"
	"math"
	"time"

	"github.com/familytree-generator/internal/data"
	"github.com/familytree-generator/internal/model"
	"github.com/familytree-generator/pkg/rand"
)

// PersonGenerator handles the creation of individual persons
type PersonGenerator struct {
	rng      *rand.SeededRandom
	repo     *data.Repository
	prob     *ProbabilityEngine
	country  string
	idCounter uint64
	countryOptions []string
}

// NewPersonGenerator creates a new person generator
func NewPersonGenerator(rng *rand.SeededRandom, repo *data.Repository, country string) *PersonGenerator {
	stats := repo.GetCountryStats(country)
	return &PersonGenerator{
		rng:       rng,
		repo:      repo,
		prob:      NewProbabilityEngine(rng, stats, repo, country),
		country:   country,
		idCounter: 0,
		countryOptions: repo.GetAvailableCountrySlugs(),
	}
}

// PersonOptions configures how a person is generated
type PersonOptions struct {
	Gender       model.Gender
	BirthYear    int
	Generation   int
	FatherID     *string
	MotherID     *string
	LastName     string // If empty, will be generated
	WealthIndex  *float64
}

// GeneratePerson creates a new person with the given options
func (g *PersonGenerator) GeneratePerson(opts PersonOptions) *model.Person {
	g.idCounter++
	id := fmt.Sprintf("P%05d", g.idCounter)

	// Determine gender if not specified
	gender := opts.Gender
	if gender == "" {
		gender = g.prob.Gender()
	}

	// Generate name
	firstName := g.generateFirstName(gender, opts.BirthYear)
	lastName := opts.LastName
	if lastName == "" {
		lastName = g.generateLastName()
	}

	// Generate birth date
	birthDate := g.generateBirthDate(opts.BirthYear)

	// Create person
	person := model.NewPerson(id, firstName, lastName, gender, birthDate, g.country, opts.Generation)

	// Set parents
	person.FatherID = opts.FatherID
	person.MotherID = opts.MotherID

	// Check if born outside marriage
	person.BornOutsideMarriage = g.prob.ShouldBeBornOutsideMarriage(opts.BirthYear)
	person.Underweight = g.prob.ShouldBeUnderweight()
	person.Residence = g.determineResidenceForCountry(g.country, opts.BirthYear)
	person.GDPPerCapita = g.repo.GetGDPPerCapita(g.country)
	person.WealthIndex = g.getWealthIndex(opts.WealthIndex)
	g.assignWealth(person)

	// Generate health profile
	person.Health = g.prob.GenerateHealthProfile()

	// Determine if person has died
	deathAge := g.prob.CalculateDeathAge(person.Health, opts.BirthYear)

	// Check infant mortality
	if g.prob.ShouldDieInInfancy() {
		deathAge = 0
		infantDeathDate := birthDate.AddDate(0, g.rng.IntRange(0, 11), g.rng.IntRange(0, 28))
		person.DeathDate = &infantDeathDate
	} else if g.prob.ShouldDieInYouth(opts.BirthYear) {
		// Died in childhood (ages 1-14)
		deathAge = g.rng.IntRange(1, 14)
		youthDeathDate := birthDate.AddDate(deathAge, g.rng.IntRange(0, 11), g.rng.IntRange(1, 28))
		person.DeathDate = &youthDeathDate
	} else if deathAge <= person.Age(time.Now()) {
		// Person has died of old age
		deathDate := birthDate.AddDate(deathAge, g.rng.IntRange(0, 11), g.rng.IntRange(1, 28))
		person.DeathDate = &deathDate
	}

	// Determine education and employment based on current age
	currentAge := person.Age(time.Now())
	if person.DeathDate != nil {
		currentAge = person.AgeAtDeath()
	}

	person.Education = g.prob.DetermineEducation()
	person.Employment = g.prob.DetermineEmployment(currentAge)

	// Initialize marital status
	person.MaritalStatus = model.Single

	// Add birth event
	birthEvent := model.NewLifeEvent(model.EventBirth, birthDate, g.country)
	person.Events = append(person.Events, birthEvent)

	g.maybeMigrate(person)

	// Add death event if deceased
	if person.DeathDate != nil {
		deathEvent := model.NewLifeEvent(model.EventDeath, *person.DeathDate, person.CurrentCountry)
		person.Events = append(person.Events, deathEvent)
	}

	return person
}

// generateFirstName selects a first name based on country and gender
func (g *PersonGenerator) generateFirstName(gender model.Gender, birthYear int) string {
	genderStr := string(gender)
	names := g.repo.GetForenamesByGender(g.country, genderStr)

	if len(names) == 0 {
		// Fallback names
		if gender == model.Male {
			return rand.Choice(g.rng, []string{"John", "James", "William", "Michael", "David"})
		}
		return rand.Choice(g.rng, []string{"Mary", "Elizabeth", "Sarah", "Emma", "Anna"})
	}

	// Weighted selection favoring higher-ranked names
	weights := make([]float64, len(names))
	for i, n := range names {
		// Higher rank (lower index) = higher weight
		weight := 1.0 / float64(n.Index+1)
		if n.Year > 0 && birthYear > 0 {
			yearDiff := math.Abs(float64(n.Year - birthYear))
			weight *= 1.0 / (1.0 + yearDiff/10.0)
		}
		weights[i] = weight
	}

	idx := g.rng.WeightedChoice(weights)
	name := names[idx].RomanizedName
	if name == "" {
		name = names[idx].LocalizedName
	}

	return name
}

// generateLastName selects a last name based on country
func (g *PersonGenerator) generateLastName() string {
	surnames := g.repo.GetSurnames(g.country)

	if len(surnames) == 0 {
		// Fallback names
		return rand.Choice(g.rng, []string{"Smith", "Johnson", "Williams", "Brown", "Jones"})
	}

	// Weighted selection favoring more common surnames
	weights := make([]float64, len(surnames))
	rankCounts := make(map[int]int, len(surnames))
	for _, s := range surnames {
		rankCounts[s.Rank]++
	}
	for i, s := range surnames {
		// Higher rank (lower index) = higher weight.
		// If multiple surnames share a rank, split that weight evenly to preserve totals.
		weight := 1.0 / float64(s.Rank+1)
		if count := rankCounts[s.Rank]; count > 1 {
			weight /= float64(count)
		}
		weights[i] = weight
	}

	idx := g.rng.WeightedChoice(weights)
	name := surnames[idx].RomanizedName
	if name == "" {
		name = surnames[idx].LocalizedName
	}

	return name
}

// generateBirthDate creates a birth date for the given year
func (g *PersonGenerator) generateBirthDate(year int) time.Time {
	month := g.rng.IntRange(1, 12)
	day := g.rng.IntRange(1, 28) // Safe day range for all months

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func (g *PersonGenerator) getWealthIndex(value *float64) float64 {
	if value != nil && *value > 0 {
		return *value
	}
	return g.randomWealthIndex()
}

func (g *PersonGenerator) randomWealthIndex() float64 {
	sigma := 0.6
	mu := -0.5 * sigma * sigma
	value := math.Exp(mu + sigma*g.rng.NormFloat64())
	if value < 0.3 {
		return 0.3
	}
	if value > 4 {
		return 4
	}
	return value
}

func (g *PersonGenerator) blendWealthIndex(parent float64, weight float64) float64 {
	if parent <= 0 {
		return g.randomWealthIndex()
	}
	if weight < 0 {
		weight = 0
	}
	if weight > 1 {
		weight = 1
	}
	random := g.randomWealthIndex()
	value := parent*weight + random*(1-weight)
	if value < 0.3 {
		return 0.3
	}
	if value > 4 {
		return 4
	}
	return value
}

func (g *PersonGenerator) assignWealth(person *model.Person) {
	person.GDPPerCapita = g.repo.GetGDPPerCapita(person.CurrentCountry)
	person.FamilyWealth = person.GDPPerCapita * person.WealthIndex
	person.IsRich = person.WealthIndex >= 1.5
}

func (g *PersonGenerator) determineResidenceForCountry(country string, year int) model.ResidenceType {
	share := g.repo.GetUrbanShare(country, year)
	if g.rng.Chance(share) {
		return model.Urban
	}
	return model.Rural
}

func (g *PersonGenerator) maybeMigrate(person *model.Person) {
	if !g.prob.ShouldMigrate() || len(g.countryOptions) < 2 {
		return
	}

	migrationAge := g.rng.IntRange(18, 45)
	migrationDate := person.BirthDate.AddDate(migrationAge, g.rng.IntRange(0, 11), g.rng.IntRange(1, 28))
	if person.DeathDate != nil && !person.DeathDate.After(migrationDate) {
		return
	}

	var destination string
	for i := 0; i < 10; i++ {
		candidate := rand.Choice(g.rng, g.countryOptions)
		if candidate != "" && candidate != person.BirthCountry {
			destination = candidate
			break
		}
	}
	if destination == "" {
		for _, candidate := range g.countryOptions {
			if candidate != person.BirthCountry {
				destination = candidate
				break
			}
		}
	}
	if destination == "" {
		return
	}

	origin := person.CurrentCountry
	person.CurrentCountry = destination
	person.Events = append(person.Events, model.NewLifeEvent(model.EventMigration, migrationDate, destination).WithDescription(origin))
	person.Residence = g.determineResidenceForCountry(destination, migrationDate.Year())
	g.assignWealth(person)
}

// GenerateSpouse creates a spouse for the given person
func (g *PersonGenerator) GenerateSpouse(person *model.Person) *model.Person {
	// Determine spouse gender
	var spouseGender model.Gender
	if person.Gender == model.Male {
		spouseGender = model.Female
	} else {
		spouseGender = model.Male
	}

	// Spouse typically within a few years of age
	ageDiff := g.rng.IntRange(-5, 5)
	spouseBirthYear := person.BirthDate.Year() + ageDiff

	spouseWealth := g.blendWealthIndex(person.WealthIndex, 0.7)
	spouse := g.GeneratePerson(PersonOptions{
		Gender:     spouseGender,
		BirthYear:  spouseBirthYear,
		Generation: person.Generation,
		WealthIndex: &spouseWealth,
	})

	return spouse
}

// GenerateChild creates a child for the given parents
func (g *PersonGenerator) GenerateChild(father, mother *model.Person, childIndex int) *model.Person {
	// Calculate birth year based on mother's age
	birthYear := g.prob.CalculateChildBirthYear(mother.BirthDate.Year(), childIndex)

	// Child inherits father's last name
	lastName := father.LastName
	parentWealth := (father.WealthIndex + mother.WealthIndex) / 2
	childWealth := g.blendWealthIndex(parentWealth, 0.7)

	child := g.GeneratePerson(PersonOptions{
		BirthYear:  birthYear,
		Generation: father.Generation + 1,
		FatherID:   &father.ID,
		MotherID:   &mother.ID,
		LastName:   lastName,
		WealthIndex: &childWealth,
	})

	return child
}

// GenerateParent creates a parent for the given person
func (g *PersonGenerator) GenerateParent(child *model.Person, gender model.Gender) *model.Person {
	birthYear := g.prob.CalculateParentBirthYear(child.BirthDate.Year(), gender)
	parentWealth := g.blendWealthIndex(child.WealthIndex, 0.6)

	opts := PersonOptions{
		Gender:     gender,
		BirthYear:  birthYear,
		Generation: child.Generation - 1,
		WealthIndex: &parentWealth,
	}

	// If generating father and child has a last name, father should have same
	if gender == model.Male && child.LastName != "" {
		opts.LastName = child.LastName
	}

	return g.GeneratePerson(opts)
}

// GenerateSibling creates a sibling for the given person
func (g *PersonGenerator) GenerateSibling(person *model.Person, father, mother *model.Person, siblingIndex int) *model.Person {
	// Siblings born within a few years
	ageDiff := g.rng.IntRange(-8, 8)
	birthYear := person.BirthDate.Year() + ageDiff

	// Ensure sibling is born after parents are at least 18
	minBirthYear := mother.BirthDate.Year() + 18
	if birthYear < minBirthYear {
		birthYear = minBirthYear + siblingIndex*2
	}

	parentWealth := (father.WealthIndex + mother.WealthIndex) / 2
	siblingWealth := g.blendWealthIndex(parentWealth, 0.8)
	sibling := g.GeneratePerson(PersonOptions{
		BirthYear:  birthYear,
		Generation: person.Generation,
		FatherID:   &father.ID,
		MotherID:   &mother.ID,
		LastName:   father.LastName,
		WealthIndex: &siblingWealth,
	})

	return sibling
}

// GetProbabilityEngine returns the probability engine for external use
func (g *PersonGenerator) GetProbabilityEngine() *ProbabilityEngine {
	return g.prob
}

// GetCurrentID returns the current ID counter
func (g *PersonGenerator) GetCurrentID() uint64 {
	return g.idCounter
}
