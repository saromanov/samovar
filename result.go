package samovar

import (
	"encoding/json"
	"gopkg.in/redis.v3"
	"time"
)

type Result struct {
	Title  string
	Result interface{}
	Date   *time.Time
}

//Store result provides write function result to db
func (res *Result) storeResult(client *redis.Client) {
	results, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	client.HSet("samovar", res.Title, string(results))
}
