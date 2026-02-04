package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/familytree-generator/internal/model"
)


func WriteJSON(tree *model.FamilyTree, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(tree); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}

	return nil
}


func WriteJSONCompact(tree *model.FamilyTree, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(tree); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}

	return nil
}


func TreeToJSON(tree *model.FamilyTree) ([]byte, error) {
	return json.MarshalIndent(tree, "", "  ")
}


func TreeToJSONCompact(tree *model.FamilyTree) ([]byte, error) {
	return json.Marshal(tree)
}


type VisualizationData struct {
	ID          string                `json:"id"`
	RootID      string                `json:"root_id"`
	Country     string                `json:"country"`
	Generations int                   `json:"generations"`
	Seed        int64                 `json:"seed"`
	ReferenceYear int                `json:"reference_year"`
	Nodes       []VisualizationNode   `json:"nodes"`
	Edges       []VisualizationEdge   `json:"edges"`
	Stats       VisualizationStats    `json:"stats"`
}


type VisualizationNode struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	FirstName           string  `json:"first_name"`
	LastName            string  `json:"last_name"`
	Gender              string  `json:"gender"`
	BirthYear           int     `json:"birth_year"`
	DeathYear           *int    `json:"death_year,omitempty"`
	IsAlive             bool    `json:"is_alive"`
	Generation          int     `json:"generation"`
	MaritalStatus       string  `json:"marital_status"`
	MarriageAge         int     `json:"marriage_age,omitempty"`
	NumberOfChildren    int     `json:"number_of_children"`
	Education           string  `json:"education"`
	Employment          string  `json:"employment"`
	AlcoholConsumption  float64 `json:"alcohol_consumption"`
	TobaccoUse          bool    `json:"tobacco_use"`
	BornOutsideMarriage bool    `json:"born_outside_marriage"`
	IsSingleParent      bool    `json:"is_single_parent"`
	Underweight         bool    `json:"underweight"`
	Residence           string  `json:"residence"`
	GDPPerCapita        float64 `json:"gdp_per_capita"`
	WealthIndex         float64 `json:"wealth_index"`
	FamilyWealth        float64 `json:"family_wealth"`
	IsRich              bool    `json:"is_rich"`
	Country             string  `json:"country"`
	CurrentCountry      string  `json:"current_country"`
}


type VisualizationEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"` 
}


type VisualizationStats struct {
	TotalPersons        int     `json:"total_persons"`
	TotalFamilies       int     `json:"total_families"`
	LivingPersons       int     `json:"living_persons"`
	DeceasedPersons     int     `json:"deceased_persons"`
	AverageAge          float64 `json:"average_age"`
	OldestPerson        int     `json:"oldest_person_age"`
	TotalChildren       int     `json:"total_children"`
	AverageChildren     float64 `json:"average_children"`
	DivorceCount        int     `json:"divorce_count"`
	SingleCount         int     `json:"single_count"`
	MarriedCount        int     `json:"married_count"`
	MaleCount           int     `json:"male_count"`
	FemaleCount         int     `json:"female_count"`
	BirthsOutsideMarriage int   `json:"births_outside_marriage"`
	TertiaryEducation   int     `json:"tertiary_education"`
	EmployedCount       int     `json:"employed_count"`
}


func WriteVisualizationJSON(tree *model.FamilyTree, filepath string) error {
	data := TreeToVisualizationData(tree)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}

	return nil
}


func TreeToVisualizationData(tree *model.FamilyTree) *VisualizationData {
	persons := tree.GetAllPersons()
	referenceYear := 0
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
		referenceYear = latestEvent.Year()
	} else {
		for _, p := range persons {
			birthYear := p.BirthDate.Year()
			if birthYear > referenceYear {
				referenceYear = birthYear
			}
		}
		if referenceYear == 0 {
			referenceYear = time.Now().Year()
		}
	}

	data := &VisualizationData{
		ID:          tree.ID,
		RootID:      tree.RootPersonID,
		Country:     tree.Country,
		Generations: tree.Generations,
		Seed:        tree.Seed,
		ReferenceYear: referenceYear,
		Nodes:       make([]VisualizationNode, 0),
		Edges:       make([]VisualizationEdge, 0),
	}

	var totalAge float64
	var ageCount int
	var oldestAge int

	
	for _, p := range persons {
		var deathYear *int
		if p.DeathDate != nil {
			year := p.DeathDate.Year()
			deathYear = &year
		}

		node := VisualizationNode{
			ID:                  p.ID,
			Name:                p.FullName(),
			FirstName:           p.FirstName,
			LastName:            p.LastName,
			Gender:              string(p.Gender),
			BirthYear:           p.BirthDate.Year(),
			DeathYear:           deathYear,
			IsAlive:             p.IsAlive(),
			Generation:          p.Generation,
			MaritalStatus:       string(p.MaritalStatus),
			MarriageAge:         p.MarriageAge,
			NumberOfChildren:    p.NumberOfChildren,
			Education:           string(p.Education),
			Employment:          string(p.Employment),
			AlcoholConsumption:  p.Health.AlcoholConsumption,
			TobaccoUse:          p.Health.TobaccoUse,
			BornOutsideMarriage: p.BornOutsideMarriage,
			IsSingleParent:      p.IsSingleParent,
			Underweight:         p.Underweight,
			Residence:           string(p.Residence),
			GDPPerCapita:        p.GDPPerCapita,
			WealthIndex:         p.WealthIndex,
			FamilyWealth:        p.FamilyWealth,
			IsRich:              p.IsRich,
			Country:             p.BirthCountry,
			CurrentCountry:      p.CurrentCountry,
		}
		data.Nodes = append(data.Nodes, node)

		
		if p.IsAlive() {
			data.Stats.LivingPersons++
		} else {
			data.Stats.DeceasedPersons++
		}

		
		if p.Gender == model.Male {
			data.Stats.MaleCount++
		} else {
			data.Stats.FemaleCount++
		}

		
		switch p.MaritalStatus {
		case model.Single:
			data.Stats.SingleCount++
		case model.Married, model.Remarried:
			data.Stats.MarriedCount++
		case model.Divorced:
			data.Stats.DivorceCount++
		}

		
		if p.Education == model.Tertiary {
			data.Stats.TertiaryEducation++
		}

		
		if p.Employment == model.Employed {
			data.Stats.EmployedCount++
		}

		
		if p.BornOutsideMarriage {
			data.Stats.BirthsOutsideMarriage++
		}

		age := p.AgeAtDeath()
		if age < 0 {
			age = referenceYear - p.BirthDate.Year()
			if age < 0 {
				age = 0
			}
		}
		totalAge += float64(age)
		ageCount++
		if age > oldestAge {
			oldestAge = age
		}
	}

	
	for _, p := range persons {
		if p.FatherID != nil {
			data.Edges = append(data.Edges, VisualizationEdge{
				Source: *p.FatherID,
				Target: p.ID,
				Type:   "parent",
			})
		}
		if p.MotherID != nil {
			data.Edges = append(data.Edges, VisualizationEdge{
				Source: *p.MotherID,
				Target: p.ID,
				Type:   "parent",
			})
		}
	}

	
	seen := make(map[string]bool)
	for _, p := range persons {
		for _, spouseID := range p.SpouseIDs {
			key := p.ID + "-" + spouseID
			reverseKey := spouseID + "-" + p.ID
			if !seen[key] && !seen[reverseKey] {
				data.Edges = append(data.Edges, VisualizationEdge{
					Source: p.ID,
					Target: spouseID,
					Type:   "spouse",
				})
				seen[key] = true
			}
		}
	}

	
	data.Stats.TotalPersons = tree.PersonCount()
	data.Stats.TotalFamilies = tree.FamilyCount()
	if ageCount > 0 {
		data.Stats.AverageAge = totalAge / float64(ageCount)
	}
	data.Stats.OldestPerson = oldestAge

	
	for _, f := range tree.GetAllFamilies() {
		data.Stats.TotalChildren += f.ChildCount()
	}
	if tree.FamilyCount() > 0 {
		data.Stats.AverageChildren = float64(data.Stats.TotalChildren) / float64(tree.FamilyCount())
	}

	return data
}
