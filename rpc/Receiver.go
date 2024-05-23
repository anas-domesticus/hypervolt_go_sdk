package rpc

import (
	"encoding/json"
	"fmt"
	"sync"
)

type responseReceiver struct {
	messageBuffer      chan RawMessage
	responseChan       chan RawMessage
	updateChan         chan RawMessage
	connection         WebsocketWrapperIface
	mutex              sync.Mutex
	responseMap        map[string][]byte
	updateChargerState func(message RawMessage) error
}

func (r *responseReceiver) receiverLoop() error {
	for {
		_, message, err := r.connection.ReadMessage()
		if err != nil {
			return err
		}
		data := RawMessage{}
		err = json.Unmarshal(message, &data)
		if err != nil {
			return err
		}
		r.messageBuffer <- data
	}
}

func (r *responseReceiver) dispatcherLoop() error {
	for {
		message := <-r.messageBuffer
		_, ok := message["result"]
		if ok {
			r.responseChan <- message
			continue
		}
		r.updateChan <- message
	}
}

func (r *responseReceiver) responseLoop() error {
	for {
		message := <-r.responseChan
		_, ok := message["id"]
		if !ok {
			fmt.Println("message lacking ID received in response loop")
			continue
		}
		val, ok := message["id"].(string)
		if !ok {
			fmt.Println("ID received incorrect type")
		}
		bytes, err := json.Marshal(message)
		if err != nil {
			fmt.Println("failed to marshal message")
		}
		r.mutex.Lock()
		r.responseMap[val] = bytes
		r.mutex.Unlock()
	}
}

func (r *responseReceiver) updateLoop() error {
	for {
		message := <-r.updateChan
		err := r.updateChargerState(message)
		if err != nil {
			return err
		}
	}
}

type RawMessage map[string]interface{}