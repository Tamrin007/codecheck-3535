package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	upgrader = websocket.Upgrader{
		// 別ホストからのリクエストを許可
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

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
func botFilter(msg string) (string, string) {
	mentionFlags := []string{"bot ", "bot　", "@bot ", "@bot　", "bot:"}
	for _, flag := range mentionFlags {
		if strings.HasPrefix(msg, flag) {
			return strings.TrimLeft(msg, flag), "bot"
		}
	}
	return msg, "message"
}

func wsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return err
		}
		var req request
		json.Unmarshal(msg, &req)

		var resp response
		req.Text, resp.Type = botFilter(req.Text)

		resp.Text = req.Text
		resp.Success = true
		if req.Text == "ping" && resp.Type == "bot" {
			resp.Text = "pong"
		}

		respJSON, err := json.Marshal(resp)
		fmt.Println(string(respJSON))
		if err != nil {
			return err
		}

		err = ws.WriteMessage(websocket.TextMessage, respJSON)
		if err != nil {
			return err
		}
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.GET("/", wsHandler)
	e.Logger.Fatal(e.Start(":1323"))
}
