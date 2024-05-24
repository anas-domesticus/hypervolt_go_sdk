package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/types"
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
	syncConnection          WebsocketWrapperIface
	token                   string
	syncReceiver            responseReceiver
	syncReceiverLoopRunning bool
	responseMap             map[string][]byte
	currentSession          types.SessionParams
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
	c, err := NewClientWithConnection(syncConnection)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewClientWithConnection(w WebsocketWrapperIface) (*Client, error) {
	responseMap := make(map[string][]byte)
	c := &Client{
		syncConnection: w,
		syncReceiver: responseReceiver{
			messageBuffer: make(chan RawMessage, 25),
			responseChan:  make(chan RawMessage),
			updateChan:    make(chan RawMessage),
			connection:    w,
			responseMap:   responseMap,
		},
		responseMap: responseMap,
	}
	c.syncReceiver.updateChargerState = c.updateChargerState
	return c, nil
}

func (c *Client) updateChargerState(msg RawMessage) error {
	val, ok := msg["method"]
	if !ok {
		return errors.New("no method in message")
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	switch val {
	case "get.session":
		var response types.GetSession
		err := json.Unmarshal(bytes, &response)
		if err != nil {
			return err
		}
		return c.handleGetSession(&response)
	}
	return nil
}

func (c *Client) handleGetSession(msg *types.GetSession) error {
	c.currentSession = msg.Params
	return nil
}

func (c *Client) GetCurrentSession() *types.SessionParams {
	return &c.currentSession
}

func (c *Client) StartResponseLoop() {
	c.syncReceiverLoopRunning = true
	c.syncReceiver.StartLoops()
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
				c.syncReceiver.mutex.Lock()
				delete(c.responseMap, id)
				c.syncReceiver.mutex.Unlock()
				return val, nil
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}
