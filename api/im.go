package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golangIM/dao"
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
func groupWebSocketHandler(c *gin.Context) {
	groupID := c.Param("groupID")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		utils.RespFail(c, "username not found")

		return
	}
	username, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username is not a string"})
		return
	}
	groupMembers := dao.Getgroupmember(groupID) // 获取群组中的成员列表
	connectedClients[username] = conn

	defer func() {
		conn.Close()
		delete(connectedClients, username)
	}()

	fmt.Println(username, "connected")

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		if messageType == websocket.TextMessage {
			for _, conn := range connectedClients {
				// 检查连接是否存在于 groupMembers 中
				for _, member := range groupMembers {
					if conn == connectedClients[member.MemberName] {
						// 发送消息给匹配的连接
						err := conn.WriteMessage(websocket.TextMessage, msg)
						if err != nil {
							fmt.Println("Write error:", err)
						}
						break
					}
				}
			}
		}
	}
}
func groupuploadImage(c *gin.Context) {
	groupID := c.Param("groupID")
	groupMembers := dao.Getgroupmember(groupID) // 获取群组中的成员列表

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

	for _, conn := range connectedClients {
		// 检查连接是否存在于 groupMembers 中
		for _, member := range groupMembers {
			if conn == connectedClients[member.MemberName] {
				// 发送消息给匹配的连接
				err := conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					fmt.Println("Write error:", err)
				}
				break
			}
		}
	}

	utils.RespSuccess(c, "File uploaded and sent successfully")
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
func invitefriend(c *gin.Context) {
	newmember := c.PostForm("id")
	groupid := c.PostForm("groupid")
	dao.Addmember(newmember, groupid)
	utils.RespSuccess(c, "邀请成功")
}
func creategroup(c *gin.Context) {
	groupname := c.Param("groupname")
	dao.Creategroup(groupname)
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		utils.RespFail(c, "username not found")

		return
	}
	username, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username is not a string"})
		return
	}
	userid := dao.Selectid(username)
	groupid := dao.Selectgroupid(groupname)
	dao.Addmember(userid, groupid)
	utils.RespSuccess(c, "建群成功")
}
