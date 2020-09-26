package apis

import (
	"ControlCenter-Client/src/model"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func GetUpdate() bool {
	data := getData()
	url := data.Base.PollAddress + "/server/update/" + data.Base.Token
	method := "GET"
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		log.Print("状态获取失败! 请检查服务端状态")
		log.Print(err)
		return false
	}
	if res != nil {
		cdata := getData()
		log.Print(":: Get update from " + data.Base.PollAddress)
		decoder := json.NewDecoder(res.Body)
		Getdata := model.UpdateInfo{}
		err := decoder.Decode(&Getdata)
		if err != nil {
			log.Print("Error:", err)
		}
		data.Certificates = Getdata.Certificates
		newdata := Getdata.ConfFile
		newdata2 := newdata
		for _, i := range cdata.ConfFile {
			exist := false
			for _, j := range newdata2 {
				if i.ID == j.ID && i.Name == j.Name && i.Path == j.Path {
					exist = true
					break
				}
			}
			if !exist {
				i.Deleted = true
				newdata = append(newdata, i)
			}
		}
		data.ConfFile = newdata
		file, _ := os.Create("data.json")
		defer file.Close()
		databy, _ := json.Marshal(data)
		io.WriteString(file, string(databy))
		log.Print("OK!")
		return true
	}
	return false
}
