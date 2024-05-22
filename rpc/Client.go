package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	inProgressUrl   = "wss://api.hypervolt.co.uk/ws/charger/%d/session/in-progress"
	syncUrl         = "wss://api.hypervolt.co.uk/ws/charger/%d/sync"
	responseTimeout = 30 * time.Second
)

type Client struct {
	syncConnection          *websocket.Conn
	token                   string
	syncReceiver            responseReceiver
	syncReceiverLoopRunning bool
	responseMap             map[string][]byte
}

func NewClient(chargerID int) (*Client, error) {
	headers := http.Header{}
	headers.Add("Origin", "https://hypervolt.co.uk")
	headers.Add("Host", "api.hypervolt.co.uk")
	headers.Add("User-Agent", "home-assistant-hypervolt-charger/2.2.2")
	syncConnection, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf(syncUrl, chargerID), headers)
	if err != nil {
		return nil, err
	}
	responseMap := make(map[string][]byte)
	return &Client{
		syncConnection: syncConnection,
		syncReceiver: responseReceiver{
			messageBuffer: make(chan RawMessage, 200),
			responseChan:  make(chan RawMessage),
			updateChan:    make(chan RawMessage),
			connection:    syncConnection,
			responseMap:   responseMap,
		},
		responseMap: responseMap,
	}, nil
}

func (c *Client) StartResponseLoop() {
	c.syncReceiverLoopRunning = true
	go func() {
		err := c.syncReceiver.receiverLoop()
		if err != nil {
			fmt.Println("Error in receiver loop:", err)
		}
		time.Sleep(2 * time.Second)
	}()
	go func() {
		err := c.syncReceiver.dispatcherLoop()
		if err != nil {
			fmt.Println("Error in dispatcher loop:", err)
		}
		time.Sleep(2 * time.Second)
	}()
	go func() {
		err := c.syncReceiver.responseLoop()
		if err != nil {
			fmt.Println("Error in response loop:", err)
		}
		time.Sleep(2 * time.Second)
	}()
	go func() {
		err := c.syncReceiver.updateLoop()
		if err != nil {
			fmt.Println("Error in response loop:", err)
		}
		time.Sleep(2 * time.Second)
	}()
}

func (c *Client) Close() error {
	c.syncReceiverLoopRunning = false
	return c.syncConnection.Close()
}

func (c *Client) sendMessage(msg any) error {
	if !c.syncReceiverLoopRunning {
		return errors.New("receiver loop is not running")
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.syncConnection.WriteMessage(websocket.TextMessage, msgBytes)
}

type IDGetter interface {
	GetID() string
}

func (c *Client) sendMessageAndWaitForResponse(msg IDGetter) ([]byte, error) {
	err := c.sendMessage(msg)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), responseTimeout)
	defer cancel()
	return c.waitForResponse(ctx, msg.GetID())
}

func (c *Client) waitForResponse(ctx context.Context, id string) ([]byte, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("operation timed out")
		default:
			c.syncReceiver.mutex.Lock()
			val, ok := c.responseMap[id]
			c.syncReceiver.mutex.Unlock()
			if ok {
				return val, nil
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}
