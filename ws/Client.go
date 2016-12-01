package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 別ホストからのリクエストを許可
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	bot  chan []byte
}

type request struct {
	Text string `json:"text"`
}

type response struct {
	Success bool   `json:"success"`
	Type    string `json:"type"`
	Text    string `json:"text"`
}

// ボットへのメンションだった場合：メンション部を除いたメッセージと Type: "bot" を返す
// そうでない場合：元のメッセージと Type: "message" を返す
func botFilter(msg string) (string, bool) {
	mentionFlags := []string{"bot ", "bot　", "@bot ", "@bot　", "bot:"}
	for _, flag := range mentionFlags {
		if strings.HasPrefix(msg, flag) {
			return strings.TrimLeft(msg, flag), true
		}
	}
	return msg, false
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		var req request
		json.Unmarshal(message, &req)

		var resp response
		resp.Success = true
		resp.Type = "message"
		resp.Text = req.Text

		respJSON, err := json.Marshal(resp)
		if err != nil {
			log.Printf("error: %v", err)
		}
		c.hub.broadcast <- respJSON

		if msg, isBot := botFilter(req.Text); isBot {
			var botMsg response
			botMsg.Success = true
			botMsg.Type = "bot"
			if msg == "ping" {
				botMsg.Text = "pong"
			}
			botJSON, err := json.Marshal(botMsg)
			if err != nil {
				log.Printf("error: %v", err)
			}
			c.hub.botcast <- botJSON
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case message, ok := <-c.bot:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.bot)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.bot)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func Handler(hub *Hub, c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := &Client{
		hub:  hub,
		conn: ws,
		send: make(chan []byte, 256),
		bot:  make(chan []byte, 256),
	}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
	return nil
}
