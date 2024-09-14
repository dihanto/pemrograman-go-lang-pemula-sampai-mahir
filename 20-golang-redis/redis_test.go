package golangredis

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func TestConnection(t *testing.T) {
	assert.NotNil(t, client)

	err := client.Close()
	assert.Nil(t, err)
}

var ctx = context.Background()

func TestPing(t *testing.T) {
	result, err := client.Ping(ctx).Result()
	assert.Nil(t, err)
	assert.Equal(t, "PONG", result)
}

func TestString(t *testing.T) {
	client.SetEx(ctx, "name", "dihanto", 3*time.Second)

	result, err := client.Get(ctx, "name").Result()
	assert.Nil(t, err)
	assert.Equal(t, "dihanto", result)

	time.Sleep(4 * time.Second)

	result, err = client.Get(ctx, "name").Result()
	assert.NotNil(t, err)
}

func TestList(t *testing.T) {
	client.RPush(ctx, "list", "dihanto")
	client.RPush(ctx, "list", "budi")
	client.RPush(ctx, "list", "joko")

	assert.Equal(t, "dihanto", client.LPop(ctx, "list").Val())
	assert.Equal(t, "budi", client.LPop(ctx, "list").Val())
	assert.Equal(t, "joko", client.LPop(ctx, "list").Val())

	client.Del(ctx, "list")
}

func TestSet(t *testing.T) {
	client.SAdd(ctx, "students", "dihanto")
	client.SAdd(ctx, "students", "dihanto")
	client.SAdd(ctx, "students", "budi")
	client.SAdd(ctx, "students", "budi")
	client.SAdd(ctx, "students", "joko")
	client.SAdd(ctx, "students", "joko")

	assert.Equal(t, int64(3), client.SCard(ctx, "students").Val())
	assert.Equal(t, []string{"joko", "dihanto", "budi"}, client.SMembers(ctx, "students").Val())
}

func TestSortedSet(t *testing.T) {
	client.ZAdd(ctx, "scores", redis.Z{Score: 10, Member: "dihanto"})
	client.ZAdd(ctx, "scores", redis.Z{Score: 8, Member: "joko"})
	client.ZAdd(ctx, "scores", redis.Z{Score: 9, Member: "budi"})

	assert.Equal(t, []string{"dihanto", "budi", "joko"}, client.ZRevRange(ctx, "scores", 0, -1).Val())
	assert.Equal(t, "dihanto", client.ZPopMax(ctx, "scores").Val()[0].Member)
	assert.Equal(t, "budi", client.ZPopMax(ctx, "scores").Val()[0].Member)
	assert.Equal(t, "joko", client.ZPopMax(ctx, "scores").Val()[0].Member)

}

func TestHash(t *testing.T) {
	client.HSet(ctx, "user:1", "id", "1")
	client.HSet(ctx, "user:1", "name", "dihanto")
	client.HSet(ctx, "user:1", "email", "dihanto@go.dev")

	user := client.HGetAll(ctx, "user:1").Val()

	assert.Equal(t, "1", user["id"])
	assert.Equal(t, "dihanto", user["name"])
	assert.Equal(t, "dihanto@go.dev", user["email"])
	client.Del(ctx, "user:1")
}

func TestGeoPoint(t *testing.T) {
	client.GeoAdd(ctx, "sellers", &redis.GeoLocation{
		Longitude: 107.107042,
		Latitude:  -6.346297,
		Name:      "Toko A",
	})
	client.GeoAdd(ctx, "sellers", &redis.GeoLocation{
		Longitude: 107.096663,
		Latitude:  -6.329913,
		Name:      "Toko B",
	})

	distance := client.GeoDist(ctx, "sellers", "Toko A", "Toko B", "km").Val()
	assert.Equal(t, 2.1536, distance)

	seller := client.GeoSearch(ctx, "sellers", &redis.GeoSearchQuery{
		Longitude:  107.100740,
		Latitude:   -6.337548,
		Radius:     100,
		RadiusUnit: "km",
	}).Val()

	assert.Equal(t, []string{"Toko A", "Toko B"}, seller)

}

func TestHyperLogLog(t *testing.T) {
	client.PFAdd(ctx, "visitors", "dihanto", "budi", "joko")
	client.PFAdd(ctx, "visitors", "usman", "amri", "budi")
	client.PFAdd(ctx, "visitors", "ruli", "dihanto", "joko")

	total := client.PFCount(ctx, "visitors").Val()
	assert.Equal(t, int64(6), total)
}

func TestPipeline(t *testing.T) {
	_, err := client.Pipelined(ctx, func(r redis.Pipeliner) error {
		r.SetEx(ctx, "name", "dihanto", 3*time.Second)
		r.Get(ctx, "name")
		return nil
	})

	assert.Nil(t, err)
	assert.Equal(t, "dihanto", client.Get(ctx, "name").Val())
}

func TestTransaction(t *testing.T) {
	_, err := client.TxPipelined(ctx, func(r redis.Pipeliner) error {
		r.SetEx(ctx, "name", "dihanto", 3*time.Second)
		r.Get(ctx, "name")
		return nil
	})

	assert.Nil(t, err)
	assert.Equal(t, "dihanto", client.Get(ctx, "name").Val())
}

func TestPublishStream(t *testing.T) {
	for i := 0; i < 10; i++ {
		err := client.XAdd(ctx, &redis.XAddArgs{
			Stream: "members",
			Values: map[string]interface{}{"name": "dihanto", "email": "dihanto@go.dev"}}).Err()
		assert.Nil(t, err)
	}
}

func TestCreateConsumerGroup(t *testing.T) {
	client.XGroupCreate(ctx, "members", "group1", "0")
	client.XGroupCreateConsumer(ctx, "members", "group1", "consumer-1")
	client.XGroupCreateConsumer(ctx, "members", "group1", "consumer-2")
}

func TestGetStream(t *testing.T) {
	result := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    "group1",
		Consumer: "consumer-1",
		Streams:  []string{"members", ">"},
		Count:    2,
		Block:    5 * time.Second,
	}).Val()

	for _, stream := range result {
		for _, message := range stream.Messages {
			fmt.Println(message.ID, message.Values)
		}
	}
}

func TestSubscribePubSub(t *testing.T) {
	pubsub := client.Subscribe(ctx, "channel-1")
	defer pubsub.Close()
	for i := 0; i < 10; i++ {
		message, _ := pubsub.ReceiveMessage(ctx)
		fmt.Println(message.Channel, message.Payload)
	}
}

func TestPublishPubSub(t *testing.T) {
	for i := 0; i < 10; i++ {
		err := client.Publish(ctx, "channel-1", "Hello "+strconv.Itoa(i)).Err()
		assert.Nil(t, err)
	}
}
