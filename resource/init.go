package resource

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"gopkg.in/olivere/elastic.v3"
	"gopkg.in/redis.v4"
	"io/ioutil"
	"master/constant"
	"time"
)

type ResourceConfig struct {
	IsEnablePostgres bool
	IsEnableRabbit   bool
	IsEnableRedis    bool
	IsEnableElastic  bool
	PostgreSQLLogger bool
}

type Resource struct {
	Config     ResourceConfig
	PostgreSql *gorm.DB
	Rabbit     *amqp.Connection
	Redis      *redis.Client
	Elastic    *elastic.Client
}

func (r *Resource) Close() {
	fmt.Println("[INFO] Close all connect resouce...")
	if r.Config.IsEnablePostgres {
		r.PostgreSql.Close()
	}
	if r.Config.IsEnableRedis {
		r.Redis.Close()
	}
	if r.Config.IsEnableRabbit {
		r.Rabbit.Close()
	}
}

func initPostgreSQL() (*gorm.DB, error) {
	dbSQL, err := gorm.Open("postgres", "postgres://"+constant.PostgresUser+
		":"+constant.PostgresPassword+"@"+constant.PostgresHost+
		":"+constant.PostgresPort+"/"+constant.PostgresDB+
		"?sslmode=disable")
	if err != nil {
		return dbSQL, err
	}
	dbSQL.DB()
	dbSQL.DB().SetMaxIdleConns(10)
	dbSQL.DB().SetMaxOpenConns(150)
	return dbSQL, err
}

func initRedis() (*redis.Client, error) {
	addr := constant.RedisHost + ":" + constant.RedisPort
	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		PoolSize:    500,
		DB:          0,
		Password:    constant.RedisPass,
		PoolTimeout: 5 * time.Second,
		IdleTimeout: 30 * time.Second,
	})
	_, err := client.Ping().Result()
	return client, err
}

func initRabbit() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://" +
		constant.RabbitUser +
		":" + constant.RabbitPass +
		"@" + constant.RabbitHost +
		":" + constant.RabbitPort + "/")
	return conn, err
}

func initElastic() *elastic.Client {
	client, err := elastic.NewClient(elastic.SetURL(constant.ElasticHost + ":" + constant.ElasticPort))
	if err != nil {
		panic(err)
	}
	if client == nil {
		panic(errors.New("no elastic search"))
	}
	err = setupIndex(client)
	if err != nil {
		panic(err)
	}
	return client
}

func deleteIndex(client *elastic.Client) error {
	_, err := client.DeleteIndex(constant.ElasticNamspace).Do()
	return err
}

func createIndex(client *elastic.Client) error {
	settings, err := readTextFromFile("../pixai-master/resource/elastic.json")
	if err != nil {
		return err
	}
	_, err = client.CreateIndex(constant.ElasticNamspace).Body(settings).Do()
	return err
}

func setupIndex(client *elastic.Client) error {
	exists, err := client.IndexExists(constant.ElasticNamspace).Do()
	if err != nil {
		return err
	}

	if !exists {
		return createIndex(client)
	}
	return nil
}

func readTextFromFile(fileName string) (string, error) {
	content, err := ioutil.ReadFile(fileName)
	return string(content), err
}

func Init(config ResourceConfig) (*Resource, error) {
	r := Resource{}
	r.Config = config
	if config.IsEnablePostgres {
		pq, err := initPostgreSQL()
		if err != nil {
			fmt.Println("[ERROR] Connect PostgreSQL Failed...", err)
			return nil, err
		}
		if config.PostgreSQLLogger {
			pq.LogMode(config.PostgreSQLLogger)
		}
		r.PostgreSql = pq
	}
	if config.IsEnableRedis {
		redisClient, err := initRedis()
		if err != nil {
			fmt.Println("[ERROR] Connect Reids Failed...", err)
			return nil, err
		}
		r.Redis = redisClient
	}
	if config.IsEnableRabbit {
		rabbit, err := initRabbit()
		if err != nil {
			fmt.Println("[ERROR] Connect rabbit Failed...", err)
			return nil, err
		}
		r.Rabbit = rabbit
	}
	if config.IsEnableElastic {
		elastic := initElastic()
		r.Elastic = elastic
	}
	return &r, nil
}
