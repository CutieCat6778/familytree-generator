package generator

import (
	"fmt"
	"time"

	"github.com/familytree-generator/internal/model"
	"github.com/familytree-generator/pkg/rand"
)


type FamilyBuilder struct {
	rng         *rand.SeededRandom
	personGen   *PersonGenerator
	familyCounter uint64
}


func NewFamilyBuilder(rng *rand.SeededRandom, personGen *PersonGenerator) *FamilyBuilder {
	return &FamilyBuilder{
		rng:         rng,
		personGen:   personGen,
		familyCounter: 0,
	}
}


func (b *FamilyBuilder) CreateFamily(husband, wife *model.Person) *model.Family {
	b.familyCounter++
	id := fmt.Sprintf("F%05d", b.familyCounter)

	
	marriageYear := b.calculateMarriageYear(husband, wife)
	marriageDate := time.Date(marriageYear, time.Month(b.rng.IntRange(1, 12)), b.rng.IntRange(1, 28), 0, 0, 0, 0, time.UTC)

	family := model.NewFamily(id, marriageDate)
	family.SetHusband(husband.ID)
	family.SetWife(wife.ID)

	
	husband.SpouseIDs = append(husband.SpouseIDs, wife.ID)
	wife.SpouseIDs = append(wife.SpouseIDs, husband.ID)

	
	husband.MarriageAge = marriageYear - husband.BirthDate.Year()
	wife.MarriageAge = marriageYear - wife.BirthDate.Year()

	
	husband.MaritalStatus = model.Married
	wife.MaritalStatus = model.Married

	
	prob := b.personGen.GetProbabilityEngine()
	if prob.ShouldGetDivorced(marriageYear) {
		divorceYear := prob.CalculateDivorceYear(marriageYear)
		divorceDate := time.Date(divorceYear, time.Month(b.rng.IntRange(1, 12)), b.rng.IntRange(1, 28), 0, 0, 0, 0, time.UTC)

		
		husbandAlive := husband.DeathDate == nil || husband.DeathDate.After(divorceDate)
		wifeAlive := wife.DeathDate == nil || wife.DeathDate.After(divorceDate)

		if husbandAlive && wifeAlive {
			family.DivorceDate = &divorceDate
			husband.MaritalStatus = model.Divorced
			wife.MaritalStatus = model.Divorced

			
			divorceEvent := model.NewLifeEvent(model.EventDivorce, divorceDate, husband.CurrentCountry).
				WithRelatedID(wife.ID)
			husband.Events = append(husband.Events, divorceEvent)

			divorceEventWife := model.NewLifeEvent(model.EventDivorce, divorceDate, wife.CurrentCountry).
				WithRelatedID(husband.ID)
			wife.Events = append(wife.Events, divorceEventWife)
		}
	}

	
	marriageEvent := model.NewLifeEvent(model.EventMarriage, marriageDate, husband.CurrentCountry).
		WithRelatedID(wife.ID)
	husband.Events = append(husband.Events, marriageEvent)

	marriageEventWife := model.NewLifeEvent(model.EventMarriage, marriageDate, wife.CurrentCountry).
		WithRelatedID(husband.ID)
	wife.Events = append(wife.Events, marriageEventWife)

	return family
}


func (b *FamilyBuilder) calculateMarriageYear(husband, wife *model.Person) int {
	
	husbandMarriageAge := b.personGen.GetProbabilityEngine().CalculateMarriageAge(model.Male, husband.BirthDate.Year())
	wifeMarriageAge := b.personGen.GetProbabilityEngine().CalculateMarriageAge(model.Female, wife.BirthDate.Year())

	husbandMarriageYear := husband.BirthDate.Year() + husbandMarriageAge
	wifeMarriageYear := wife.BirthDate.Year() + wifeMarriageAge

	
	marriageYear := husbandMarriageYear
	if wifeMarriageYear > marriageYear {
		marriageYear = wifeMarriageYear
	}

	
	if husband.DeathDate != nil && marriageYear > husband.DeathDate.Year() {
		marriageYear = husband.DeathDate.Year() - 1
	}
	if wife.DeathDate != nil && marriageYear > wife.DeathDate.Year() {
		marriageYear = wife.DeathDate.Year() - 1
	}

	return marriageYear
}


func (b *FamilyBuilder) GenerateChildren(family *model.Family, husband, wife *model.Person, tree *model.FamilyTree) []*model.Person {
	prob := b.personGen.GetProbabilityEngine()
	
	numChildren := prob.CalculateChildrenCount(family.MarriedDate.Year())

	children := make([]*model.Person, 0, numChildren)

	for i := 0; i < numChildren; i++ {
		child := b.personGen.GenerateChild(husband, wife, i)

		
		if child.BirthDate.Before(family.MarriedDate) {
			
			yearsAfterMarriage := b.rng.IntRange(1, 3) + i*b.rng.IntRange(2, 4)
			child.BirthDate = family.MarriedDate.AddDate(yearsAfterMarriage, b.rng.IntRange(0, 11), b.rng.IntRange(1, 28))
		}

		
		if wife.DeathDate != nil && child.BirthDate.After(*wife.DeathDate) {
			continue 
		}

		
		motherAge := child.BirthDate.Year() - wife.BirthDate.Year()
		if motherAge < 16 || motherAge > 50 {
			continue 
		}

		
		family.AddChild(child.ID)
		husband.ChildrenIDs = append(husband.ChildrenIDs, child.ID)
		wife.ChildrenIDs = append(wife.ChildrenIDs, child.ID)
		husband.NumberOfChildren++
		wife.NumberOfChildren++

		tree.AddPerson(child)
		children = append(children, child)
	}

	return children
}


func (b *FamilyBuilder) LinkSpouses(husband, wife *model.Person, tree *model.FamilyTree) *model.Family {
	family := b.CreateFamily(husband, wife)
	tree.AddFamily(family)
	return family
}


func (b *FamilyBuilder) GenerateSiblings(person *model.Person, father, mother *model.Person, tree *model.FamilyTree) []*model.Person {
	prob := b.personGen.GetProbabilityEngine()
	
	numSiblings := prob.CalculateSiblingCount(person.BirthDate.Year())

	siblings := make([]*model.Person, 0, numSiblings)

	for i := 0; i < numSiblings; i++ {
		sibling := b.personGen.GenerateSibling(person, father, mother, i)

		
		motherAge := sibling.BirthDate.Year() - mother.BirthDate.Year()
		if motherAge < 16 || motherAge > 50 {
			continue
		}

		
		if mother.DeathDate != nil && sibling.BirthDate.After(*mother.DeathDate) {
			continue
		}

		
		father.ChildrenIDs = append(father.ChildrenIDs, sibling.ID)
		mother.ChildrenIDs = append(mother.ChildrenIDs, sibling.ID)

		tree.AddPerson(sibling)
		siblings = append(siblings, sibling)
	}

	return siblings
}
