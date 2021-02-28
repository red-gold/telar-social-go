package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	micros "github.com/red-gold/ts-serverless/micros"
	"github.com/red-gold/ts-serverless/micros/circles/handlers"
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
		server.POST("/", handlers.CreateCircleHandle(db), coreServer.RouteProtectionCookie)
		server.POST("/following/:userId", handlers.CreateFollowingHandle(db), coreServer.RouteProtectionHMAC)
		server.PUT("/", handlers.UpdateCircleHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/:circleId", handlers.DeleteCircleHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/my", handlers.GetMyCircleHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/id/:circleId", handlers.GetCircleHandle(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
