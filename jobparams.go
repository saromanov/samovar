package samovar

import (
	"encoding/json"
)

//JobParams provides struct set to queue
type JobParams struct {
	Name      string
	Arguments interface{}
	Delay     uint
}

func GetJobArguments(item string) JobParams {
	var params JobParams
	err := json.Unmarshal([]byte(item), &params)
	if err != nil {
		panic(err)
	}

	return params
}

func PrepareParameters(jp *JobParams) string {
	results, err := json.Marshal(jp)
	if err != nil {
		panic(err)
	}

	return string(results)
}
