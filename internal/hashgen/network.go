package hashgen

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"

	"github.com/ilyakaznacheev/rand-hashing/internal/pkg/config"
	"github.com/ilyakaznacheev/rand-hashing/internal/pkg/types"
)

const (
	redisLisrKey = "randhash"
)

// StartGeneration generates SHA-3 sums, saves them into Redis and sends results via WebSocket
//
// key - base key
//
// n - number of keys, randomly generated by base key
func StartGeneration(confPath, key string, n int) error {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "hashgen",
	}

	// read config file
	conf, err := config.ReadConfig(confPath)
	if err != nil {
		return err
	}

	// connect to client
	ws, err := newWSConnection(u.String())
	if err != nil {
		return err
	}
	defer ws.Close()

	// create redis handler
	rh := newRedisHandler(*conf)

	// handle keyboard interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// generage sums
	sumChan := startHashing(key, n)

	log.Println("start generation")

	// send results via ws
	for {
		select {
		case <-time.After(3 * time.Second):
			if kh, ok := <-sumChan; !ok {
				log.Println("work done")
				os.Exit(0)
			} else {
				strHash := fmt.Sprintf("%x", kh.hash)
				// save into Redis list
				go rh.saveToRedis(kh.key, strHash)
				// send via ws
				err = ws.sendMessage(kh.key, strHash)
				log.Println("ws error: ", err)
			}
		case <-interrupt:
			log.Println("keyboard interrupt")
			os.Exit(0)
		}
	}
}

type wsConnection struct {
	*websocket.Conn
}

func newWSConnection(url string) (*wsConnection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return &wsConnection{conn}, nil
}

func (c *wsConnection) sendMessage(key, hash string) error {
	msg := types.Message{
		Number: key,
		Hash:   hash,
	}

	err := c.WriteJSON(msg)
	if err != nil {
		return err
	}
	return nil
}

type redisHandler struct {
	c *redis.Client
}

func newRedisHandler(conf config.Config) *redisHandler {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	return &redisHandler{client}
}

func (r *redisHandler) saveToRedis(key, hash string) {
	msg := types.Message{
		Number: key,
		Hash:   hash,
	}
	jsonMsg, _ := json.Marshal(&msg)

	r.c.RPush(redisLisrKey, string(jsonMsg))
}
