package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait       = 10 * time.Second
	pingInterval   = (pongWait * 9) / 10
	maxMessageSize = 512
)

var connectionUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func isWebSocketRequest(r *http.Request) bool {
	connectionHeaderValue := r.Header.Get("Connection")
	return connectionHeaderValue == "Upgrade"
}

func connectWebSocket(w http.ResponseWriter, r *http.Request, s IService, sc *serverConnector) error {
	slog.Info("new ws connection from", "address", r.RemoteAddr, "url", r.URL)
	conn, err := connectionUpgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade connection error", "err", err)
		return err
	}
	sc.wsconn = conn
	go sc.readMessages(s)
	go sc.writeMessages(s)
	return nil
}

func (sc *serverConnector) writeMessages(s IService) {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		sc.wsconn.Close()
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-s.GetOutgoingMessagesQueue():
			if !ok {
				slog.Warn("write close message problem", "conn", sc.wsconn.RemoteAddr())
				if err := sc.wsconn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					slog.Error("close message error", "conn", sc.wsconn.RemoteAddr(), "err", err)
					return
				}
			}
			slog.Debug("sending message", "msg", message)
			if err := sc.wsconn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				slog.Error("error sending message", "msg", message, "err", err)
			}
		case <-ticker.C:
			slog.Debug("sending PING")
			if err := sc.wsconn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				slog.Error("sending PING error", "err", err)
				return
			}
		}
	}
}

func (sc *serverConnector) readMessages(s IService) {
	defer func() {
		sc.wsconn.Close()
	}()

	if err := configureWebSocketConnection(sc.wsconn); err != nil {
		slog.Error("configuring ws connection error", "err", err)
		return
	}

	for {
		messageType, payload, err := sc.wsconn.ReadMessage()
		if err != nil {
			slog.Error("reading message error", "err", err)
			break
		}
		parsedPayload, err := parsePayload(messageType, payload)
		if err != nil {
			slog.Error("payload error", "err", err)
		} else {
			slog.Info("web socket message", "msg", parsedPayload)
			if err := s.OnWebSocketMessage(string(parsedPayload)); err != nil {
				slog.Warn("processing web socket message error", "msg", parsedPayload, "err", err)
			}
		}

	}
}

func parsePayload(messageType int, payload []byte) (string, error) {
	switch messageType {
	case websocket.TextMessage:
		return string(payload), nil
	default:
		return "", errors.Join(ErrorUnsupportedMessageType, fmt.Errorf("message type: %v", messageType))
	}
}

func configureWebSocketConnection(conn *websocket.Conn) error {
	if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		slog.Error("setting ws deadline", "err", err)
		return err
	}
	conn.SetPongHandler(pongHandler(conn))
	conn.SetReadLimit(maxMessageSize)
	return nil
}

func pongHandler(conn *websocket.Conn) func(string) error {
	return func(s string) error {
		slog.Debug("PONG")
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	}
}
