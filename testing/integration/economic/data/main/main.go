package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Wrapper struct {
	Data []EconomicData `json:"data"`
}

type EconomicData struct {
	Date   string `json:"date"`
	Value  string `json:"value"`
	Change string `json:"change"`
}

const sqlPath = "../../../sql"
const inputFilePath = "../"

func main() {
	fmt.Println("Converting test JSON Data file -> SQL files inserting data")

	err := os.Mkdir(sqlPath, 0755)
	if err != nil {
		if !strings.Contains(err.Error(), "exists") {
			panic(err)
		}
	}

	files, err := ioutil.ReadDir("../")
	if err != nil {
		panic(err)
	}

	uberFile, err := os.Create(sqlPath + "/" + "000010_load_economic_data.up.sql")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			fileData, err := ioutil.ReadFile(inputFilePath + f.Name())
			if err != nil {
				panic(err)
			}
			econData := Wrapper{}
			err = json.Unmarshal(fileData, &econData)
			if err != nil {
				panic(err)
			}

			tableName := strings.Replace(f.Name(), ".json", "", 1)
			filePath := sqlPath + "/" + strings.Replace(f.Name(), ".json", ".sql", 1)
			fmt.Println("Creating SQL file " + filePath)

			f, err := os.Create(filePath)
			if err != nil {
				panic(err)
			}
			for _, d := range econData.Data {
				if d.Value == "." {
					continue
				}
				date, err := time.Parse(time.RFC3339, d.Date)
				if err != nil {
					date2, err2 := time.Parse("2006-01-02", d.Date)
					if err2 != nil {
						panic(err2)
					}
					date = date2
				}
				insert := fmt.Sprintf("INSERT INTO %s (time, value) VALUES ('%s', %s);\n", tableName, date.Format(time.RFC3339), d.Value)

				fmt.Fprint(f, insert)
				fmt.Fprint(uberFile, insert)
			}
			fmt.Fprint(uberFile, "\n")
		}
	}
}
