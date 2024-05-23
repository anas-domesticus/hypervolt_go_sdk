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
	"strconv"
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
	c, err := NewClientWithConnection(syncConnection)
	if err != nil {
		return nil, err
	}
	c.chargerState.ChargerID = strconv.Itoa(chargerID)
	return c, nil
}

func NewClientWithConnection(w WebsocketWrapperIface) (*Client, error) {
	responseMap := make(map[string][]byte)
	chargerState := state.HypervoltDeviceState{}
	c := &Client{
		syncConnection: w,
		syncReceiver: responseReceiver{
			messageBuffer: make(chan RawMessage, 25),
			responseChan:  make(chan RawMessage),
			updateChan:    make(chan RawMessage),
			connection:    w,
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
	case "sync.apply":
		var response types.SyncApply
		err := json.Unmarshal(bytes, &response)
		if err != nil {
			return err
		}
		return c.handleSyncApply(&response)
	}
	return nil
}

func (c *Client) handleGetSession(msg *types.GetSession) error {
	if msg == nil {
		return errors.New("cannot handle nil message")
	}
	if msg.Params.Charging != nil {
		c.chargerState.IsCharging = *msg.Params.Charging
		if *msg.Params.Charging {
			c.chargerState.CarPlugged = true
		}
	}
	if msg.Params.TrueMilliAmps != nil {
		c.chargerState.CurrentSessionCurrentMilliamps = *msg.Params.TrueMilliAmps
	}
	if msg.Params.Voltage != nil {
		c.chargerState.CurrentSessionVoltage = *msg.Params.Voltage
	}
	if msg.Params.Voltage != nil && msg.Params.TrueMilliAmps != nil {
		c.chargerState.CurrentSessionPower = float64(*msg.Params.Voltage) * float64(*msg.Params.TrueMilliAmps)
	}
	if msg.Params.WattHours != nil {
		c.chargerState.SessionWatthours = *msg.Params.WattHours
	}
	if msg.Params.CarbonSavedGrams != nil {
		c.chargerState.SessionCarbonSavedGrams = *msg.Params.CarbonSavedGrams
	}
	if msg.Params.CtCurrent != nil {
		c.chargerState.CurrentSessionCtCurrent = *msg.Params.CtCurrent // TODO: check this is right
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

func (c *Client) handleSyncApply(msg *types.SyncApply) error {
	if msg == nil {
		return errors.New("cannot handle nil message")
	}
	for i := range msg.Params {
		if msg.Params[i].Brightness != nil {
			c.chargerState.LEDBrightness = *msg.Params[i].Brightness
		}
		if msg.Params[i].LockState != nil {
			switch *msg.Params[i].LockState {
			case "locked":
				c.chargerState.LockState = state.LOCKED
			case "unlocked":
				c.chargerState.LockState = state.UNLOCKED
			case "pending":
				c.chargerState.LockState = state.PENDING_LOCK
			}
		}
		if msg.Params[i].ReleaseState != nil {
			switch *msg.Params[i].ReleaseState {
			case "default":
				c.chargerState.ReleaseState = state.DEFAULT
			case "released":
				c.chargerState.ReleaseState = state.RELEASED
			}
		}
		if msg.Params[i].MaxCurrent != nil {
			c.chargerState.MaxCurrentMilliamps = *msg.Params[i].MaxCurrent
		}
		if msg.Params[i].CtFlags != nil { // TODO
		}
		if msg.Params[i].SolarMode != nil {
			switch *msg.Params[i].SolarMode {
			case "boost":
				c.chargerState.ChargeMode = state.BOOST
			case "eco":
				c.chargerState.ChargeMode = state.ECO
			case "super_eco": // TODO: check these
				c.chargerState.ChargeMode = state.SUPER_ECO
			}
		}
		if msg.Params[i].Features != nil { // TODO
		}
		if msg.Params[i].RandomStart != nil {
			c.chargerState.RandomStart = *msg.Params[i].RandomStart
		}
		if msg.Params[i].EffectName != nil { // TODO
		}
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
