package samovar

import
(
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
)


func StartServer(jobs *Jobs) {
	r := mux.NewRouter()
	r.HandleFunc("/job/{title}", func(w http.ResponseWriter, req *http.Request){
		title := mux.Vars(req)["title"]
		var job Job
		err := jobs.GetJob(title, &job)
		if err != nil {
			//fmt.Printf(err)
		}
		w.Write([]byte(fmt.Sprintf("%d", job.numberofcalls)))
		return
	})
	http.ListenAndServe(":8000", r)
}