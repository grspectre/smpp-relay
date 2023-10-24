package main

import (
	"bytes"
	"encoding/json"
	"github.com/ajankovic/smpp"
	"github.com/ajankovic/smpp/pdu"
	"io"
	"log"
	"net/http"
)

type Payload struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Priority    int    `json:"priority"`
	RemoteAddr  string `json:"remoteAddr"`
	Message     string `json:"message"`
	User        string `json:"user"`
	Password    string `json:"password"`
}

func sendSMS(sm *pdu.SubmitSm, ctx *smpp.Context, sid string, pwd string) {
	log.Printf("Data coding: %v", rune(sm.DataCoding))
	url := cfg.REST.Url
	payload := Payload{
		Source:      sm.SourceAddr,
		Destination: sm.DestinationAddr,
		Priority:    sm.PriorityFlag,
		RemoteAddr:  ctx.RemoteAddr(),
		Message:     UCS2Decode(sm.ShortMessage),
		User:        sid,
		Password:    pwd,
	}

	jsonData, err := json.Marshal(payload)
	log.Printf("Data: %v", string(jsonData))
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}

	// Create the HTTP client
	client := &http.Client{}

	// Create the POST request with the JSON payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(cfg.REST.HeaderKey, cfg.REST.Token)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		var data []byte
		_, err2 := Body.Read(data)
		if err2 != nil {
			log.Fatalf("read body error: %v", err)
			return
		}
		log.Printf("Body: %v", data)
		err := Body.Close()
		if err != nil {
			log.Fatalf("request error: %v", err)
		}
	}(resp.Body)
}
