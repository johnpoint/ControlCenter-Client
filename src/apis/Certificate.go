package apis

import (
	"io"
	"os"
)

func SyncCer() bool {
	sslPath := "/web/ssl/"
	data := getData()
	for i := 0; i < len(data.Certificates); i++ {
		if _, err := os.Stat(sslPath + data.Certificates[i].Domain); os.IsNotExist(err) {
			os.Mkdir(sslPath+data.Certificates[i].Domain, 0777)
		}
		fc, _ := os.Create(sslPath + data.Certificates[i].Domain + "/" + data.Certificates[i].Domain + ".fc") // .fc as fullchain
		defer fc.Close()
		_, err := io.WriteString(fc, data.Certificates[i].FullChain)
		if err != nil {
			panic(err)
		}
		key, _ := os.Create(sslPath + data.Certificates[i].Domain + "/" + data.Certificates[i].Domain + ".key")
		_, err = io.WriteString(key, data.Certificates[i].Key)
		if err != nil {
			panic(err)
		}
	}
	return true
}
