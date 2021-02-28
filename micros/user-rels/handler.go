package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	micros "github.com/red-gold/ts-serverless/micros"
	"github.com/red-gold/ts-serverless/micros/user-rels/handlers"
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
		server.POST("/follow", handlers.FollowHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/unfollow/:userId", handlers.UnfollowHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/circle/:circleId", handlers.DeleteCircle(db), coreServer.RouteProtectionCookie)
		server.PUT("/circles", handlers.UpdateRelCirclesHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/followers", handlers.GetFollowersHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/following", handlers.GetFollowingHandle(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
