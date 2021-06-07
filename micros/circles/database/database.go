// Copyright (c) 2021 Amirhossein Movahedi (@qolzam)
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package database

import (
	"context"
	"fmt"

	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/data/mongodb"
)

var Db interface{}

// Connect open database connection
func Connect(ctx context.Context) error {
	coreConfig := config.AppConfig

	switch *coreConfig.DBType {
	case config.DB_MONGO:
		mongoClient, err := mongodb.NewMongoClient(ctx, *coreConfig.MongoDBHost, *coreConfig.Database)
		Db = mongoClient
		return err
	}

	return fmt.Errorf("Please set valid database type in confing file!")
}
