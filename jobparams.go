package samovar

import (
	"encoding/json"
)

//JobParams provides struct set to queue
type JobParams struct {
	Name      string
	Arguments []interface{}
	Delay     uint
	Period    uint

	//HardDeadline represents time where task must be complete
	Hardeadline   float64
}


//getJobArguments needs for convert item data for the JobParams (worker area)
func getJobArguments(item string) JobParams {
	var params JobParams
	err := json.Unmarshal([]byte(item), &params)
	if err != nil {
		panic(err)
	}

	return params
}

//prepareParemeters converts JobParams object to the string
func prepareParameters(jp *JobParams) string {
	results, err := json.Marshal(jp)
	if err != nil {
		panic(err)
	}

	return string(results)
}
