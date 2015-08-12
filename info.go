package samovar
import
(
	"fmt"
	"gopkg.in/redis.v3"
)

//Info provides storing info about job

type Info struct {
	Title string
	Status int
}

//Store info about job
func (res *Info) storeInfo(client *redis.Client) {
	results := Marshal(res)
	client.HSet("samovar", fmt.Sprintf("%s_info", res.Title), string(results))
}