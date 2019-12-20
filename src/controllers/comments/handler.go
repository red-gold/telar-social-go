package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	"github.com/red-gold/ts-serverless/src/controllers"
	cf "github.com/red-gold/ts-serverless/src/controllers/comments/config"
	"github.com/red-gold/ts-serverless/src/controllers/comments/handlers"
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
		server.POST("/", handlers.CreateCommentHandle(db), coreServer.RouteProtectionCookie)
		server.PUT("/", handlers.UpdateCommentHandle(db), coreServer.RouteProtectionCookie)
		server.PUT("/profile", handlers.UpdateCommentProfileHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/id/:commentId/post/:postId", handlers.DeleteCommentHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/post/:postId", handlers.DeleteCommentByPostIdHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/", handlers.GetCommentsByPostIdHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/:commentId", handlers.GetCommentHandle(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
