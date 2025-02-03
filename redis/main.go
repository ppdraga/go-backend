package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val, _, err := client.Scan(ctx, 0, "fincert_inn:*", 0).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("fincert_inn", val)

}
