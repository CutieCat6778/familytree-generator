package generator

import (
	"math"

	"github.com/familytree-generator/internal/data"
	"github.com/familytree-generator/internal/model"
	"github.com/familytree-generator/pkg/rand"
)

// ProbabilityEngine handles statistical calculations for family tree generation
type ProbabilityEngine struct {
	rng     *rand.SeededRandom
	stats   *data.CountryStats
	repo    *data.Repository
	country string
}

// NewProbabilityEngine creates a new probability engine
func NewProbabilityEngine(rng *rand.SeededRandom, stats *data.CountryStats, repo *data.Repository, country string) *ProbabilityEngine {
	return &ProbabilityEngine{
		rng:     rng,
		stats:   stats,
		repo:    repo,
		country: country,
	}
}

// CalculateChildrenCount determines the number of children using TFR (Total Fertility Rate)
func (p *ProbabilityEngine) CalculateChildrenCount(year int) int {
	// Use historical TFR data if available
	tfr := p.repo.GetFertilityRate(p.country, year)

	// TFR is children per woman, add some variance
	children := p.rng.NormalDistribution(tfr, tfr*0.25)

	// Round and ensure non-negative
	result := int(math.Round(children))
	if result < 0 {
		result = 0
	}
	if result > 12 {
		result = 12
	}

	return result
}

// CalculateChildrenCountLegacy uses birth rate (fallback method)
func (p *ProbabilityEngine) CalculateChildrenCountLegacy() int {
	birthRate := p.stats.BirthRate
	avgChildren := birthRate / 8.0

	if avgChildren < 0.5 {
		avgChildren = 0.5
	}
	if avgChildren > 8 {
		avgChildren = 8
	}

	children := p.rng.NormalDistribution(avgChildren, avgChildren*0.3)
	result := int(math.Round(children))
	if result < 0 {
		result = 0
	}
	if result > 12 {
		result = 12
	}

	return result
}

// CalculateDeathAge determines when a person dies based on life expectancy
func (p *ProbabilityEngine) CalculateDeathAge(health model.HealthProfile, birthYear int) int {
	baseAge := p.stats.LifeExpectancy

	// Adjust for historical period (life expectancy was lower in the past)
	if birthYear < 1950 {
		baseAge = baseAge * 0.75 // ~25% lower
	} else if birthYear < 1980 {
		baseAge = baseAge * 0.9 // ~10% lower
	}

	// Apply health modifiers
	if health.TobaccoUse {
		baseAge -= p.rng.Float64Range(5, 10)
	}

	if health.AlcoholConsumption > 10 {
		baseAge -= p.rng.Float64Range(2, 5)
	}

	// Add variance (standard deviation of ~8 years)
	deathAge := p.rng.NormalDistribution(baseAge, 8)

	if deathAge < 1 {
		deathAge = 1
	}
	if deathAge > 120 {
		deathAge = 120
	}

	return int(math.Round(deathAge))
}

// ShouldDieInInfancy determines if a person dies before age 1
func (p *ProbabilityEngine) ShouldDieInInfancy() bool {
	imr := p.stats.InfantMortality
	probability := imr / 1000.0
	return p.rng.Chance(probability)
}

// ShouldDieInYouth determines if a person dies before age 15
func (p *ProbabilityEngine) ShouldDieInYouth(birthYear int) bool {
	youthMortality := p.repo.GetYouthMortality(p.country, birthYear)
	// Youth mortality is percentage (e.g., 5.0 = 5%)
	probability := youthMortality / 100.0
	return p.rng.Chance(probability)
}

// CalculateMarriageAge determines when a person gets married using historical data
func (p *ProbabilityEngine) CalculateMarriageAge(gender model.Gender, birthYear int) int {
	var baseAge float64

	if gender == model.Female {
		// Use historical marriage age data for women
		baseAge = p.repo.GetMarriageAgeWomen(p.country, birthYear+25) // Approximate marriage year
	} else {
		// Men typically marry ~2-3 years older than women
		womenAge := p.repo.GetMarriageAgeWomen(p.country, birthYear+25)
		baseAge = womenAge + p.rng.Float64Range(2, 4)
	}

	// Add variance
	age := p.rng.NormalDistribution(baseAge, 3)

	if age < 18 {
		age = 18
	}
	if age > 50 {
		age = 50
	}

	return int(math.Round(age))
}

// ShouldGetDivorced determines if a marriage ends in divorce
func (p *ProbabilityEngine) ShouldGetDivorced(marriageYear int) bool {
	divorceRate := p.repo.GetDivorceRate(p.country, marriageYear)
	// Divorce rate is per 1000 people
	// Rough conversion to probability per marriage
	probability := divorceRate / 1000.0 * 10 // ~10 year window
	return p.rng.Chance(probability)
}

// ShouldBeBornOutsideMarriage determines if a child is born outside marriage
func (p *ProbabilityEngine) ShouldBeBornOutsideMarriage(birthYear int) bool {
	share := p.repo.GetBirthsOutsideMarriage(p.country, birthYear)
	// Share is a percentage
	return p.rng.Chance(share / 100.0)
}

func (p *ProbabilityEngine) ShouldBeUnderweight() bool {
	share := p.repo.GetUnderweightU5(p.country)
	return p.rng.Chance(share / 100.0)
}

func (p *ProbabilityEngine) DetermineResidence(birthYear int) model.ResidenceType {
	share := p.repo.GetUrbanShare(p.country, birthYear)
	if p.rng.Chance(share) {
		return model.Urban
	}
	return model.Rural
}

// ShouldBeSingleParent determines if a person becomes a single parent
func (p *ProbabilityEngine) ShouldBeSingleParent(year int) bool {
	share := p.repo.GetSingleParentShare(p.country, year)
	return p.rng.Chance(share / 100.0)
}

// ShouldGetMarried determines if a person ever gets married
func (p *ProbabilityEngine) ShouldGetMarried(birthYear int) bool {
	marriageRate := p.repo.GetMarriageRate(p.country, birthYear+28)
	// Marriage rate is per 1000 people per year
	// Higher rate = more likely to marry
	// Base probability around 85%, adjusted by rate
	baseProbability := 0.85
	if marriageRate < 5 {
		baseProbability = 0.70 // Lower marriage rates
	} else if marriageRate > 8 {
		baseProbability = 0.95 // Higher marriage rates
	}
	return p.rng.Chance(baseProbability)
}

// ShouldMigrate determines if a person migrates based on migration rate
func (p *ProbabilityEngine) ShouldMigrate() bool {
	migRate := p.stats.MigrationRate
	probability := math.Abs(migRate) / 1000.0 * 0.5
	return p.rng.Chance(probability)
}

// DetermineEmployment determines employment status based on age and rates
func (p *ProbabilityEngine) DetermineEmployment(age int) model.EmploymentStatus {
	if age < 16 {
		return model.Child
	}

	if age < 18 {
		return model.Student
	}

	if age >= 65 {
		return model.Retired
	}

	var unemploymentRate float64
	if age < 25 {
		unemploymentRate = p.stats.YouthUnemploymentRate
	} else {
		unemploymentRate = p.stats.UnemploymentRate
	}

	if age < 26 {
		studentProb := 0.3 + (p.stats.EducationExpenditure / 100)
		if p.rng.Chance(studentProb) {
			return model.Student
		}
	}

	if p.rng.Chance(unemploymentRate / 100) {
		return model.Unemployed
	}

	return model.Employed
}

// DetermineEducation determines education level based on country development
func (p *ProbabilityEngine) DetermineEducation() model.EducationLevel {
	eduExp := p.stats.EducationExpenditure
	gdp := p.stats.GDPPerCapita

	developmentScore := (gdp / 50000) + (eduExp / 10)
	if developmentScore > 1 {
		developmentScore = 1
	}

	roll := p.rng.Float64()

	if developmentScore > 0.7 {
		if roll < 0.05 {
			return model.NoEducation
		} else if roll < 0.20 {
			return model.Primary
		} else if roll < 0.55 {
			return model.Secondary
		} else {
			return model.Tertiary
		}
	} else if developmentScore > 0.4 {
		if roll < 0.10 {
			return model.NoEducation
		} else if roll < 0.35 {
			return model.Primary
		} else if roll < 0.75 {
			return model.Secondary
		} else {
			return model.Tertiary
		}
	} else {
		if roll < 0.20 {
			return model.NoEducation
		} else if roll < 0.55 {
			return model.Primary
		} else if roll < 0.85 {
			return model.Secondary
		} else {
			return model.Tertiary
		}
	}
}

// GenerateHealthProfile creates a health profile based on country statistics
func (p *ProbabilityEngine) GenerateHealthProfile() model.HealthProfile {
	alcohol := p.stats.AlcoholConsumption + p.rng.NormalDistribution(0, 2)
	if alcohol < 0 {
		alcohol = 0
	}
	return model.HealthProfile{
		AlcoholConsumption: alcohol,
		TobaccoUse:         p.rng.Chance(p.stats.TobaccoUse / 100),
	}
}

// CalculateChildBirthYear calculates when a child is born based on parent ages
func (p *ProbabilityEngine) CalculateChildBirthYear(motherBirthYear, childIndex int) int {
	marriageAge := p.CalculateMarriageAge(model.Female, motherBirthYear)
	motherAgeAtFirstChild := marriageAge + p.rng.IntRange(1, 3)

	// Space children 2-4 years apart
	spacing := childIndex * p.rng.IntRange(2, 4)

	return motherBirthYear + motherAgeAtFirstChild + spacing
}

// CalculateParentBirthYear calculates when a parent was born based on child's birth year
func (p *ProbabilityEngine) CalculateParentBirthYear(childBirthYear int, parentGender model.Gender) int {
	var ageGap int
	if parentGender == model.Female {
		ageGap = p.rng.IntRange(22, 32)
	} else {
		ageGap = p.rng.IntRange(25, 38)
	}

	return childBirthYear - ageGap
}

// CalculateSiblingCount determines number of siblings using TFR
func (p *ProbabilityEngine) CalculateSiblingCount(year int) int {
	count := p.CalculateChildrenCount(year)
	if count > 0 {
		count--
	}
	return count
}

// ShouldRemarry determines if a divorced/widowed person remarries
func (p *ProbabilityEngine) ShouldRemarry() bool {
	return p.rng.Chance(0.40) // ~40% remarriage rate
}

// Gender randomly determines gender with 50/50 probability
func (p *ProbabilityEngine) Gender() model.Gender {
	if p.rng.Bool() {
		return model.Male
	}
	return model.Female
}

// CalculateDivorceYear determines when a divorce happens if it does
func (p *ProbabilityEngine) CalculateDivorceYear(marriageYear int) int {
	// Most divorces happen within first 10 years
	yearsMarried := p.rng.IntRange(2, 15)
	return marriageYear + yearsMarried
}

// DetermineMaritalStatus determines the final marital status of a person
func (p *ProbabilityEngine) DetermineMaritalStatus(person *model.Person, hasSpouse bool, isDivorced bool) model.MaritalStatus {
	if !hasSpouse {
		return model.Single
	}
	if isDivorced {
		if len(person.SpouseIDs) > 1 {
			return model.Remarried
		}
		return model.Divorced
	}
	// Check if spouse is dead (would make them widowed)
	return model.Married
}
