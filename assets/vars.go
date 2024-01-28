package assets

import (
	"familytree-gen/dev/types"
	"fmt"
	"log"
	"os"
)

var (
	CountriesCodeData map[string]string
	ForenamesData     types.ForenameCountries
	SurnamesData      types.SurnamesCountries
	TFRData           map[string]interface{}
)

func init() {
	var err error
	var ok bool
	data, err := readJSON("./assets/countries-code.json")
	if err != nil {
		log.Println(err)
	}
	CountriesCodeData, ok = data.(map[string]string)
	if !ok {
		log.Println("error: CountriesCodeData is not map[string]string")
	}
	data, err = readJSON("./assets/forenames.json")
	if err != nil {
		log.Println(err)
	}
	ForenamesData, ok = data.(types.ForenameCountries)
	if !ok {
		log.Println("error: ForenamesData is not ForenameCountries")
	}
	data, err = readJSON("./assets/surnames.json")
	if err != nil {
		log.Println(err)
	}
	SurnamesData, ok = data.(types.SurnamesCountries)
	if !ok {
		log.Println("error: SurnamesData is not SurnamesCountries")
	}
	data, err = readJSON("./assets/tfr.json")
	if err != nil {
		log.Println(err)
	}
	TFRData, ok = data.(map[string]interface{})
	if !ok {
		log.Println("error: TFRData is not map[string]float64")
	}
}

func readJSON(path string) (interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var data interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in readJSON", r)
			data = nil
			err = fmt.Errorf("error: %v", r)
		}
	}()
	return data, nil
}
