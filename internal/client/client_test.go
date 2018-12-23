package client

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/ilyakaznacheev/rand-hashing/internal/pkg/types"
)

func TestClient(t *testing.T) {
	tests := []struct {
		msg []types.Message
	}{
		{
			msg: []types.Message{
				{
					Number: "34230604",
					Hash:   "7f0b34d6070a97be239234bce620cbe8cee3bbc174e412dc501dbb4ee9122605",
				},
			},
		},
		{
			msg: []types.Message{
				{
					Number: "34230604",
					Hash:   "7f0b34d6070a97be239234bce620cbe8cee3bbc174e412dc501dbb4ee9122605",
				},
				{
					Number: "34230621",
					Hash:   "3ee49d26fa38c09fd439b51ce857f7f46c40e595c40b1b0e3fffed6e890a1139",
				},
				{
					Number: "34230359",
					Hash:   "e707a01fc84fcf31d9f50dd483322e13865cb8dc6b14cf6bb14babe5425a2e3d",
				},
			},
		},
		{
			msg: []types.Message{
				{
					Number: "1111",
					Hash:   "aaaa",
				},
				{
					Number: "2222",
					Hash:   "bbbb",
				},
				{
					Number: "3333",
					Hash:   "cccc",
				},
				{
					Number: "4444",
					Hash:   "dddd",
				},
				{
					Number: "5555",
					Hash:   "eeee",
				},
			},
		},
	}

	for _, tt := range tests {
		// create request handler with local writer
		w := &bytes.Buffer{}
		msg := make(chan receivedMessage)
		stop := make(chan struct{})
		h := &Handler{
			msgChan: msg,
			out:     w,
			stop:    stop,
			getID:   func(int) string { return "" },
		}

		// serve incoming messages
		// and close printing when connection will be closed
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ReadMessage(w, r)
			h.StopPrinting()
		}))

		// create test client
		u := "ws" + strings.TrimPrefix(s.URL, "http")
		c, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Fatal(err)
		}

		// send messages and close connection
		go func() {
			for _, msg := range tt.msg {
				c.WriteJSON(&msg)
			}
			c.Close()
		}()

		// print all messages
		h.PrintMessages()

		s.Close()

		act := w.String()

		for _, msg := range tt.msg {
			text := msg.Number + " : " + msg.Hash
			if !strings.Contains(act, text) {
				t.Errorf("console output doesn't contain message %s", text)
			}
		}
	}
}

func TestNewHandler(t *testing.T) {
	h := newHandler()
	if h == nil {
		t.Error("request handler creation failed")
	}
}
