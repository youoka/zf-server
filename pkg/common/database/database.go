package database

import (
	"context"
	"fmt"
	go_redis "github.com/go-redis/redis/v8"
	"time"
	"zf-server/pkg/common/config"
)

var DB DataBases

type DataBases struct {
	MysqlDB mysqlDB
	RDB     go_redis.UniversalClient
}

func init() {
	fmt.Println("init mysql redis ")

	initMysqlDB()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if config.Config.Redis.EnableCluster {
		DB.RDB = go_redis.NewClusterClient(&go_redis.ClusterOptions{
			Addrs:    config.Config.Redis.DBAddress,
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			PoolSize: 50,
		})
		_, err := DB.RDB.Ping(ctx).Result()
		if err != nil {
			fmt.Println("redis cluster failed address ", config.Config.Redis.DBAddress)
			panic(err.Error() + " redis cluster " + config.Config.Redis.DBUserName + config.Config.Redis.DBPassWord)
		}
	} else {
		DB.RDB = go_redis.NewClient(&go_redis.Options{
			Addr:     config.Config.Redis.DBAddress[0],
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			DB:       0,                              // use default DB
			PoolSize: 100,                            // 连接池大小
		})
		_, err := DB.RDB.Ping(ctx).Result()
		if err != nil {
			panic(err.Error() + " redis " + config.Config.Redis.DBAddress[0] + config.Config.Redis.DBUserName + config.Config.Redis.DBPassWord)
		}
	}
	fmt.Println("init mysql redis ok ")
}
