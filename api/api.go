package api

import (
	"../backend"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
	"net/http"
	"fmt"
	"encoding/json"
)

const (
	Serveraddr = ":8000"
)

type JobShow struct {
}

func StartServer() {
	r := mux.NewRouter()
	back := backend.InitRedis("localhost:6379")

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

		res := client.GetResult(title)
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
