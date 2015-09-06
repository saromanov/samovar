package samovar

import (
	"errors"
	"log"
	"os"
	"fmt"
)

type Logging struct {
	path string
}

func InitLog(path string) *Logging {
	err := createFile(path)
	if err != nil {
		log.Fatal(err)
	}
	l := new(Logging)
	l.path = path
	return l
}

//LogWrite provides default writes to log file
func (l *Logging) LogWrite(msg string) {
	log.Println(msg)
}

//SetNewPath provides change output file for logging
func (l *Logging) SetNewPath(newpath string) error {
	err := createFile(newpath)
	if err != nil {
		return err
	}
	l.path = newpath
	return nil
}

func createFile(path string) error {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.New(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()
	log.SetOutput(f)
	return nil
}
