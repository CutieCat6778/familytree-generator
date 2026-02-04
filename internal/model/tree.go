package model

import (
	"time"
)

type FamilyTree struct {
	ID           string             `json:"id"`
	RootPersonID string             `json:"root_person_id"`
	Persons      map[string]*Person `json:"persons"`
	Families     map[string]*Family `json:"families"`
	Generations  int                `json:"generations"`
	Country      string             `json:"country"`
	GeneratedAt  time.Time          `json:"generated_at"`
	Seed         int64              `json:"seed"`
}

func NewFamilyTree(id, country string, generations int, seed int64) *FamilyTree {
	return &FamilyTree{
		ID:          id,
		Persons:     make(map[string]*Person),
		Families:    make(map[string]*Family),
		Generations: generations,
		Country:     country,
		GeneratedAt: time.Now(),
		Seed:        seed,
	}
}

func (t *FamilyTree) SetRootPerson(p *Person) {
	t.RootPersonID = p.ID
	t.AddPerson(p)
}

func (t *FamilyTree) GetRootPerson() *Person {
	return t.Persons[t.RootPersonID]
}

func (t *FamilyTree) AddPerson(p *Person) {
	t.Persons[p.ID] = p
}

func (t *FamilyTree) GetPerson(id string) *Person {
	return t.Persons[id]
}

func (t *FamilyTree) AddFamily(f *Family) {
	t.Families[f.ID] = f
}

func (t *FamilyTree) GetFamily(id string) *Family {
	return t.Families[id]
}

func (t *FamilyTree) PersonCount() int {
	return len(t.Persons)
}

func (t *FamilyTree) FamilyCount() int {
	return len(t.Families)
}

func (t *FamilyTree) GetAncestors(personID string) []*Person {
	ancestors := make([]*Person, 0)
	person := t.GetPerson(personID)
	if person == nil {
		return ancestors
	}

	visited := make(map[string]bool)
	t.collectAncestors(person, &ancestors, visited)
	return ancestors
}

func (t *FamilyTree) collectAncestors(p *Person, ancestors *[]*Person, visited map[string]bool) {
	if p.FatherID != nil && !visited[*p.FatherID] {
		father := t.GetPerson(*p.FatherID)
		if father != nil {
			visited[*p.FatherID] = true
			*ancestors = append(*ancestors, father)
			t.collectAncestors(father, ancestors, visited)
		}
	}
	if p.MotherID != nil && !visited[*p.MotherID] {
		mother := t.GetPerson(*p.MotherID)
		if mother != nil {
			visited[*p.MotherID] = true
			*ancestors = append(*ancestors, mother)
			t.collectAncestors(mother, ancestors, visited)
		}
	}
}

func (t *FamilyTree) GetDescendants(personID string) []*Person {
	descendants := make([]*Person, 0)
	person := t.GetPerson(personID)
	if person == nil {
		return descendants
	}

	visited := make(map[string]bool)
	t.collectDescendants(person, &descendants, visited)
	return descendants
}

func (t *FamilyTree) collectDescendants(p *Person, descendants *[]*Person, visited map[string]bool) {
	for _, childID := range p.ChildrenIDs {
		if !visited[childID] {
			child := t.GetPerson(childID)
			if child != nil {
				visited[childID] = true
				*descendants = append(*descendants, child)
				t.collectDescendants(child, descendants, visited)
			}
		}
	}
}

func (t *FamilyTree) GetSiblings(personID string) []*Person {
	siblings := make([]*Person, 0)
	person := t.GetPerson(personID)
	if person == nil {
		return siblings
	}

	if person.FatherID != nil {
		father := t.GetPerson(*person.FatherID)
		if father != nil {
			for _, childID := range father.ChildrenIDs {
				if childID != personID {
					if child := t.GetPerson(childID); child != nil {
						siblings = append(siblings, child)
					}
				}
			}
		}
	}

	return siblings
}

func (t *FamilyTree) GetGeneration(gen int) []*Person {
	persons := make([]*Person, 0)
	for _, p := range t.Persons {
		if p.Generation == gen {
			persons = append(persons, p)
		}
	}
	return persons
}

func (t *FamilyTree) GetAllPersons() []*Person {
	persons := make([]*Person, 0, len(t.Persons))
	for _, p := range t.Persons {
		persons = append(persons, p)
	}
	return persons
}

func (t *FamilyTree) GetAllFamilies() []*Family {
	families := make([]*Family, 0, len(t.Families))
	for _, f := range t.Families {
		families = append(families, f)
	}
	return families
}
