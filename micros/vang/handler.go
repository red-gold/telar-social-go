package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	micros "github.com/red-gold/ts-serverless/micros"
	"github.com/red-gold/ts-serverless/micros/vang/handlers"
)

func init() {

	micros.InitConfig()
}

// Cache state
var server *coreServer.ServerRouter
var db interface{}

// Handler function
func Handle(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	// Start
	if db == nil {
		var startErr error
		db, startErr = micros.Start(ctx)
		if startErr != nil {
			fmt.Printf("Error startup: %s", startErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(startErr.Error()))
		}
	}

	// Server Routing
	if server == nil {
		server = coreServer.NewServerRouter()
		server.POST("/messages", handlers.SaveMessages(db), coreServer.RouteProtectionCookie)
		server.PUT("/message", handlers.UpdateMessageHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/message/:messageId", handlers.DeleteMessageHandle(db), coreServer.RouteProtectionCookie)
		server.POST("/room/active", handlers.ActivePeerRoom(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
