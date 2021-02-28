package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	micros "github.com/red-gold/ts-serverless/micros"
	"github.com/red-gold/ts-serverless/micros/gallery/handlers"
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
		server.POST("/", handlers.CreateMediaHandle(db), coreServer.RouteProtectionCookie)
		server.POST("/list", handlers.CreateMediaListHandle(db), coreServer.RouteProtectionCookie)
		server.PUT("/", handlers.UpdateMediaHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/id/:mediaId", handlers.DeleteMediaHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/dir/:dir", handlers.DeleteDirectoryHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/", handlers.QueryAlbumHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/id/:mediaId", handlers.GetMediaHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/dir/:dir", handlers.GetMediaByDirectoryHandle(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
