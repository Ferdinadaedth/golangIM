package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		fmt.Printf("Received from %s: %s\n", userID, msg)

		// Parse the message to find the target user
		// The message format could be "targetUserID:message"
		parts := strings.SplitN(string(msg), ":", 2)
		if len(parts) == 2 {
			targetUserID := parts[0]
			messageToSend := parts[1]

			// Check if the target user is connected
			targetConn, exists := connectedClients[targetUserID]
			if exists {
				// Send the message to the target user
				if err := targetConn.WriteMessage(websocket.TextMessage, []byte(messageToSend)); err != nil {
					fmt.Println("Write error:", err)
				}
			} else {
				fmt.Printf("Target user %s is not connected.\n", targetUserID)
			}
		}
	}
}
func websocketHandler1(c *gin.Context) {
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
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		fmt.Printf("Received from %s: %s\n", userID, msg)

		// Parse the message to find the message type and content
		parts := strings.SplitN(string(msg), ":", 2)
		if len(parts) == 2 {
			messageType := parts[0]
			messageContent := parts[1]

			if messageType == "image" {
				// Broadcast the image data to other clients
				for otherID, otherConn := range connectedClients {
					if otherID != userID {
						if err := otherConn.WriteMessage(websocket.TextMessage, []byte("image:"+messageContent)); err != nil {
							fmt.Println("Write error:", err)
							break
						}
					}
				}
			}
		}
	}
}
