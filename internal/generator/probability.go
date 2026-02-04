package generator

import (
	"math"

	"github.com/familytree-generator/internal/data"
	"github.com/familytree-generator/internal/model"
	"github.com/familytree-generator/pkg/rand"
)


type ProbabilityEngine struct {
	rng     *rand.SeededRandom
	stats   *data.CountryStats
	repo    *data.Repository
	country string
}


func NewProbabilityEngine(rng *rand.SeededRandom, stats *data.CountryStats, repo *data.Repository, country string) *ProbabilityEngine {
	return &ProbabilityEngine{
		rng:     rng,
		stats:   stats,
		repo:    repo,
		country: country,
	}
}


func (p *ProbabilityEngine) CalculateChildrenCount(year int) int {
	
	tfr := p.repo.GetFertilityRate(p.country, year)

	
	children := p.rng.NormalDistribution(tfr, tfr*0.25)

	
	result := int(math.Round(children))
	if result < 0 {
		result = 0
	}
	if result > 12 {
		result = 12
	}

	return result
}


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


func (p *ProbabilityEngine) CalculateDeathAge(health model.HealthProfile, birthYear int) int {
	baseAge := p.stats.LifeExpectancy

	
	if birthYear < 1950 {
		baseAge = baseAge * 0.75 
	} else if birthYear < 1980 {
		baseAge = baseAge * 0.9 
	}

	
	if health.TobaccoUse {
		baseAge -= p.rng.Float64Range(5, 10)
	}

	if health.AlcoholConsumption > 10 {
		baseAge -= p.rng.Float64Range(2, 5)
	}

	
	deathAge := p.rng.NormalDistribution(baseAge, 8)

	if deathAge < 1 {
		deathAge = 1
	}
	if deathAge > 120 {
		deathAge = 120
	}

	return int(math.Round(deathAge))
}


func (p *ProbabilityEngine) ShouldDieInInfancy() bool {
	imr := p.stats.InfantMortality
	probability := imr / 1000.0
	return p.rng.Chance(probability)
}


func (p *ProbabilityEngine) ShouldDieInYouth(birthYear int) bool {
	youthMortality := p.repo.GetYouthMortality(p.country, birthYear)
	
	probability := youthMortality / 100.0
	return p.rng.Chance(probability)
}


func (p *ProbabilityEngine) CalculateMarriageAge(gender model.Gender, birthYear int) int {
	var baseAge float64

	if gender == model.Female {
		
		baseAge = p.repo.GetMarriageAgeWomen(p.country, birthYear+25) 
	} else {
		
		womenAge := p.repo.GetMarriageAgeWomen(p.country, birthYear+25)
		baseAge = womenAge + p.rng.Float64Range(2, 4)
	}

	
	age := p.rng.NormalDistribution(baseAge, 3)

	if age < 18 {
		age = 18
	}
	if age > 50 {
		age = 50
	}

	return int(math.Round(age))
}


func (p *ProbabilityEngine) ShouldGetDivorced(marriageYear int) bool {
	divorceRate := p.repo.GetDivorceRate(p.country, marriageYear)
	
	
	probability := divorceRate / 1000.0 * 10 
	return p.rng.Chance(probability)
}


func (p *ProbabilityEngine) ShouldBeBornOutsideMarriage(birthYear int) bool {
	share := p.repo.GetBirthsOutsideMarriage(p.country, birthYear)
	
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


func (p *ProbabilityEngine) ShouldBeSingleParent(year int) bool {
	share := p.repo.GetSingleParentShare(p.country, year)
	return p.rng.Chance(share / 100.0)
}


func (p *ProbabilityEngine) ShouldGetMarried(birthYear int) bool {
	marriageRate := p.repo.GetMarriageRate(p.country, birthYear+28)
	
	
	
	baseProbability := 0.85
	if marriageRate < 5 {
		baseProbability = 0.70 
	} else if marriageRate > 8 {
		baseProbability = 0.95 
	}
	return p.rng.Chance(baseProbability)
}


func (p *ProbabilityEngine) ShouldMigrate() bool {
	migRate := p.stats.MigrationRate
	probability := math.Abs(migRate) / 1000.0 * 0.5
	return p.rng.Chance(probability)
}


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


func (p *ProbabilityEngine) CalculateChildBirthYear(motherBirthYear, childIndex int) int {
	marriageAge := p.CalculateMarriageAge(model.Female, motherBirthYear)
	motherAgeAtFirstChild := marriageAge + p.rng.IntRange(1, 3)

	
	spacing := childIndex * p.rng.IntRange(2, 4)

	return motherBirthYear + motherAgeAtFirstChild + spacing
}


func (p *ProbabilityEngine) CalculateParentBirthYear(childBirthYear int, parentGender model.Gender) int {
	var ageGap int
	if parentGender == model.Female {
		ageGap = p.rng.IntRange(22, 32)
	} else {
		ageGap = p.rng.IntRange(25, 38)
	}

	return childBirthYear - ageGap
}


func (p *ProbabilityEngine) CalculateSiblingCount(year int) int {
	count := p.CalculateChildrenCount(year)
	if count > 0 {
		count--
	}
	return count
}


func (p *ProbabilityEngine) ShouldRemarry() bool {
	return p.rng.Chance(0.40) 
}


func (p *ProbabilityEngine) Gender() model.Gender {
	if p.rng.Bool() {
		return model.Male
	}
	return model.Female
}


func (p *ProbabilityEngine) CalculateDivorceYear(marriageYear int) int {
	
	yearsMarried := p.rng.IntRange(2, 15)
	return marriageYear + yearsMarried
}


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
	
	return model.Married
}
