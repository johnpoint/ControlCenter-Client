package apis

import (
	"encoding/json"
	. "github.com/johnpoint/ControlCenter-Client/src/model"
	"log"
	"os"
)

const ClientVersion = "1.9.5"

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
