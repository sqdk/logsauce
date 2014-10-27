package logsauce

import (
	"github.com/blakesmith/go-grok/grok"
	"log"
)

//"encoding/json"

func init() {
	g := initGrok()

	text := "Tue May 15 11:21:42 [conn1047685] moveChunk deleted: 7157"
	pattern := "%{DAY}"
	err := g.Compile(pattern)
	if err != nil {
		log.Fatal("Error:", err)
	}
	match := g.Match(text)
	log.Println(match)

	captures := match.Captures()

	log.Println(captures)

}

func initGrok() *grok.Grok {
	g := grok.New()
	g.AddPatternsFromFile("/home/cp/gopath/src/github.com/sqdk/logsauce/logsauce-server/patterns.txt")
	return g
}

func tokenizeLinesForPeriod(hostName, filePath string, startUnix, endUnix int64) {
	lines, err := getLoglinesForPeriodForHostnameAndFilepath(hostName, filePath, startUnix, endUnix)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(lines)

}
