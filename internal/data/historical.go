package data

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type HistoricalRecord struct {
	Entity string
	Code   string
	Year   int
	Value  float64
}

type HistoricalDataset struct {
	Name    string
	Records []HistoricalRecord
	ByCode  map[string][]HistoricalRecord
	ByYear  map[int][]HistoricalRecord
}

type HistoricalData struct {
	FertilityRate         *HistoricalDataset
	MarriageAgeWomen      *HistoricalDataset
	DivorceRate           *HistoricalDataset
	YouthMortality        *HistoricalDataset
	BirthsOutsideMarriage *HistoricalDataset
	MarriageRate          *HistoricalDataset
	SingleParentShare     *HistoricalDataset
	UrbanPopulationShare  *HistoricalDataset
}

func LoadHistoricalCSV(filepath string) (*HistoricalDataset, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	reader := csv.NewReader(bytes.NewReader(data))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parsing CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("no data rows in file")
	}

	dataset := &HistoricalDataset{
		Name:   filepath,
		ByCode: make(map[string][]HistoricalRecord),
		ByYear: make(map[int][]HistoricalRecord),
	}

	for _, row := range records[1:] {
		if len(row) < 4 {
			continue
		}

		year, err := strconv.Atoi(strings.TrimSpace(row[2]))
		if err != nil {
			continue
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
		if err != nil {
			continue
		}

		record := HistoricalRecord{
			Entity: strings.TrimSpace(row[0]),
			Code:   strings.TrimSpace(row[1]),
			Year:   year,
			Value:  value,
		}

		dataset.Records = append(dataset.Records, record)
		dataset.ByCode[record.Code] = append(dataset.ByCode[record.Code], record)
		dataset.ByYear[year] = append(dataset.ByYear[year], record)
	}

	for code := range dataset.ByCode {
		sort.Slice(dataset.ByCode[code], func(i, j int) bool {
			return dataset.ByCode[code][i].Year < dataset.ByCode[code][j].Year
		})
	}

	return dataset, nil
}

func LoadUrbanShareCSV(filepath string) (*HistoricalDataset, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	reader := csv.NewReader(bytes.NewReader(data))
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parsing CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("no data rows in file")
	}

	dataset := &HistoricalDataset{
		Name:   filepath,
		ByCode: make(map[string][]HistoricalRecord),
		ByYear: make(map[int][]HistoricalRecord),
	}

	for _, row := range records[1:] {
		if len(row) < 5 {
			continue
		}

		year, err := strconv.Atoi(strings.TrimSpace(row[2]))
		if err != nil {
			continue
		}

		urban, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
		if err != nil {
			continue
		}
		rural, err := strconv.ParseFloat(strings.TrimSpace(row[4]), 64)
		if err != nil {
			continue
		}
		total := urban + rural
		if total <= 0 {
			continue
		}

		record := HistoricalRecord{
			Entity: strings.TrimSpace(row[0]),
			Code:   strings.TrimSpace(row[1]),
			Year:   year,
			Value:  urban / total,
		}

		dataset.Records = append(dataset.Records, record)
		dataset.ByCode[record.Code] = append(dataset.ByCode[record.Code], record)
		dataset.ByYear[year] = append(dataset.ByYear[year], record)
	}

	for code := range dataset.ByCode {
		sort.Slice(dataset.ByCode[code], func(i, j int) bool {
			return dataset.ByCode[code][i].Year < dataset.ByCode[code][j].Year
		})
	}

	return dataset, nil
}

func (d *HistoricalDataset) GetValue(code string, year int) (float64, bool) {
	records, ok := d.ByCode[code]
	if !ok || len(records) == 0 {
		return 0, false
	}

	for _, r := range records {
		if r.Year == year {
			return r.Value, true
		}
	}

	var before, after *HistoricalRecord
	for i := range records {
		if records[i].Year < year {
			before = &records[i]
		} else if records[i].Year > year && after == nil {
			after = &records[i]
			break
		}
	}

	if before != nil && after != nil {

		ratio := float64(year-before.Year) / float64(after.Year-before.Year)
		return before.Value + ratio*(after.Value-before.Value), true
	} else if before != nil {

		return before.Value, true
	} else if after != nil {

		return after.Value, true
	}

	return 0, false
}

func (d *HistoricalDataset) GetLatestValue(code string) (float64, int, bool) {
	records, ok := d.ByCode[code]
	if !ok || len(records) == 0 {
		return 0, 0, false
	}

	latest := records[len(records)-1]
	return latest.Value, latest.Year, true
}

func (d *HistoricalDataset) GetValueOrDefault(code string, year int, defaultVal float64) float64 {
	if val, ok := d.GetValue(code, year); ok {
		return val
	}
	return defaultVal
}

func (d *HistoricalDataset) GetAvailableYearRange() (int, int) {
	if len(d.Records) == 0 {
		return 0, 0
	}

	minYear, maxYear := d.Records[0].Year, d.Records[0].Year
	for _, r := range d.Records {
		if r.Year < minYear {
			minYear = r.Year
		}
		if r.Year > maxYear {
			maxYear = r.Year
		}
	}
	return minYear, maxYear
}

func (d *HistoricalDataset) GetAvailableCountries() []string {
	codes := make([]string, 0, len(d.ByCode))
	for code := range d.ByCode {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes
}

func LoadHistoricalData(dataDir string) (*HistoricalData, error) {
	h := &HistoricalData{}
	var err error

	h.FertilityRate, err = LoadHistoricalCSV(filepath.Join(dataDir, "children-born-per-woman.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading fertility rate: %w", err)
	}

	h.MarriageAgeWomen, err = LoadHistoricalCSV(filepath.Join(dataDir, "age-at-marriage-women.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading marriage age: %w", err)
	}

	h.DivorceRate, err = LoadHistoricalCSV(filepath.Join(dataDir, "divorces-per-1000-people.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading divorce rate: %w", err)
	}

	h.YouthMortality, err = LoadHistoricalCSV(filepath.Join(dataDir, "youth-mortality-rate.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading youth mortality: %w", err)
	}

	h.BirthsOutsideMarriage, err = LoadHistoricalCSV(filepath.Join(dataDir, "share-of-births-outside-marriage.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading births outside marriage: %w", err)
	}

	h.MarriageRate, err = LoadHistoricalCSV(filepath.Join(dataDir, "marriage-rate-per-1000-inhabitants.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading marriage rate: %w", err)
	}

	h.SingleParentShare, err = LoadHistoricalCSV(filepath.Join(dataDir, "share-of-single-parent-households.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading single parent share: %w", err)
	}

	h.UrbanPopulationShare, err = LoadUrbanShareCSV(filepath.Join(dataDir, "urban-and-rural-population.csv"))
	if err != nil {
		return nil, fmt.Errorf("loading urban population share: %w", err)
	}

	return h, nil
}

var slugToISO3 = map[string]string{
	"afghanistan": "AFG", "albania": "ALB", "algeria": "DZA", "argentina": "ARG",
	"armenia": "ARM", "australia": "AUS", "austria": "AUT", "azerbaijan": "AZE",
	"bangladesh": "BGD", "belarus": "BLR", "belgium": "BEL", "bolivia": "BOL",
	"bosnia-and-herzegovina": "BIH", "brazil": "BRA", "bulgaria": "BGR",
	"cambodia": "KHM", "cameroon": "CMR", "canada": "CAN", "chile": "CHL",
	"china": "CHN", "colombia": "COL", "costa-rica": "CRI", "croatia": "HRV",
	"cuba": "CUB", "czech-republic": "CZE", "denmark": "DNK", "dominican-republic": "DOM",
	"ecuador": "ECU", "egypt": "EGY", "el-salvador": "SLV", "estonia": "EST",
	"ethiopia": "ETH", "finland": "FIN", "france": "FRA", "georgia": "GEO",
	"germany": "DEU", "ghana": "GHA", "greece": "GRC", "guatemala": "GTM",
	"honduras": "HND", "hungary": "HUN", "iceland": "ISL", "india": "IND",
	"indonesia": "IDN", "iran": "IRN", "iraq": "IRQ", "ireland": "IRL",
	"israel": "ISR", "italy": "ITA", "jamaica": "JAM", "japan": "JPN",
	"jordan": "JOR", "kazakhstan": "KAZ", "kenya": "KEN", "korea": "KOR",
	"kuwait": "KWT", "latvia": "LVA", "lebanon": "LBN", "lithuania": "LTU",
	"luxembourg": "LUX", "malaysia": "MYS", "malta": "MLT", "mexico": "MEX",
	"moldova": "MDA", "mongolia": "MNG", "montenegro": "MNE", "morocco": "MAR",
	"nepal": "NPL", "netherlands": "NLD", "new-zealand": "NZL", "nicaragua": "NIC",
	"nigeria": "NGA", "norway": "NOR", "pakistan": "PAK", "panama": "PAN",
	"paraguay": "PRY", "peru": "PER", "philippines": "PHL", "poland": "POL",
	"portugal": "PRT", "romania": "ROU", "russia": "RUS", "russian-federation": "RUS",
	"saudi-arabia": "SAU", "senegal": "SEN", "serbia": "SRB", "singapore": "SGP",
	"slovakia": "SVK", "slovenia": "SVN", "south-africa": "ZAF", "spain": "ESP",
	"sri-lanka": "LKA", "sweden": "SWE", "switzerland": "CHE", "taiwan": "TWN",
	"tanzania": "TZA", "thailand": "THA", "tunisia": "TUN", "turkey": "TUR",
	"ukraine": "UKR", "united-arab-emirates": "ARE", "united-kingdom": "GBR",
	"united-states": "USA", "uruguay": "URY", "uzbekistan": "UZB",
	"venezuela": "VEN", "vietnam": "VNM", "yemen": "YEM", "zambia": "ZMB",
	"zimbabwe": "ZWE", "faroe-islands": "FRO",
}

func GetISO3FromSlug(slug string) string {
	if code, ok := slugToISO3[slug]; ok {
		return code
	}
	return ""
}
