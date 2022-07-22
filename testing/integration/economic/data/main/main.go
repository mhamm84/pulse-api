package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Wrapper struct {
	Data []EconomicData `json:"data"`
}

type EconomicData struct {
	Date   time.Time `json:"date"`
	Value  string    `json:"value"`
	Change string    `json:"change"`
}

const sqlPath = "./sql"
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
			filePath := sqlPath + "/" + strings.Replace(f.Name(), ".json", ".sql", 1)
			fmt.Println("Creating SQL file " + filePath)
			f, err := os.Create(filePath)
			if err != nil {
				panic(err)
			}
			tableName := strings.Replace(f.Name(), ".json", "", 1)
			for _, d := range econData.Data {
				fmt.Fprintf(f, "INSERT INTO %s (time, value) VALUES ('%s', %s);\n", tableName, d.Date.Format(time.RFC3339), d.Value)
			}
		}
	}
}
