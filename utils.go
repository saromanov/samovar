package samovar

import (
	"fmt"
	"os"
	"os/signal"
	"encoding/json"
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
