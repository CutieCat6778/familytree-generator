package types

type Forenames struct {
	Gender    string   `json:"gender"`
	ID        string   `json:"id"`
	Localized []string `json:"localized"`
	Rank      int      `json:"rank"`
	Romanized []string `json:"romanized"`
}

type ForenameCountry struct {
	Names []Forenames `json:"names"`
}

type ForenameCountries map[string][]ForenameCountry

type Surname struct {
	Count     *int     `json:"count"`
	ID        string   `json:"id"`
	Localized []string `json:"localized"`
	Rank      *int     `json:"rank"`
	Romanized []string `json:"romanized"`
}

type SurnamesCountries map[string][]Surname
