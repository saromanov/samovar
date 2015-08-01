package samovar

import (
	"encoding/json"
	"fmt"
	"gopkg.in/redis.v3"
	"log"
	"time"
	"hash/crc32"
)

type Result struct {
	Title  string
	Result interface{}
	Date   time.Time
	DataChecksum uint32
}

//Store result provides write function result to db
func (res *Result) storeResult(client *redis.Client) {
	results, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	client.HSet("samovar", res.Title, string(results))
}

func getResult(client *redis.Client, title string) interface{} {
	var res Result
	result, err := client.HGet("samovar", title).Result()
	if err != nil {
		log.Printf(fmt.Sprintf("Job by title %s not found", title))
	}
	errunm := json.Unmarshal([]byte(result), &res)
	if errunm != nil {
		panic(err)
	}
	return res

}

func getChecksum(item []byte) uint32 {
	crc := crc32.New(crc32.MakeTable(crc32.Castagnoli))
    crc.Write([]byte(item))
    return crc.Sum32()
}
