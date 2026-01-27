package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Migan178/misschord-backend/internal/chat"
	"github.com/Migan178/misschord-backend/internal/repository"
	"github.com/Migan178/misschord-backend/internal/repository/ent"
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *chat.Hub, authMiddleware *jwt.GinJWTMiddleware, w http.ResponseWriter, r *http.Request) {
	var user *ent.User

	// goto를 쓰는게 중첩된 if보다 나을 듯
	// 쿠키의 토큰으로 인증하고 (브라우저에서), 인증이 실패하면 websocket 데이터로 인증
	{
		tokenCookie, err := r.Cookie("token_jwt")
		if err != nil {
			goto Upgrade
		}

		token, err := authMiddleware.ParseTokenString(tokenCookie.String())
		if err != nil {
			goto Upgrade
		}

		claims := jwt.ExtractClaimsFromToken(token)
		id, err := strconv.Atoi(claims["id"].(string))
		if err != nil {
			goto Upgrade
		}

		if dbUser, err := repository.GetDatabase().Users.Get(r.Context(), id); err == nil {
			user = dbUser
		}
	}

Upgrade:
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}

	client := chat.NewClient(hub, user, ws)

	client.Start()
}
