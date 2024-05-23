package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/types"
	"github.com/anas-domesticus/hypervolt_go_sdk/state"
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
	chargerState            state.HypervoltDeviceState
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
	chargerState := state.HypervoltDeviceState{}
	c := &Client{
		syncConnection: syncConnection,
		syncReceiver: responseReceiver{
			messageBuffer: make(chan RawMessage, 25),
			responseChan:  make(chan RawMessage),
			updateChan:    make(chan RawMessage),
			connection:    syncConnection,
			responseMap:   responseMap,
		},
		responseMap:  responseMap,
		chargerState: chargerState,
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
	if msg == nil {
		return errors.New("cannot handle nil message")
	}
	if msg.Params.Charging != nil {
		c.chargerState.IsCharging = *msg.Params.Charging
		c.chargerState.CarPlugged = true
	}
	if msg.Params.TrueMilliAmps != nil {
		c.chargerState.CurrentSessionCurrentMilliamps = *msg.Params.TrueMilliAmps
	}
	if msg.Params.Voltage != nil {
		c.chargerState.CurrentSessionVoltage = *msg.Params.Voltage
	}
	if msg.Params.WattHours != nil {
		c.chargerState.SessionWatthours = *msg.Params.WattHours
	}
	if msg.Params.CarbonSavedGrams != nil {
		c.chargerState.SessionCarbonSavedGrams = *msg.Params.CarbonSavedGrams
	}
	if msg.Params.CtCurrent != nil {
		c.chargerState.CurrentSessionCtCurrent = *msg.Params.CtCurrent
	}
	if msg.Params.EvPower != nil {
		c.chargerState.EVPower = *msg.Params.EvPower
	}
	if msg.Params.GridPower != nil {
		c.chargerState.GridPower = *msg.Params.GridPower
	}
	if msg.Params.HousePower != nil {
		c.chargerState.HousePower = *msg.Params.HousePower
	}
	if msg.Params.GenerationPower != nil {
		c.chargerState.GenerationPower = *msg.Params.GenerationPower
	}
	return nil
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
