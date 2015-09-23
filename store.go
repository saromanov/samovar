package samovar
import
(
	"gopkg.in/redis.v3"
	"encoding/json"
	"log"
)

const (
	STORE_TITLE = "samovar_store"
)

type Store struct {
	client  *redis.Client
}

func InitStore()*Store {
	st := new(Store)
	st.client = initRedis("localhost:6379")
	return st
}

func (store *Store) Set(key, value string){
	store.client.HSet(STORE_TITLE, key, value)
}

func (store *Store) Get(key string) Result {
	value, err := store.client.HGet(STORE_TITLE, key).Result()
	if err != nil {
		log.Fatal(err)
	}

	var res Result
	errunm := json.Unmarshal([]byte(value), &res)
	if errunm != nil {
		log.Fatal(err)
	}

	return res
}