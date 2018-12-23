package hashgen

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"

	"github.com/ilyakaznacheev/rand-hashing/internal/pkg/types"
)

func TestNetwork(t *testing.T) {
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
		msgChan := make(chan types.Message)

		// create mock server
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{}
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()

			for {
				msg := types.Message{}
				err := c.ReadJSON(&msg)
				if err != nil {
					break
				}
				msgChan <- msg
			}
		}))

		// connect to mock server
		u := "ws" + strings.TrimPrefix(s.URL, "http")
		wsCon, err := newWSConnection(u)
		if err != nil {
			t.Error(err)
		}

		// send and receive test results
		for _, msg := range tt.msg {
			wsCon.sendMessage(msg.Number, msg.Hash)
			act := <-msgChan
			if act != msg {
				t.Errorf("wrong message %v, expected %v", act, msg)
			}
		}

		wsCon.Close()
		s.Close()
	}
}
