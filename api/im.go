package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golangIM/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connectedClients = make(map[string]*websocket.Conn)

func websocketHandler(c *gin.Context) {
	userID := c.Param("userID")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}

	connectedClients[userID] = conn
	defer func() {
		conn.Close()
		delete(connectedClients, userID)
	}()

	fmt.Println(userID, "connected")

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		if messageType == websocket.TextMessage {
			fmt.Printf("Received from %s: %s\n", userID, msg)
			parts := strings.SplitN(string(msg), ":", 2)
			if len(parts) == 2 {
				targetUserID := parts[0]
				messageToSend := parts[1]
				targetConn, exists := connectedClients[targetUserID]
				if exists {
					if err := targetConn.WriteMessage(websocket.TextMessage, []byte(messageToSend)); err != nil {
						fmt.Println("Write error:", err)
					}
				} else {
					fmt.Printf("Target user %s is not connected.\n", targetUserID)
				}
			}

		} else if messageType == websocket.BinaryMessage {
			// 处理图片数据消息

			// 将收到的图片数据广播给所有连接的客户端
			for _, clientConn := range connectedClients {
				err := clientConn.WriteMessage(websocket.BinaryMessage, msg)
				if err != nil {
					fmt.Println("Write error:", err)
				}
			}
		}
	}
}

func uploadImage(c *gin.Context) {
	userID := c.Param("userID")
	conn, exists := connectedClients[userID]
	if !exists {
		utils.RespFail(c, "User not connected")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.RespFail(c, "Error uploading file")
		return
	}

	// 读取上传的文件数据
	fileData, err := file.Open()
	if err != nil {
		utils.RespFail(c, "Error reading file")
		return
	}
	defer fileData.Close()

	// 将文件数据发送给WebSocket连接的客户端
	data, err := ioutil.ReadAll(fileData)
	if err != nil {
		utils.RespFail(c, "Error reading file data")
		return
	}

	err = conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		utils.RespFail(c, "Error sending file data over WebSocket")
		return
	}

	utils.RespSuccess(c, "File uploaded and sent successfully")
}
