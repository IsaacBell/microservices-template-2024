package conf

import (
	"encoding/csv"
	"log"
	"os"
)

func ConfigDir() string {
	_, err := os.Stat("./configs")
	if os.IsNotExist(err) {
		return "../../configs/"
	}

	return "./configs/"
}

func CategoriesDir() string {
	return ConfigDir() + "transactions-personal-finance-category-taxonomy.csv"
}

func PersonalFinanceCategories() [][]string {
	// Skip the header row
	return readCsvFile(CategoriesDir())[1:]
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
