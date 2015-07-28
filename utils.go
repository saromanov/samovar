package samovar

import
(
	"os"
	"os/signal"
	"fmt"
)

func detectExit(){
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func(){
        <-sig
        fmt.Println("Good bye!")
        os.Exit(1)
	}()
}