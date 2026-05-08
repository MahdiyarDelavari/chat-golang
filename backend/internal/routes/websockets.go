package routes

import (
	"backend/internal/middlewares"
	"backend/internal/models"
	"backend/internal/realtime"
	"backend/internal/utils"
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
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


func heartbeat(ctx context.Context, client *realtime.Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err := client.Conn.Ping(pingCtx)
			if err != nil {
				log.Println("Ping failed, disconnecting client.")
				cancel()
				return
			}
			cancel()

			client.SendEvent(realtime.Event{
			EventType: realtime.EventHeartbeat,
			Payload:   nil,
		})
		}
	}
}


func writePump(ctx context.Context, client *realtime.Client) {
	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-client.SendChannel():
			if !ok {
				return
			}

			writeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			_ = wsjson.Write(writeCtx, client.Conn, event)
			cancel()
		}
	}
}

func readPump(ctx context.Context, cancel context.CancelFunc, hub *realtime.Hub, client *realtime.Client) {
	defer cancel()
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("Recovered from panic in readPump for client %d: %v", client.User.ID, r)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var event realtime.Event
		err := wsjson.Read(ctx, client.Conn, &event)
		if err != nil {
			return
		}

		handleIncomingEvent(hub, client, event)
	}
}