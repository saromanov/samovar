package samovar
import
(
	"fmt"
	"gopkg.in/redis.v3"
	"errors"
	"encoding/json"
)

//Info provides storing info about job

const (
	Completed = 1
	Waiting = 2
)

type Info struct {
	Title string
	JobID string
	Status int
}

//Store info about job
func (res *Info) storeInfo(client *redis.Client) {
	results := Marshal(res)
	client.HSet("samovar", fmt.Sprintf("%s%s_info", res.Title, res.JobID), string(results))
}

func getInfo(client *redis.Client, title string)(Info, error) {
	var info Info
	result, err := client.HGet("samovar", fmt.Sprintf("%s_info", title)).Result()
	if err != nil {
		return Info{}, errors.New(fmt.Sprintf("Job by title %s not found", title))
	}
	errunm := json.Unmarshal([]byte(result), &info)
	if errunm != nil {
		return Info{}, errunm
	}

	return info, nil
}