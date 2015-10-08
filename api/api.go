package api

import (
	"../backend"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
	"net/http"
	"encoding/json"
	"time"
)

const (
	Serveraddr = ":8000"
)

type JobShow struct {
}

type Result struct {
	ID           string
	Title        string
	Result       interface{}
	Date         time.Time
	DataChecksum uint32
	Status       uint32
}

//JobParams provides struct set to queue
type JobParams struct {
	Name      string
	Arguments interface{}
	Delay     uint
	Period    uint
	//HardDeadline represents time where task must be complete
	Hardeadline   float64
}

//prepareParemeters converts JobParams object to the string
func prepareParameters(jp *JobParams) string {
	results, err := json.Marshal(jp)
	if err != nil {
		panic(err)
	}

	return string(results)
}


func StartServer() {
	r := mux.NewRouter()
	back := backend.InitRedisBackend()
	rediscli := backend.InitRedis("localhost:6379")

	//Get information about job
	r.HandleFunc("/statjob/{title}", func(w http.ResponseWriter, req *http.Request) {
		/*title := mux.Vars(req)["title"]
		var job []*Job
		err := jobs.GetJob(title, &job)
		if err != nil {
			//fmt.Printf(err)
		}
		//w.Write([]byte(fmt.Sprintf("%d", job.numberofcalls)))
		return*/
	})

	//Start new job
	r.HandleFunc("/jobrun/{title}", func(w http.ResponseWriter, req *http.Request) {
		title := mux.Vars(req)["title"]
		if title == "" {
			return
		}

		back.PublishJob(prepareParameters(&JobParams{
				Name:      title,
		}), "default")

		w.Write([]byte(fmt.Sprintf("Job %s was sending", title)))
	})

	r.HandleFunc("/jobresult/{title}", func(w http.ResponseWriter, req *http.Request) {
		title := mux.Vars(req)["title"]
		if title == "" {
			return
		}

		res, err := getResult(rediscli, title)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error for getting result")))
			return
		}
		w.Write([]byte(fmt.Sprintf("%s", res)))
	})

	http.ListenAndServe(":8000", r)
}

func getResult(client *redis.Client, title string) (interface{}, error) {
	var res Result
	result, err := client.HGet("samovar", fmt.Sprintf("%s_result", title)).Result()
	if err != nil {
		return nil, fmt.Errorf("Job by title %s not found", title)
	}
	errunm := json.Unmarshal([]byte(result), &res)
	if errunm != nil {
		return nil, errunm
	}
	return res, nil

}
