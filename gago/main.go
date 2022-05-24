package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
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

func sendmsg2(c *whatsmeow.Client, receiver types.JID) error {
	_, err := client.SendMessage(receiver, "", &waProto.Message{
		InteractiveMessage: &waProto.InteractiveMessage{
			Header: &waProto.Header{
				Title:              proto.String("Título"),
				Subtitle:           proto.String("subtítulo"),
				HasMediaAttachment: proto.Bool(false),
				// Media: &waProto.Header_ImageMessage{
				// 	ImageMessage: &waProto.ImageMessage{
				// 		Mimetype:  proto.String("image/png"),
				// 		StaticUrl: proto.String("https://pbs.twimg.com/profile_images/1498641868397191170/6qW2XkuI_400x400.png"),
				// 		Caption:   proto.String("cool stuff"),
				// 	},
				// },
			},
			Body: &waProto.InteractiveMessageBody{
				Text: proto.String("Corpo"),
			},
			Footer: &waProto.Footer{
				Text: proto.String("Pé"),
			},
			InteractiveMessage: &waProto.InteractiveMessage_NativeFlowMessage{
				NativeFlowMessage: &waProto.NativeFlowMessage{
					Buttons: []*waProto.NativeFlowButton{{
						Name: proto.String("1"),
						// ButtonParamsJson: proto.String("1"),
					}, {
						Name: proto.String("2"),
						// ButtonParamsJson: proto.String("1"),
					}, {
						Name: proto.String("3"),
						// ButtonParamsJson: proto.String("1"),
					}},
					// MessageParamsJson: proto.String("message params json"),
					MessageVersion: proto.Int32(26),
				},
			},
		},
	})

	return err
}

func sendlink(s, url string, receiver types.JID) error {
	_, err := client.SendMessage(receiver, "", &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String(url),
			MatchedText: proto.String(url),
			Description: proto.String("Redirecionamento para " + s),
			Title:       proto.String(s),
			PreviewType: waProto.ExtendedTextMessage_NONE.Enum(),
		},
	})
	return err
}

func sendbuttons(receiver types.JID) error {
	_, err := client.SendMessage(receiver, "", &waProto.Message{
		// Conversation: proto.String(message),
		ButtonsMessage: &waProto.ButtonsMessage{
			HeaderType: waProto.ButtonsMessage_TEXT.Enum(),
			Header: &waProto.ButtonsMessage_Text{
				Text: *proto.String("Atual Printer"),
			},
			ContentText: proto.String(`
Utilize das opções para a navegação.
`),
			FooterText: proto.String("Impressão Digital"),
			Buttons: []*waProto.Button{{
				ButtonId: proto.String("canvas"),
				ButtonText: &waProto.ButtonText{
					DisplayText: proto.String("canvas"),
				},
				Type: waProto.Button_RESPONSE.Enum(),
			}, {
				ButtonId: proto.String("etiquetas"),
				ButtonText: &waProto.ButtonText{
					DisplayText: proto.String("etiquetas"),
				},
				Type: waProto.Button_RESPONSE.Enum(),
			}, {
				ButtonId: proto.String("adesivos"),
				ButtonText: &waProto.ButtonText{
					DisplayText: proto.String("adesivos"),
				},
				Type: waProto.Button_RESPONSE.Enum(),
			}, {
				ButtonId: proto.String("banners"),
				ButtonText: &waProto.ButtonText{
					DisplayText: proto.String("banners e faixas"),
				},
				Type: waProto.Button_RESPONSE.Enum(),
			}, {
				ButtonId: proto.String("lonas"),
				ButtonText: &waProto.ButtonText{
					DisplayText: proto.String("lonas"),
				},
				Type: waProto.Button_RESPONSE.Enum(),
			}, {
				ButtonId: proto.String("pspvc"),
				ButtonText: &waProto.ButtonText{
					DisplayText: proto.String("ps/pvc adesivado"),
				},
				Type: waProto.Button_RESPONSE.Enum(),
			}, {
				ButtonId: proto.String("recorte"),
				ButtonText: &waProto.ButtonText{
					DisplayText: proto.String("recorte eletrônico"),
				},
				Type: waProto.Button_RESPONSE.Enum(),
			}},
		},
	})

	return err
}

func stringHandler(c *whatsmeow.Client, v *events.Message, s string) {
	switch strings.ToLower(s) {
	case "oi", "menu", "opcoes", "opções", "bom dia", "boa tarde", "boa noite", "olá", "salve":
		err := sendbuttons(v.Info.Chat)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "site", "homepage", "home", "link":
		url := "https://www.atualprinter.com.br/"
		err := sendlink("HomePage", url, v.Info.Chat)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "canvas", "canva":
		url := "https://www.atualprinter.com.br/canvas-g79"
		err := sendlink("Canvas", url, v.Info.Chat)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "etiquetas", "etiqueta", "etiquetar":
		url := "https://www.atualprinter.com.br/canvas-g79"
		err := sendlink("Etiquetas", url, v.Info.Chat)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "adesivos", "adesivo", "adesivar":
		url := "https://www.atualprinter.com.br/adesivos-g16"
		err := sendlink("Adesivos", url, v.Info.Chat)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "banners", "banner", "baner", "baners", "faixa", "faixas", "banners e faixas":
		url := "https://www.atualprinter.com.br/banners-e-faixas-g3"
		err := sendlink("Banners", url, v.Info.Chat)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "lona", "lonas":
		url := "https://www.atualprinter.com.br/lonas-g27"
		err := sendlink("Lonas", url, v.Info.Chat)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "ps", "pvc", "adesivado", "ps/pvc adesivado":
		url := "https://www.atualprinter.com.br/pspvc-adesivado-g4"
		err := sendlink("PS/PVC Adesivado", url, v.Info.Chat)
		if err != nil {
			fmt.Println(err.Error())
		}
	case "recorte", "eletronico", "recorte eletronico", "recorte eletrônico", "eletrônico":
		url := "https://www.atualprinter.com.br/recorte-eletronico-g78"
		err := sendlink("Recorte eletrônico", url, v.Info.Chat)
		if err != nil {
			fmt.Println(err.Error())
		}
	default:
		err := sendmsg("Não encontrei a opção selecionada *"+s+"*.", v.Info.Chat)
		if err != nil {
			fmt.Println(err.Error())
		}
		stringHandler(c, v, "menu")
	}
}

func mainEventHandler(c *whatsmeow.Client, evt interface{}) {
	has := func(v interface{}) bool { return !reflect.ValueOf(v).IsZero() }
	switch v := evt.(type) {
	case *events.Message:
		msg := v.Message
		fmt.Printf("ola %s\n", msg.String())
		switch {
		case has(msg.GetButtonsResponseMessage()):
			fmt.Printf("has %s\n", msg.GetButtonsResponseMessage())
			s := msg.GetButtonsResponseMessage().GetSelectedButtonId()
			stringHandler(c, v, s)
		default:
			fmt.Printf("default %s\n", msg.GetButtonsMessage().GetText())
			s := msg.GetConversation()
			stringHandler(c, v, s)
		}
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
	client.AddEventHandler(func(evt interface{}) {
		mainEventHandler(client, evt)
	})

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
