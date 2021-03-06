package logsauce

import (
	"bufio"
	"code.google.com/p/go.exp/fsnotify"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func WatchFiles(filesToWatch []string, loadExistingData bool, serverAddress, clientToken string) {
	executionDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadFile(executionDir + "datawatch.state")
	var fileIndex map[string]int64

	if err != nil {
		fileIndex = make(map[string]int64)
	} else {
		err := json.Unmarshal(b, &fileIndex)
		if err != nil {
			log.Fatal("State information corrupt. Delete datawatch.state file or repair")
		}
	}

	ticker := time.NewTicker(time.Second * 1)
	go func(fileIndexes *map[string]int64) {
		for {
			select {
			case <-ticker.C:
				log.Println("Writing state")
				b, _ := json.Marshal(fileIndex)
				err = ioutil.WriteFile(executionDir+"datawatch.state", b, 0644)
				if err != nil {
					log.Println("Unable to write state information to disk")
				}
			}
		}

		println("Timer expired")
	}(&fileIndex)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(filesToWatch); i++ {
		err := watcher.Watch(filesToWatch[i])
		fileIndex[filesToWatch[i]] = 0
		if err != nil {
			log.Println(err)
		}
	}

	if loadExistingData {
		for i := 0; i < len(filesToWatch); i++ {
			file, err := os.Open(filesToWatch[i])
			if err != nil {
				log.Println(err)
				break
			}

			scanner := bufio.NewScanner(file)
			length := int64(0)

			for scanner.Scan() {
				length += int64(len(scanner.Text()))
				SendLogLine(serverAddress, scanner.Text(), filesToWatch[i], "", clientToken)
			}

			fileIndex[filesToWatch[i]] = length

		}
	}

	log.Println("Watching")

	type KV struct {
		K string
		V int64
	}

	filechangeschan := make(chan *KV)

	go func(fileChangesChannel chan *KV) {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)

				if ev.IsDelete() {
					fileIndex[ev.Name] = 0 //Reset byte counter

					err := errors.New("")
					for err != nil {
						log.Println("Waiting for file to rotate")
						err = watcher.Watch(ev.Name)
						time.Sleep(1 * time.Second)
					}
				}

				if ev.IsModify() {
					file, err := os.Open(ev.Name)
					if err != nil {
						log.Println(err)
						break
					}

					data, err := ioutil.ReadAll(file)
					farthestNewline := int64(0)

					//Advance cursor
					for i := len(data); int64(i) > 0; i-- {
						if rune(data[i-1]) == '\n' {
							farthestNewline = int64(i - 1)
							break
						}

						if fileIndex[ev.Name] == int64(i-1) {
							farthestNewline = int64(len(data) - 1)
							break
						}
					}

					var newData []byte
					if fileIndex[ev.Name] == farthestNewline {
						newData = []byte{}
					} else {
						newData = data[fileIndex[ev.Name]:farthestNewline]
					}

					log.Println(string(newData))

					fileIndex[ev.Name] = farthestNewline
				}

			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}(filechangeschan)

	for {
		select {
		case <-filechangeschan:

		}
	}

	looper := make(chan int)
	<-looper
}

func WatchFile(out chan<- string, fileToWatch string) {
	file, err := os.Open(fileToWatch)
	if err != nil {
		log.Println(err)
		return
	}

	//currentLine := 0

	c := time.Tick(1 * time.Second)
	for _ = range c {
		fileInfo, err := file.Stat()
		if err != nil {
			log.Println(err)
			break
		}

		log.Printf("%s\n", fileInfo.Size())
	}
}
