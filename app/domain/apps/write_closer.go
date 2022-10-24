package apps

import "github.com/gorilla/websocket"

type WriteCloser struct {
	*websocket.Conn
}

func (wc *WriteCloser) Read(p []byte) (n int, err error) {
	messageType, r, err := wc.NextReader()
	if err != nil {
		return 0, err
	}
	if messageType != websocket.TextMessage {
		return 0, nil
	}
	return r.Read(p)
}

func (wc *WriteCloser) Write(p []byte) (n int, err error) {
	return len(p), wc.Conn.WriteMessage(websocket.TextMessage, p)
}
