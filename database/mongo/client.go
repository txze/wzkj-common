package mongo

import (
	"context"
	"fmt"
	"sync"

	"github.com/txze/wzkj-common/logger"

	"github.com/qiniu/qmgo"
)

type MongoClient struct {
	Client *qmgo.Client
	DB     *qmgo.Database
}

var defaultMongoClient = &MongoClient{}

func InitMongo(username, password, host, port string, database string) {
	var once sync.Once
	once.Do(func() {

		var auth string
		if username != "" {
			auth = fmt.Sprintf("%s@%s", username, password)
		}
		var uri = fmt.Sprintf("mongodb://%s%s:%s", auth, host, port)

		client, err := qmgo.NewClient(context.Background(), &qmgo.Config{Uri: uri})
		if err != nil {
			panic(err)
		}
		defaultMongoClient.Client = client
		defaultMongoClient.DB = client.Database(database)

		logger.Info("Mongodb connect success",
			logger.String("uri", host+":"+port),
			logger.String("database", database))
	})
}
func DB() *qmgo.Database {
	return defaultMongoClient.DB
}
