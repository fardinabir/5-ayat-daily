package db

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"one-minute-quran/models"
	"os"
	"strconv"
	"testing"
)

func TestArabic_InsertData(t *testing.T) {
	viper.SetConfigFile("../config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	cc := ConnectDB()

	// Open the CSV file
	file, err := os.Open("/Users/pathaoltd/Documents/QuranAyats/QuranAllVerseSheet1.csv")
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Initialize a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV records:", err)
		return
	}

	// Open a database connection (SQLite in this example)
	db := cc

	// Auto-migrate the schema (create the table if it doesn't exist)
	db.AutoMigrate(&models.Ayah{})

	// Iterate through CSV records and write to the database
	for i, record := range records {
		// Assuming the CSV has three columns: ID, Name, Value
		id, _ := strconv.ParseUint(record[0], 10, 64)
		suraNo, _ := strconv.ParseInt(record[1], 10, 64)
		verseNo, _ := strconv.ParseInt(record[2], 10, 64)
		arabic := record[3]
		bangla := record[4]
		english := record[5]
		//value, _ := strconv.Atoi(strings.TrimSpace(record[2]))

		// Create a new instance of your model
		data := models.Ayah{
			SuraNo:          int(suraNo),
			VerseNo:         int(verseNo),
			AyahTextArabic:  arabic,
			AyahTextBangla:  bangla,
			AyahTextEnglish: english,
		}

		data.ID = uint(id)

		if i < 5 {
			fmt.Println(id, suraNo, verseNo, bangla)
		}
		// Insert or update the record in the database
		db.Create(&data)
	}

	fmt.Println("CSV data has been imported into the database.")
}
