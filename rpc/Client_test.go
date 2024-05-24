package rpc

import (
	"context"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/mocks"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/types"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateChargerState(t *testing.T) {
	t.Run("with valid get session message", func(t *testing.T) {
		mockClient := Client{}

		session := &types.SessionParams{Charging: false, TrueMilliAmps: 0, Voltage: 0, WattHours: 0, CarbonSavedGrams: 0, CtCurrent: 1900, EvPower: 0, GridPower: 287, HousePower: 287, GenerationPower: 0}

		message := RawMessage{
			"jsonrpc": "2.0",
			"method":  "get.session",
			"params": map[string]interface{}{
				"charging":           false,
				"true_milli_amps":    0,
				"voltage":            0,
				"watt_hours":         0,
				"carbon_saved_grams": 0,
				"ct_current":         1900,
				"ev_power":           0,
				"grid_power":         287,
				"house_power":        287,
				"generation_power":   0,
			},
		}

		err := mockClient.updateChargerState(message)
		assert.NoError(t, err)

		// Check if GetCurrentSession return same session data
		assert.Equal(t, session, mockClient.GetCurrentSession())
	})

	t.Run("with unsupported message method", func(t *testing.T) {
		mockClient := Client{}
		message := RawMessage{"method": "unsupported"}

		err := mockClient.updateChargerState(message)
		assert.NoError(t, err) // since we do not handle unknown methods we should not get error
	})

	t.Run("with malformed json", func(t *testing.T) {
		mockClient := Client{}
		message := RawMessage{"method": "get.session", "params": "{"}

		err := mockClient.updateChargerState(message)
		assert.Error(t, err)
	})

	t.Run("with missing method field", func(t *testing.T) {
		mockClient := Client{}
		message := RawMessage{}

		err := mockClient.updateChargerState(message)
		assert.Error(t, err)
		assert.Equal(t, "no method in message", err.Error())
	})
}

func TestWaitForResponse(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T, c *Client, id string)
	}{
		{
			name: "CancelledContext",
			test: func(t *testing.T, c *Client, id string) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // immediate cancellation
				_, err := c.waitForResponse(ctx, id)
				assert.Error(t, err)
				assert.Equal(t, "operation timed out", err.Error())
			},
		},
		{
			name: "ValuePresentInResponseMap",
			test: func(t *testing.T, c *Client, id string) {
				expectedResponse := []byte("Test response")
				c.responseMap[id] = expectedResponse
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()
				resp, err := c.waitForResponse(ctx, id)
				assert.NoError(t, err)
				assert.Equal(t, expectedResponse, resp)
			},
		},
		{
			name: "ValueNotPresentInResponseMap",
			test: func(t *testing.T, c *Client, id string) {
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()
				_, err := c.waitForResponse(ctx, id)
				assert.Error(t, err)
				assert.Equal(t, "operation timed out", err.Error())
			},
		},
	}

	// mocking required fields in Client struct for testing
	wWrapperIfaceMock := mocks.NewWebsocketWrapperIface(t)
	c := &Client{
		syncConnection: wWrapperIfaceMock,
		responseMap:    make(map[string][]byte),
		syncReceiver: responseReceiver{
			mutex: sync.Mutex{},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable.
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			id := "test_id"
			tt.test(t, c, id)
		})
	}
}
