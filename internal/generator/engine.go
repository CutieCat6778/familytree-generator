package generator

import (
	"fmt"
	"time"

	"github.com/familytree-generator/internal/data"
	"github.com/familytree-generator/internal/model"
	"github.com/familytree-generator/pkg/rand"
)

type Config struct {
	Country            string
	Generations        int
	Seed               int64
	StartYear          int
	RootGender         model.Gender
	IncludeExtended    bool
	LifeExpectancyMode LifeExpectancyMode
}

func DefaultConfig() Config {
	return Config{
		Country:            "germany",
		Generations:        3,
		Seed:               time.Now().UnixNano(),
		StartYear:          1970,
		RootGender:         "",
		IncludeExtended:    false,
		LifeExpectancyMode: LifeExpectancyTotal,
	}
}

type Engine struct {
	config    Config
	repo      *data.Repository
	rng       *rand.SeededRandom
	tree      *model.FamilyTree
	personGen *PersonGenerator
	familyBld *FamilyBuilder
}

func NewEngine(config Config, repo *data.Repository) *Engine {
	rng := rand.New(config.Seed)

	e := &Engine{
		config: config,
		repo:   repo,
		rng:    rng,
	}

	e.personGen = NewPersonGenerator(rng, repo, config.Country, config.LifeExpectancyMode)
	e.familyBld = NewFamilyBuilder(rng, e.personGen)

	return e
}

func (e *Engine) Generate() (*model.FamilyTree, error) {

	if err := e.repo.ValidateCountry(e.config.Country); err != nil {
		return nil, fmt.Errorf("invalid country: %w", err)
	}

	treeID := fmt.Sprintf("tree_%d", e.config.Seed)
	e.tree = model.NewFamilyTree(treeID, e.config.Country, e.config.Generations, e.config.Seed)

	rootGender := e.config.RootGender
	if rootGender == "" {
		rootGender = e.personGen.GetProbabilityEngine().Gender()
	}

	root := e.personGen.GeneratePerson(PersonOptions{
		Gender:     rootGender,
		BirthYear:  e.config.StartYear,
		Generation: 0,
	})

	e.tree.SetRootPerson(root)

	e.generateAncestors(root, e.config.Generations-1)

	e.generateDescendants(root, e.config.Generations-1)

	e.applyReferenceYearMortality()

	return e.tree, nil
}

func (e *Engine) applyReferenceYearMortality() {
	referenceYear := e.calculateReferenceYear()
	if referenceYear == 0 {
		return
	}

	refDate := time.Date(referenceYear, time.December, 31, 0, 0, 0, 0, time.UTC)
	prob := e.personGen.GetProbabilityEngine()

	for _, p := range e.tree.GetAllPersons() {
		maxAge := prob.MaxAllowedAge(p.BirthDate.Year(), p.Gender)
		ageAtReference := referenceYear - p.BirthDate.Year()
		if ageAtReference < 0 {
			ageAtReference = 0
		}

		if p.DeathDate != nil && p.DeathDate.Before(p.BirthDate) {
			adjusted := p.BirthDate
			p.DeathDate = &adjusted
			ensureDeathEvent(p)
		}

		if ageAtReference > maxAge {
			proposed := e.personGen.randomDateAtAge(p.BirthDate, maxAge)
			if proposed.After(refDate) {
				proposed = refDate
			}
			if p.DeathDate == nil || p.DeathDate.After(proposed) {
				p.DeathDate = &proposed
				ensureDeathEvent(p)
			}
		} else if p.DeathDate != nil {
			if p.AgeAtDeath() > maxAge {
				proposed := e.personGen.randomDateAtAge(p.BirthDate, maxAge)
				if proposed.After(refDate) {
					proposed = refDate
				}
				p.DeathDate = &proposed
				ensureDeathEvent(p)
			}
		}
	}
}

func (e *Engine) calculateReferenceYear() int {
	persons := e.tree.GetAllPersons()
	if len(persons) == 0 {
		return 0
	}

	maxGeneration := 0
	for _, p := range persons {
		if p.Generation > maxGeneration {
			maxGeneration = p.Generation
		}
	}

	var latestEvent time.Time
	for _, p := range persons {
		if p.Generation != maxGeneration {
			continue
		}
		for _, ev := range p.Events {
			if ev.Date.After(latestEvent) {
				latestEvent = ev.Date
			}
		}
		if p.BirthDate.After(latestEvent) {
			latestEvent = p.BirthDate
		}
		if p.DeathDate != nil && p.DeathDate.After(latestEvent) {
			latestEvent = *p.DeathDate
		}
	}

	if !latestEvent.IsZero() {
		return latestEvent.Year()
	}

	referenceYear := 0
	for _, p := range persons {
		birthYear := p.BirthDate.Year()
		if birthYear > referenceYear {
			referenceYear = birthYear
		}
	}
	if referenceYear == 0 {
		referenceYear = time.Now().Year()
	}
	return referenceYear
}

func ensureDeathEvent(person *model.Person) {
	if person.DeathDate == nil {
		return
	}

	for i := range person.Events {
		if person.Events[i].Type == model.EventDeath {
			person.Events[i].Date = *person.DeathDate
			person.Events[i].Location = person.CurrentCountry
			return
		}
	}

	deathEvent := model.NewLifeEvent(model.EventDeath, *person.DeathDate, person.CurrentCountry)
	person.Events = append(person.Events, deathEvent)
}

func (e *Engine) generateAncestors(person *model.Person, remainingGenerations int) {
	if remainingGenerations <= 0 {
		return
	}

	father := e.personGen.GenerateParent(person, model.Male)
	mother := e.personGen.GenerateParent(person, model.Female)

	e.tree.AddPerson(father)
	e.tree.AddPerson(mother)

	person.FatherID = &father.ID
	person.MotherID = &mother.ID
	father.ChildrenIDs = append(father.ChildrenIDs, person.ID)
	mother.ChildrenIDs = append(mother.ChildrenIDs, person.ID)

	family := e.familyBld.LinkSpouses(father, mother, e.tree)

	family.AddChild(person.ID)

	if e.config.IncludeExtended {
		e.familyBld.GenerateSiblings(person, father, mother, e.tree)
	}

	e.generateAncestors(father, remainingGenerations-1)
	e.generateAncestors(mother, remainingGenerations-1)
}

func (e *Engine) generateDescendants(person *model.Person, remainingGenerations int) {
	if remainingGenerations <= 0 {
		return
	}

	if person.DeathDate != nil {
		deathAge := person.AgeAtDeath()
		if deathAge < 18 {
			return
		}
	}

	spouse := e.personGen.GenerateSpouse(person)
	e.tree.AddPerson(spouse)

	var husband, wife *model.Person
	if person.Gender == model.Male {
		husband = person
		wife = spouse
	} else {
		husband = spouse
		wife = person
	}

	family := e.familyBld.LinkSpouses(husband, wife, e.tree)

	children := e.familyBld.GenerateChildren(family, husband, wife, e.tree)

	for _, child := range children {
		e.generateDescendants(child, remainingGenerations-1)
	}
}

func (e *Engine) GetTree() *model.FamilyTree {
	return e.tree
}

func (e *Engine) GetConfig() Config {
	return e.config
}
