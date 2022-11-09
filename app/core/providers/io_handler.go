package providers

import (
	"bytes"
	"io"
	"sync"
	"unicode"

	"github.com/gorilla/websocket"
)

type ShellLogger struct {
	sync.Mutex
	Base        io.ReadWriteCloser
	Term        io.Writer
	addComplete bool
}

func (sl *ShellLogger) Read(p []byte) (n int, err error) {
	if n, err = sl.Base.Read(p); err != nil || n == 0 {
		return
	} else {
		sl.Term.Write(p[:n])
		sl.Lock()
		defer sl.Unlock()
		sl.addComplete = p[n-1] == '\t'
		return
	}
}

func (sl *ShellLogger) Write(p []byte) (n int, err error) {
	n, err = sl.Base.Write(p)
	sl.Lock()
	defer sl.Unlock()
	if sl.addComplete {
		for _, c := range string(p) {
			if unicode.IsPrint(c) {
				sl.Term.Write([]byte(string(c)))
			}
		}
		if len(p) == 0 || p[len(p)-1] != '\a' {
			sl.addComplete = false
		}
	}
	return
}
func (sl *ShellLogger) Close() error {
	return sl.Base.Close()
}

type OptionalWriter struct {
	bytes.Buffer
	disable bool
}

func (ow *OptionalWriter) Write(p []byte) (int, error) {
	if ow.disable {
		return len(p), nil
	}
	return ow.Buffer.Write(p)
}

func (ow *OptionalWriter) Close() error {
	return nil
}

func (ow *OptionalWriter) Disable() *OptionalWriter {
	ow.disable = true
	return ow
}

type WebSocketWriter struct {
	*websocket.Conn
}

func (wsw *WebSocketWriter) Read(p []byte) (n int, err error) {
	if messageType, r, err := wsw.NextReader(); err != nil {
		return 0, err
	} else if messageType != websocket.TextMessage {
		return 0, nil
	} else {
		return r.Read(p)
	}
}

func (wsw *WebSocketWriter) Write(p []byte) (n int, err error) {
	return len(p), wsw.Conn.WriteMessage(websocket.TextMessage, p)
}
