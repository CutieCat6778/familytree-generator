package generator

import (
	"fmt"
	"time"

	"github.com/familytree-generator/internal/data"
	"github.com/familytree-generator/internal/model"
	"github.com/familytree-generator/pkg/rand"
)

// Config holds configuration for tree generation
type Config struct {
	Country         string       // Country slug (e.g., "united-states")
	Generations     int          // Number of generations to generate
	Seed            int64        // Random seed for reproducibility
	StartYear       int          // Birth year of the root person
	RootGender      model.Gender // Gender of root person (empty = random)
	IncludeExtended bool         // Include extended family (siblings, cousins)
}

// DefaultConfig returns the default configuration
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

// Engine is the main family tree generation engine
type Engine struct {
	config    Config
	repo      *data.Repository
	rng       *rand.SeededRandom
	tree      *model.FamilyTree
	personGen *PersonGenerator
	familyBld *FamilyBuilder
}

// NewEngine creates a new generation engine
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

// Generate creates a complete family tree
func (e *Engine) Generate() (*model.FamilyTree, error) {
	// Validate country
	if err := e.repo.ValidateCountry(e.config.Country); err != nil {
		return nil, fmt.Errorf("invalid country: %w", err)
	}

	// Create tree
	treeID := fmt.Sprintf("tree_%d", e.config.Seed)
	e.tree = model.NewFamilyTree(treeID, e.config.Country, e.config.Generations, e.config.Seed)

	// Generate root person
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

	// Generate ancestors (going backwards)
	e.generateAncestors(root, e.config.Generations-1)

	// Generate descendants (going forwards)
	e.generateDescendants(root, e.config.Generations-1)

	return e.tree, nil
}

// generateAncestors recursively generates parents, grandparents, etc.
func (e *Engine) generateAncestors(person *model.Person, remainingGenerations int) {
	if remainingGenerations <= 0 {
		return
	}

	// Generate parents
	father := e.personGen.GenerateParent(person, model.Male)
	mother := e.personGen.GenerateParent(person, model.Female)

	// Add to tree
	e.tree.AddPerson(father)
	e.tree.AddPerson(mother)

	// Link to child
	person.FatherID = &father.ID
	person.MotherID = &mother.ID
	father.ChildrenIDs = append(father.ChildrenIDs, person.ID)
	mother.ChildrenIDs = append(mother.ChildrenIDs, person.ID)

	// Create family
	family := e.familyBld.LinkSpouses(father, mother, e.tree)

	// Add person as child of this family
	family.AddChild(person.ID)

	// Generate siblings if extended family is enabled
	if e.config.IncludeExtended {
		e.familyBld.GenerateSiblings(person, father, mother, e.tree)
	}

	// Recursively generate ancestors for parents
	e.generateAncestors(father, remainingGenerations-1)
	e.generateAncestors(mother, remainingGenerations-1)
}

// generateDescendants recursively generates children, grandchildren, etc.
func (e *Engine) generateDescendants(person *model.Person, remainingGenerations int) {
	if remainingGenerations <= 0 {
		return
	}

	// Check if person died before they could have children
	if person.DeathDate != nil {
		deathAge := person.AgeAtDeath()
		if deathAge < 18 {
			return // Too young to have children
		}
	}

	// Generate spouse
	spouse := e.personGen.GenerateSpouse(person)
	e.tree.AddPerson(spouse)

	// Create family
	var husband, wife *model.Person
	if person.Gender == model.Male {
		husband = person
		wife = spouse
	} else {
		husband = spouse
		wife = person
	}

	family := e.familyBld.LinkSpouses(husband, wife, e.tree)

	// Generate children
	children := e.familyBld.GenerateChildren(family, husband, wife, e.tree)

	// Recursively generate descendants for children
	for _, child := range children {
		e.generateDescendants(child, remainingGenerations-1)
	}
}

// GetTree returns the generated tree
func (e *Engine) GetTree() *model.FamilyTree {
	return e.tree
}

// GetConfig returns the engine configuration
func (e *Engine) GetConfig() Config {
	return e.config
}
