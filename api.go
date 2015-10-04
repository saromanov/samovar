package samovar

import (
	"./backend"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type JobShow struct {
}

func StartServer(jobs *Jobs) {
	r := mux.NewRouter()
	client := new(Client)
	client.client = backend.InitRedis("localhost:6379")

	//Get information about job
	r.HandleFunc("/statjob/{title}", func(w http.ResponseWriter, req *http.Request) {
		title := mux.Vars(req)["title"]
		var job []*Job
		err := jobs.GetJob(title, &job)
		if err != nil {
			//fmt.Printf(err)
		}
		//w.Write([]byte(fmt.Sprintf("%d", job.numberofcalls)))
		return
	})

	//Start new job
	r.HandleFunc("/jobrun/{title}", func(w http.ResponseWriter, req *http.Request) {
		title := mux.Vars(req)["title"]
		if title == "" {
			return
		}

		client.Send(title, &JobOptions{})
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
