package db

import (
	"log"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MustInitMongo(uri, dbName string) {
	if err := mgm.SetDefaultConfig(nil, dbName, options.Client().ApplyURI(uri)); err != nil {
		log.Fatalf("mongo init failed: %v", err)
	}
	log.Println("mongo: connected (mgm)")
}
