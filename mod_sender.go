package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ajankovic/smpp"
	"github.com/ajankovic/smpp/pdu"
	"hash/fnv"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
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

func hashStr() string {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	rnd := rand.Int63()
	str := fmt.Sprintf("%d.%d", now.Unix(), rnd)
	h := fnv.New32a()
	_, err := h.Write([]byte(str))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum32())
}

func sendSMS(sm *pdu.SubmitSm, ctx *smpp.Context, sid string, pwd string) {
	message := sm.ShortMessage

	// idx := fmt.Sprintf("%s %s", ctx.SessionID(), sm.DestinationAddr)
	//	message = messageOrEmpty(idx, message)
	//	if message == "" {
	//		return
	//	}

	if sm.DataCoding == 8 {
		message = UCS2Decode(message)
	}
	url := cfg.REST.Url
	payload := Payload{
		Source:      sm.SourceAddr,
		Destination: sm.DestinationAddr,
		Priority:    sm.PriorityFlag,
		RemoteAddr:  ctx.RemoteAddr(),
		Message:     message,
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

	var msgId = fmt.Sprintf("%s_%d", hashStr(), resp.StatusCode)
	resBody, err := io.ReadAll(resp.Body)
	log.Printf("%v\n", string(resBody))

	respSm := sm.Response(msgId)
	if err := ctx.Respond(respSm, pdu.StatusOK); err != nil {
		log.Printf("Server can't respond to the submit_sm request: %+v", err)
	}
}
