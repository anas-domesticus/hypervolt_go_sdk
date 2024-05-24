package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/anas-domesticus/hypervolt_go_sdk/rpc/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReceiverMessageRouting(t *testing.T) {
	testCases := []struct {
		name                 string
		message              string
		readMessageErr       error
		expectResponse       bool
		expectResponseKey    string
		expectUpdateFuncCall bool
	}{
		{
			name:                 "Auth response",
			message:              `{"id":"1716312852461187","jsonrpc":"2.0","result":{"authenticated":true}}`,
			readMessageErr:       nil,
			expectResponse:       true,
			expectResponseKey:    "1716312852461187",
			expectUpdateFuncCall: false,
		},
		{
			name:                 "get.session message",
			message:              `{"jsonrpc":"2.0","method":"get.session","params":{"charging":false,"true_milli_amps":0,"voltage":0,"watt_hours":0,"carbon_saved_grams":0,"ct_current":1900,"ev_power":0,"grid_power":306,"house_power":306,"generation_power":0}}`,
			readMessageErr:       nil,
			expectResponse:       false,
			expectResponseKey:    "",
			expectUpdateFuncCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			updateFuncWasCalled := false
			updateFuncMessage := make(RawMessage)
			wantFuncMessage := make(RawMessage)

			if tc.expectUpdateFuncCall {
				err := json.Unmarshal([]byte(tc.message), &wantFuncMessage)
				if err != nil {
					fmt.Println(err)
				}
			}
			wsMock := mocks.NewWebsocketWrapperIface(t)
			testReceiver := responseReceiver{
				messageBuffer: make(chan RawMessage, 25),
				responseChan:  make(chan RawMessage),
				updateChan:    make(chan RawMessage),
				connection:    wsMock,
				responseMap:   make(map[string][]byte),
				updateChargerState: func(message RawMessage) error {
					updateFuncWasCalled = true
					updateFuncMessage = message
					return nil
				},
			}

			wsMock.On("ReadMessage").Return(0, []byte(tc.message), tc.readMessageErr)
			_ = testReceiver.StartLoops()

			time.Sleep(100 * time.Millisecond)

			val, ok := testReceiver.responseMap[tc.expectResponseKey]
			assert.Equal(t, ok, tc.expectResponse)
			if tc.expectResponse {
				assert.Equal(t, string(val), tc.message)
			}

			assert.Equal(t, updateFuncWasCalled, tc.expectUpdateFuncCall)
			assert.Equal(t, updateFuncMessage, wantFuncMessage)

		})
	}
}
