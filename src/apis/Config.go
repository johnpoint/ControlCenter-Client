package apis

import (
	. "ControlCenter-Client/src/model"
	"encoding/json"
	"log"
	"os"
)

const ClientVersion = "2.0.3"

func getData() Data {
	file, _ := os.Open("data.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	data := Data{}
	err := decoder.Decode(&data)
	if err != nil {
		log.Print("Error:", err)
	}
	return data
}
