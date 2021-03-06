package samovar

import (
	"encoding/json"
	"fmt"
	"gopkg.in/redis.v3"
	"hash/crc32"
	"log"
	"time"
	"errors"
)


const (
	STARTED = iota
	FINISHED
	FAILED 
)

type Result struct {
	ID           string
	JobID        string
	Title        string
	Result       interface{}
	Date         time.Time
	DataChecksum uint32
	Status       uint32
}

//Store result provides write function result to db
func (res *Result) storeResult(client *redis.Client) {
	err := res.checkResultItem()
	if err != nil {
		log.Fatal(err)
	}

	results := Marshal(res)
	client.HSet("samovar", fmt.Sprintf("%s%s_result", res.Title, res.JobID), string(results))
}

//storeResultbyId provides setting only id and reslt of the job
func (res *Result) storeResultById(client *redis.Client) {
	value := Marshal(res.Result)
	client.HSet("samovar", fmt.Sprintf("%s_idresult", res.ID), string(value))
}

//return result of the job by id
func getResultById(client *redis.Client, id string) (interface{}, error) {
	var res interface{}
	result, err := client.HGet("samovar", fmt.Sprintf("%s_idresult", id)).Result()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ID %s not found", id))
	}

	errunm := json.Unmarshal([]byte(result), &res)
	if errunm != nil {
		return nil, errunm
	}
	return res, nil
}

func getResult(client *redis.Client, title string) (interface{}, error) {
	var res Result
	result, err := client.HGet("samovar", fmt.Sprintf("%s_result", title)).Result()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Job by title %s not found", title))
	}
	errunm := json.Unmarshal([]byte(result), &res)
	if errunm != nil {
		return nil, errunm
	}
	return res, nil

}

func getChecksum(item []byte) uint32 {
	crc := crc32.New(crc32.MakeTable(crc32.Castagnoli))
	crc.Write([]byte(item))
	return crc.Sum32()
}

//This helful method returns error if the fields of Result is invalid
func (res*Result) checkResultItem() error {
	if res.Title == " " || res.Title == "" {
		return errors.New("Title is empty string")
	}

	if res.ID == "" {
		return errors.New("ID is empty")
	}

	return nil
}
