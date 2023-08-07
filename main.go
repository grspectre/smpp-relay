package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/ajankovic/smpp"
	"github.com/ajankovic/smpp/pdu"
)

var (
	serverAddr string
	systemID   string
	msgID      int
)

var cfg Config

func main() {
	readConfigFile(&cfg)

	err := initLog()
	if err != nil {
		log.Fatalf("Init logger error: %v", err)
		return
	}

	serverHost := cfg.SMPP.Host + ":" + strconv.Itoa(cfg.SMPP.Port)
	flag.StringVar(&serverAddr, "addr", serverHost, "server will listen on this address.")
	flag.StringVar(&systemID, "systemid", "SMS Gateway", "descriptive server identification.")
	flag.Parse()

	sessConf := smpp.SessionConf{
		Handler: smpp.HandlerFunc(func(ctx *smpp.Context) {
			log.Printf("get command: %v", ctx.CommandID())
			switch ctx.CommandID() {
			case pdu.BindTransmitterID:
				btx, err := ctx.BindTx()
				if err != nil {
					log.Printf("Invalid PDU in context error: %+v", err)
				}
				log.Printf("Incoming connection from %s with ID: %s", ctx.RemoteAddr(), btx.SystemID)
				resp := btx.Response(systemID)
				responseStatus := pdu.StatusInvPaswd
				if btx.Password == cfg.SMPP.Password && cfg.SMPP.User == btx.SystemID {
					responseStatus = pdu.StatusOK
				}
				if err := ctx.Respond(resp, responseStatus); err != nil {
					log.Printf("Server can't respond to the Binding request: %+v", err)
				}

			case pdu.BindTransceiverID:
				btrx, err := ctx.BindTRx()
				if err != nil {
					log.Printf("Invalid PDU in context error: %+v", err)
				}
				log.Printf("Incoming connection from %s with ID: %s", ctx.RemoteAddr(), btrx.SystemID)
				resp := btrx.Response(systemID)
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					log.Printf("Server can't respond to the Binding request: %+v", err)
				}

			case pdu.SubmitSmID:
				sm, err := ctx.SubmitSm()
				if err != nil {
					log.Printf("Invalid PDU in context error: %+v", err)
				}

				go sendSMS(sm, ctx)

				resp := sm.Response(fmt.Sprintf("msgID_%d", msgID))
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					log.Printf("Server can't respond to the submit_sm request: %+v", err)
				}
				msgID++

			case pdu.UnbindID:
				unb, err := ctx.Unbind()
				if err != nil {
					log.Printf("Invalid PDU in context error: %+v", err)
				}
				resp := unb.Response()
				if err := ctx.Respond(resp, pdu.StatusOK); err != nil {
					log.Printf("Server can't respond to the submit_sm request: %+v", err)
				}
				ctx.CloseSession()
			}
		}),
	}
	srv := smpp.NewServer(serverAddr, sessConf)

	log.Printf("'%s' is listening on '%s'", systemID, serverAddr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("Serving exited with error: %+v", err)
	}
	log.Printf("Server closed")
}

func initLog() error {
	log.SetFlags(log.LstdFlags)
	f, err := os.OpenFile(cfg.SYSTEM.Log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v %v", cfg.SYSTEM.Log, err)
	}
	multi := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multi)
	return err
}
