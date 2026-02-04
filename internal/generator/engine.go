package generator

import (
	"fmt"
	"time"

	"github.com/familytree-generator/internal/data"
	"github.com/familytree-generator/internal/model"
	"github.com/familytree-generator/pkg/rand"
)

type Config struct {
	Country         string
	Generations     int
	Seed            int64
	StartYear       int
	RootGender      model.Gender
	IncludeExtended bool
}

func DefaultConfig() Config {
	return Config{
		Country:         "germany",
		Generations:     3,
		Seed:            time.Now().UnixNano(),
		StartYear:       1970,
		RootGender:      "",
		IncludeExtended: false,
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

	e.personGen = NewPersonGenerator(rng, repo, config.Country)
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

	return e.tree, nil
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
