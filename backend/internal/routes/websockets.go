package routes

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/realtime"
	"backend/internal/utils"
	"context"
	"net/http"
	"strings"

	"github.com/coder/websocket"
)

func handleWebsocket(hub *realtime.Hub, w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get(middlewares.CtxAuthorization)
	if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer "){
		utils.JSON(w, http.StatusUnauthorized, false, "Missing or invalid Authorization header", nil)
		return
	}
	accessToken := strings.TrimSpace(authHeader[7:])

	userId, _, _, err := utils.VerifyJWT(accessToken)

	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, false, "Missing or invalid Authorization header", nil)
		return
	}

	user, err := models.GetUserByID(userId)
	if err != nil {
		utils.JSON(w, http.StatusUnauthorized, false, "Missing or invalid Authorization header", nil)
		return
	}

	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{""},
	}

	connection, err := websocket.Accept(w, r, opts)
	if err != nil {
		utils.JSON(w, http.StatusInternalServerError, false, "Failed to accept WebSocket connection", nil)
		return
	}

	client:= realtime.NewClient(user, connection)

	hub.RegisterClientConnection(client)
	hub.SendCurrentClients(client)

	defer func() {
		hub.UnregisterClientConnection(client)
		client.Close()
	}()
	
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go heartbeat(ctx, client)
	go writePump(ctx, client)
	readPump(ctx, cancel, hub , client)





}