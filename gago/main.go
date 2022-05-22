package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var client *whatsmeow.Client

func sendmsg(message string, receiver types.JID) error {
	_, err := client.SendMessage(receiver, "", &waProto.Message{
		Conversation: proto.String(message),
	})

	return err
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		msg := strings.ToLower(v.Message.GetConversation())
		switch msg {
		case "oi":
			err := sendmsg("oi eu sou um bot!", v.Info.Chat)
			if err != nil {
				fmt.Println(err.Error())
			}
		case "número", "numero", "aleatório":
			msg := fmt.Sprintf("aqui está um número aleatório: %d", rand.Intn(100))
			err := sendmsg(msg, v.Info.Chat)
			if err != nil {
				fmt.Println(err.Error())
			}
		case "teste":
			err := sendmsg("testado!", v.Info.Chat)
			if err != nil {
				fmt.Println(err.Error())
			}
		case "ping":
			err := sendmsg("pong!", v.Info.Chat)
			if err != nil {
				fmt.Println(err.Error())
			}
		default:
			err := sendmsg("não reconheço oque digitas\ntente alguma das opções:\n - oi\n - numero\n - test\n - ping", v.Info.Chat)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func main() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
