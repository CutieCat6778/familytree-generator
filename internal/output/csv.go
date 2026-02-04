package output

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/familytree-generator/internal/model"
)


func WriteCSV(tree *model.FamilyTree, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	
	header := []string{
		"id",
		"first_name",
		"last_name",
		"gender",
		"birth_date",
		"death_date",
		"birth_country",
		"current_country",
		"father_id",
		"mother_id",
		"spouse_ids",
		"children_ids",
		"generation",
		"education",
		"employment",
		"alcohol_consumption",
		"tobacco_use",
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	
	for _, person := range tree.GetAllPersons() {
		row := personToRow(person)
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("writing row for %s: %w", person.ID, err)
		}
	}

	return nil
}


func personToRow(p *model.Person) []string {
	deathDate := ""
	if p.DeathDate != nil {
		deathDate = p.DeathDate.Format("2006-01-02")
	}

	fatherID := ""
	if p.FatherID != nil {
		fatherID = *p.FatherID
	}

	motherID := ""
	if p.MotherID != nil {
		motherID = *p.MotherID
	}

	tobaccoUse := "false"
	if p.Health.TobaccoUse {
		tobaccoUse = "true"
	}

	return []string{
		p.ID,
		p.FirstName,
		p.LastName,
		string(p.Gender),
		p.BirthDate.Format("2006-01-02"),
		deathDate,
		p.BirthCountry,
		p.CurrentCountry,
		fatherID,
		motherID,
		strings.Join(p.SpouseIDs, ";"),
		strings.Join(p.ChildrenIDs, ";"),
		fmt.Sprintf("%d", p.Generation),
		string(p.Education),
		string(p.Employment),
		fmt.Sprintf("%.1f", p.Health.AlcoholConsumption),
		tobaccoUse,
	}
}


func WriteFamiliesCSV(tree *model.FamilyTree, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	
	header := []string{
		"id",
		"husband_id",
		"wife_id",
		"married_date",
		"divorce_date",
		"children_ids",
		"children_count",
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	
	for _, family := range tree.GetAllFamilies() {
		row := familyToRow(family)
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("writing row for %s: %w", family.ID, err)
		}
	}

	return nil
}


func familyToRow(f *model.Family) []string {
	husbandID := ""
	if f.HusbandID != nil {
		husbandID = *f.HusbandID
	}

	wifeID := ""
	if f.WifeID != nil {
		wifeID = *f.WifeID
	}

	divorceDate := ""
	if f.DivorceDate != nil {
		divorceDate = f.DivorceDate.Format("2006-01-02")
	}

	return []string{
		f.ID,
		husbandID,
		wifeID,
		f.MarriedDate.Format("2006-01-02"),
		divorceDate,
		strings.Join(f.ChildrenIDs, ";"),
		fmt.Sprintf("%d", len(f.ChildrenIDs)),
	}
}
