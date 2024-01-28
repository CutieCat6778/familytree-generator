package lib

import (
	"familytree-gen/dev/assets"
	"familytree-gen/dev/types"
	"math/rand"
)

func GenerateAPerson(_country *string, _sex *bool) types.Person {
	var country string
	var sex bool

	if _country == nil {
		country = getRandomCountry()
	} else {
		country = *_country
	}

	if _sex == nil {
		sex = rand.Float32() < 0.5
	} else {
		sex = *_sex
	}

	forename := randomForename(country, sex)
	surname := randomSurname(country)
}

func randomSurname(country string) string {
	surnames := assets.SurnamesData[country][rand.Intn(len(assets.SurnamesData[country]))]

	return surnames.Localized[0]
}

func randomForename(country string, sex bool) string {
	names := assets.ForenamesData[country][rand.Intn(len(assets.ForenamesData[country]))].Names

	var filteredNames []types.Forenames

	for _, name := range names {
		condition := name.Gender == "M"
		if condition == sex {
			filteredNames = append(filteredNames, name)
		}
	}

	randName := filteredNames[rand.Intn(len(filteredNames))].Localized

	return randName[0]
}

func getRandomCountry() string {
	values := make([]string, 0, len(assets.CountriesCodeData))
	for _, v := range assets.CountriesCodeData {
		values = append(values, v)
	}

	return values[rand.Intn(len(values))]
}
