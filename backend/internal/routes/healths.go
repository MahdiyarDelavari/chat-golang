package routes

import (
	"backend/internal/utils"
	"net/http"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func handleHealthCheckHTTP(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, true, "API is healthy:)", nil)
}

func handleHealthCheckWs(w http.ResponseWriter, r *http.Request) {
	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	}

	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "Connection closed")

	ctx := r.Context()

	for {
		var message string
		err := wsjson.Read(ctx, conn, &message)
		if err != nil {
			break
		}

		response := map[string]any{
			"data":    message,
			"from":    "server",
			"success": true,
		}

		err = wsjson.Write(ctx, conn, response)
		if err != nil {
			break
		}
	}

}