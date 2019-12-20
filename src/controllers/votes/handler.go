package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	"github.com/red-gold/ts-serverless/src/controllers"
	cf "github.com/red-gold/ts-serverless/src/controllers/votes/config"
	"github.com/red-gold/ts-serverless/src/controllers/votes/handlers"
)

func init() {

	cf.InitConfig()
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
		db, startErr = controllers.Start(ctx)
		if startErr != nil {
			fmt.Printf("Error startup: %s", startErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(startErr.Error()))
		}
	}

	// Server Routing
	if server == nil {
		server = coreServer.NewServerRouter()
		server.POST("/", handlers.CreateVoteHandle(db), coreServer.RouteProtectionCookie)
		server.PUT("/", handlers.UpdateVoteHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/id/:voteId", handlers.DeleteVoteHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/post/:postId", handlers.DeleteVoteByPostIdHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/", handlers.GetVotesByPostIdHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/:voteId", handlers.GetVoteHandle(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
