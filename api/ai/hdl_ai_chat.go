package ai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/coze-dev/coze-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatRequestParams struct {
	Message        string `json:"message"`
	ReportContent  string `json:"report_content"`
	ConversationID string `json:"conversation_id"`
	CustomerID     string `json:"customer_id"`
}

type ServiceChatResponse struct {
	Code           int64  `json:"code"`
	Type           string `json:"type"`
	Message        string `json:"message"`
	ConversationId string `json:"conversation_id"`
}

func (a *App) hdlAIChat(c *gin.Context) {
	// 检查是否携带 Authorization 头
	jwtHeader := c.Request.Header.Get("Authorization")
	if jwtHeader == "" || !strings.HasPrefix(jwtHeader, "Bearer ") {
		c.JSON(http.StatusBadRequest, &ServiceChatResponse{
			Code:    http.StatusBadRequest,
			Message: "authorization missing",
		})
		return
	}

	// 检查 JWT 有效性
	jwtContent := strings.TrimPrefix(jwtHeader, "Bearer ")
	if !a.CheckJWT(jwtContent) {
		c.JSON(http.StatusBadRequest, &ServiceChatResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid JWT",
		})
		return
	}

	// 将 HTTP 请求升级为 WebSocket 请求
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		c.JSON(http.StatusInternalServerError, &ServiceChatResponse{
			Code:    http.StatusInternalServerError,
			Message: "websocket upgrade failed",
		})
		return
	}
	defer conn.Close()

	for {
		// 读取请求参数
		var param ChatRequestParams
		err := conn.ReadJSON(&param)
		if err != nil {
			log.Println("Read error:", err)
			conn.WriteJSON(&ServiceChatResponse{
				Code: 400,
			})
			break
		}

		if param.Message == "" {
			conn.WriteJSON(&ServiceChatResponse{
				Code: 201,
				Type: "heart",
			})
			continue
		}

		if param.CustomerID == "" {
			param.CustomerID = "guest"
		}

		// 构建 AI 请求
		req := &coze.CreateChatsReq{
			BotID:  a.botID,
			UserID: param.CustomerID,
			Messages: []*coze.Message{
				coze.BuildUserQuestionText(param.Message, nil),
			},
			ConversationID: param.ConversationID,
			CustomVariables: map[string]string{
				"reportContent": param.ReportContent,
			},
		}

		// 流式读取 AI 响应
		resp, err := a.api.Chat.Stream(context.Background(), req)
		if err != nil {
			log.Printf("Error starting chats: %v\n", err)
			conn.WriteJSON(&ServiceChatResponse{
				Code: 400,
			})
			break
		}
		defer resp.Close()

		for {
			event, err := resp.Recv()

			// 流关闭，请求完成
			if errors.Is(err, io.EOF) {
				log.Println("Stream finished")
				break
			}
			if err != nil {
				log.Println(err)
				break
			}

			// 读取响应内容，并发送给前端
			if event.Event == coze.ChatEventConversationMessageDelta {
				if strings.TrimSpace(event.Message.Content) == "" {
					continue
				}
				rsp := &ServiceChatResponse{
					Code:           200,
					Type:           "text",
					Message:        event.Message.Content,
					ConversationId: event.Message.ConversationID,
				}
				conn.WriteJSON(rsp)
			} else if event.Event == coze.ChatEventConversationChatCompleted {
				rsp := &ServiceChatResponse{
					Code: 300,
				}
				conn.WriteJSON(rsp)
				break
			} else {
				continue
			}
		}
	}
}

func (a *App) CheckJWT(tokenString string) bool {
	// 解码JWT
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法是否为HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secretKey), nil
	})

	if err != nil {
		log.Printf("JWT解析失败: %v", err)
		return false
	}

	// 验证令牌有效性
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		now := time.Now().Unix()
		iat := claims.IssuedAt
		// 允许 10 分钟有效期
		if abs(now-iat) <= 600 {
			return true
		}
		return false
	} else {
		return false
	}
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
