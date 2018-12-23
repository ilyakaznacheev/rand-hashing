package client

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	"github.com/ilyakaznacheev/rand-hashing/internal/pkg/types"
)

const (
	colorGreen   = "\033[32m"
	colorRed     = "\033[31m"
	colorBlue    = "\033[34m"
	colorDefault = "\033[0m"
)

// receivedMessage is a message structure with generator ID
type receivedMessage struct {
	types.Message
	id string
}

// Handler handles ws connection and input messages in JSON format
type Handler struct {
	msgChan chan receivedMessage
	out     io.Writer
	stop    chan struct{}
	getID   func(int) string
}

// newHandler creates new ws handler, that prints messages into console
func newHandler() *Handler {
	rand.Seed(time.Now().UnixNano())

	msg := make(chan receivedMessage)
	stop := make(chan struct{}, 1)

	return &Handler{
		msgChan: msg,
		out:     os.Stdout,
		stop:    stop,
		getID:   RandomID,
	}
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

	// generate random ID for each connected generator
	id := h.getID(8)

	// notify that new generator connected
	fmt.Fprintf(h.out, "%s%s connected%s\n", colorGreen, id, colorDefault)

	for {
		msg := types.Message{}
		err = c.ReadJSON(&msg)
		if err != nil {
			// notify that generator disconnected
			fmt.Fprintf(h.out, "%s%s disconnected%s\n", colorRed, id, colorDefault)
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

MESSAGE:
	for {
		select {
		case msg := <-h.msgChan:
			fmt.Fprintf(h.out, "%s%s >%s %s : %s\n", colorBlue, msg.id, colorDefault, msg.Number, msg.Hash)

		case <-interrupt:
			fmt.Fprintln(h.out, "\nexiting client")
			os.Exit(0)

		case <-h.stop:
			break MESSAGE
		}
	}
}

// StopPrinting terminates message prinring
func (h *Handler) StopPrinting() {
	h.stop <- struct{}{}
}

// Start runs client on default localhost:8080
func Start() {
	h := newHandler()
	http.HandleFunc("/hashgen", h.ReadMessage)

	go h.PrintMessages()
	fmt.Println("starting client")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
