package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	micros "github.com/red-gold/ts-serverless/micros"
	"github.com/red-gold/ts-serverless/micros/posts/handlers"
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
		server.POST("/", handlers.CreatePostHandle(db), coreServer.RouteProtectionCookie)
		server.POST("/index", handlers.InitPostIndexHandle(db), coreServer.RouteProtectionHMAC)
		server.PUT("/", handlers.UpdatePostHandle(db), coreServer.RouteProtectionCookie)
		server.PUT("/profile", handlers.UpdatePostProfileHandle(db), coreServer.RouteProtectionCookie)
		server.PUT("/score/+1/:postId", handlers.IncrementScoreHandle(db), coreServer.RouteProtectionHMAC)
		server.PUT("/score/-1/:postId", handlers.DecrementScoreHandle(db), coreServer.RouteProtectionHMAC)
		server.PUT("/comment/+1/:postId", handlers.IncrementCommentHandle(db), coreServer.RouteProtectionHMAC)
		server.PUT("/comment/-1/:postId", handlers.DecrementCommentHandle(db), coreServer.RouteProtectionHMAC)
		server.DELETE("/:postId", handlers.DeletePostHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/", handlers.QueryPostHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/:postId", handlers.GetPostHandle(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
