package samovar

import (
	"fmt"
	"os"
	"os/signal"
	"encoding/json"
	"os/exec"
	"math/rand"
	"time"
)

func detectExit() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		fmt.Println("Good bye!")
		os.Exit(1)
	}()
}

func Marshal(res interface{}) string {
	results, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	return string(results)
}

//Idgen provides generation of unique id for job result
func Idgen() string {
	res, err := exec.Command("uuidgen").Output()
	if err != nil {
		result := ""
		for i := 0; i < 16; i++ {
			rand.Seed(time.Now().UTC().UnixNano())
			result += string(1+rand.Intn(9))
		}
		return result
	}

	return string(res)
}
