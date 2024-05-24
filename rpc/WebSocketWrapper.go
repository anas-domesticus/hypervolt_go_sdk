package rpc

import "github.com/gorilla/websocket"

//go:generate mockery --name WebsocketWrapperIface

type WebsocketWrapperIface interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
	Close() error
}

type WebsocketWrapper struct {
	w *websocket.Conn
}

func (wr *WebsocketWrapper) ReadMessage() (messageType int, p []byte, err error) {
	return wr.w.ReadMessage()
}

func (wr *WebsocketWrapper) WriteMessage(messageType int, data []byte) error {
	return wr.w.WriteMessage(messageType, data)
}

func (wr *WebsocketWrapper) Close() error {
	return wr.w.Close()
}
