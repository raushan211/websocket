package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	r := gin.Default()
	setupRoutes(r)
	r.Run()
}

type Quote struct {
	Q string `json:"q"`
	A string `json:"a"`
	H string `json:"h"`
}

func setupRoutes(r *gin.Engine) {
	r.GET("/connect", registerClient)

}

var wsupgraders = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func registerClient(c *gin.Context) {

	wsupgraders.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsupgraders.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		msg := fmt.Sprintf("Failed to set websocket upgrade: %+v", err)
		fmt.Println(msg)
		return
	}
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 5)
		mType, mByte, err := conn.ReadMessage()
		fmt.Println("mByte: ", string(mByte))
		fmt.Println("mType: ", mType)
		fmt.Println("err: ", err)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", getQuote())))
	}
	conn.Close()
}

func getQuote() string {

	url := "https://zenquotes.io/api/random"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "unable to generate code"
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "unable to generate code"
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "unable to generate code"
	}
	fmt.Println(string(body))

	data := []Quote{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return "unable to generate code"
	}
	return data[0].Q

}
