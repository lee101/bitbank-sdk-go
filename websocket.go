package bitbank

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type EventCallback func(msg map[string]interface{})

type WSClient struct {
	URL        string
	APIKey     string
	MaxRetries int

	conn       *websocket.Conn
	mu         sync.Mutex
	callbacks  map[string][]EventCallback
	done       chan struct{}
	reconnect  bool
}

func NewWSClient(url string, apiKey string) *WSClient {
	return &WSClient{
		URL:        url,
		APIKey:     apiKey,
		MaxRetries: 5,
		callbacks:  make(map[string][]EventCallback),
	}
}

func (c *Client) WebSocket() *WSClient {
	wsURL := "wss://bitbank.nz/ws"
	if c.baseURL != DefaultBaseURL {
		wsURL = c.baseURL + "/ws"
	}
	return NewWSClient(wsURL, c.apiKey)
}

func (w *WSClient) On(event string, cb EventCallback) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.callbacks[event] = append(w.callbacks[event], cb)
}

func (w *WSClient) Connect() error {
	w.reconnect = true
	w.done = make(chan struct{})
	go w.run()
	return nil
}

func (w *WSClient) Disconnect() {
	w.reconnect = false
	w.mu.Lock()
	if w.conn != nil {
		w.conn.Close()
	}
	w.mu.Unlock()
	if w.done != nil {
		select {
		case <-w.done:
		case <-time.After(5 * time.Second):
		}
	}
}

func (w *WSClient) Subscribe(pair string) error {
	return w.send(map[string]interface{}{"event": "subscribe", "currency_pair": pair})
}

func (w *WSClient) Unsubscribe(pair string) error {
	return w.send(map[string]interface{}{"event": "unsubscribe", "currency_pair": pair})
}

func (w *WSClient) SubscribeAll() error {
	return w.send(map[string]interface{}{"event": "subscribe_all_pairs"})
}

func (w *WSClient) SubscribeLive(pairs []string) error {
	msg := map[string]interface{}{"event": "subscribe_live"}
	if w.APIKey != "" {
		msg["secret"] = w.APIKey
	}
	if len(pairs) > 0 {
		msg["currency_pairs"] = pairs
	}
	return w.send(msg)
}

func (w *WSClient) Ping() error {
	return w.send(map[string]interface{}{"event": "ping"})
}

func (w *WSClient) send(data map[string]interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.conn == nil {
		return &WebSocketError{BitbankError{Message: "not connected"}}
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return w.conn.WriteMessage(websocket.TextMessage, b)
}

func (w *WSClient) run() {
	defer close(w.done)
	retries := 0
	for w.reconnect {
		conn, _, err := websocket.DefaultDialer.Dial(w.URL, nil)
		if err != nil {
			log.Printf("ws dial error: %v", err)
			retries++
			if retries > w.MaxRetries {
				return
			}
			time.Sleep(time.Duration(min(1<<retries, 30)) * time.Second)
			continue
		}
		w.mu.Lock()
		w.conn = conn
		w.mu.Unlock()
		retries = 0

		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("ws read error: %v", err)
				break
			}
			var msg map[string]interface{}
			if json.Unmarshal(data, &msg) != nil {
				continue
			}
			event, _ := msg["event"].(string)
			w.mu.Lock()
			for _, cb := range w.callbacks[event] {
				cb(msg)
			}
			for _, cb := range w.callbacks["*"] {
				cb(msg)
			}
			w.mu.Unlock()
		}

		w.mu.Lock()
		w.conn = nil
		w.mu.Unlock()

		if !w.reconnect {
			break
		}
		retries++
		if retries > w.MaxRetries {
			return
		}
		delay := min(1<<retries, 30)
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
