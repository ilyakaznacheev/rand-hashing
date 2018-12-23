package client

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"

	"github.com/ilyakaznacheev/rand-hashing/internal/pkg/types"
)

const (
	colorGreen   = "\033[32m"
	colorRed     = "\033[31m"
	colorBlue    = "\033[34m"
	colorDefault = "\033[0m"
)

type receivedMessage struct {
	types.Message
	id string
}

// Handler handles ws connection and input messages in JSON format
type Handler struct {
	msgChan chan receivedMessage
}

// newHandler creates new ws handler
func newHandler() *Handler {
	msg := make(chan receivedMessage)
	return &Handler{msg}
}

// ReadMessage handles input ws message
func (h *Handler) ReadMessage(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WS error: ", err)
		return
	}
	defer c.Close()

	id := RandomID(8)

	// notify that new generator connected
	fmt.Printf("%s%s connected%s\n", colorGreen, id, colorDefault)

	for {
		msg := types.Message{}
		err = c.ReadJSON(&msg)
		if err != nil {
			// notify that new generator disconnected
			fmt.Printf("%s%s disconnected%s\n", colorRed, id, colorDefault)
			return
		}

		// send message to cmd print
		h.msgChan <- receivedMessage{
			Message: msg,
			id:      id,
		}

	}
}

// PrintMessages prints incoming messages into terminal
func (h *Handler) PrintMessages() {
	// handle keyboard interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case msg := <-h.msgChan:
			fmt.Printf("%s%s >%s %s : %s\n", colorBlue, msg.id, colorDefault, msg.Number, msg.Hash)

		case <-interrupt:
			fmt.Println("\nexiting client")
			os.Exit(0)
		}
	}
}

// Start runs client on default localhost:8080
func Start() {
	h := newHandler()
	http.HandleFunc("/hashgen", h.ReadMessage)

	go h.PrintMessages()
	fmt.Println("starting client")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
