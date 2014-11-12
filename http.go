package logsauce

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func SendLogLine(destinationHost, line, filepath, filename, token string) {
	logline := LogLine{Line: line, Filename: filename, Filepath: filepath}

	client := &http.Client{}
	b, err := json.Marshal(&logline)
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest("POST", "http://"+destinationHost+"/logs", bytes.NewBuffer(b))
	req.Header.Add("token", token)

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		log.Println(resp.Status)
	}

}
