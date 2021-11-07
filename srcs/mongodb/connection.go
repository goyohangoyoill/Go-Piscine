package mongodb

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func preparation() (info [5]string) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./secret")
	viper.AddConfigPath("../secret")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("viper loading failed")
		return
	}
	info[0] = (viper.Get("DB_KIND")).(string)
	info[1] = (viper.Get("DB_HOST")).(string)
	info[2] = (viper.Get("DB_NAME")).(string)
	info[3] = (viper.Get("DB_USER")).(string)
	info[4] = (viper.Get("DB_PASS")).(string)
	return
}

func MongoConn() (client *mongo.Client, ctx context.Context) {
	info := preparation()
	for _, v := range info {
		if v == "" {
			_ = fmt.Errorf("viper invalid value")
			return nil, nil
		}
	}
	// timeout 기반의 Context 생성
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	// Authentication 을 위한 Client Option 구성
	clientOptions := options.Client().ApplyURI(
		"mongodb://" + info[1] + ":27017").SetAuth(
		options.Credential{
			AuthSource: "",
			Username:   info[3],
			Password:   info[4],
		},
	)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("MongoDB Connection Success")
	return client, ctx
}
