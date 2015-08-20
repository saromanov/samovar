package samovar

import
(
	"net/http"
	"github.com/gorilla/mux"
	//"fmt"
)


type JobShow struct {

}

func StartServer(jobs *Jobs) {
	r := mux.NewRouter()
	client := new(Client)
	client.backend = InitRedisBackend()
	client.client = initRedis("localhost:6379")

	//Get information about job
	r.HandleFunc("/statjob/{title}", func(w http.ResponseWriter, req *http.Request){
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
	})
	http.ListenAndServe(":8000", r)
}