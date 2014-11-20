package logsauce

import (
	"fmt"
	"github.com/mbanzon/cache"
	"log"
	"time"
)

var responseCache *cache.ICache

func init() {
	responseCache = cache.NewICache(0, 0)
}
func calcDistributionOverTime(patternName, hostname, filepath string, resolution int, fieldName string, timeStart, timeEnd int64) []ComputeResponse {
	periodLength := timeEnd - timeStart
	intervalSize := periodLength / int64(resolution)

	log.Printf("Calculating distribution with resolution: %v", resolution)

	//ts := time.Now()

	lineCount := 0

	timeCount := timeStart + intervalSize
	var responses []ComputeResponse
	for i := 0; i < resolution; i++ {
		reqString := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v", patternName, hostname, filepath, timeCount-intervalSize, timeCount, fieldName, "dist")
		if responseCache.Has(reqString) {
			cachedResponse, exists := responseCache.Get(reqString)
			if exists == true {
				log.Println("Found cached response: ", reqString)
				responses = append(responses, cachedResponse.(ComputeResponse))
				timeCount += intervalSize
				continue
			} else {
				log.Println("Not cached: ", reqString)
			}
		}

		loglines, err := getTokenizedDataForHostnameAndFilepath(patternName, hostname, filepath, timeCount-intervalSize, timeCount)
		lineCount += len(loglines)
		if err != nil {
			log.Println(err)
		}
		if len(loglines) == 0 {
			continue
		}

		var response ComputeResponse
		response.Values = make(map[string]int)
		response.TimeStart = timeCount
		response.IntervalSize = intervalSize
		response.FieldName = fieldName
		response.Filename = loglines[0].Filepath
		response.HostId = loglines[0].HostId

		if err != nil {
			log.Println(err)
			return []ComputeResponse{}
		}
		//Can be optimized if loglines are sorted by timestamp
		for k := 0; k < len(loglines); k++ {
			if loglines[k].Timestamp < timeCount && loglines[k].Timestamp >= timeCount-intervalSize {
				response.Values[loglines[k].TokenizedObject[fieldName]] += 1
			}
		}

		responses = append(responses, response)
		timeCount += intervalSize
		cacheResponse(reqString, response)
	}

	//log.Printf("Calculation done. \nTime: %v. \nTotal number of records processed: %v. \nAvg. time pr record: %v ns", time.Now().Sub(ts), lineCount, int64(lineCount)/(time.Now().Sub(ts).Nanoseconds()/int64(1000000)))
	return responses
}

func countUniqueOverTime(patternName, hostname, filepath string, resolution int, fieldName string, timeStart, timeEnd int64) []ComputeResponse {
	periodLength := timeEnd - timeStart
	intervalSize := periodLength / int64(resolution)
	var responses []ComputeResponse
	lineCount := 0

	log.Printf("Calculating distribution with resolution: %v", resolution)

	ts := time.Now()

	intermediateResponses := calcDistributionOverTime(patternName, hostname, filepath, resolution, fieldName, timeStart, timeEnd)

	for i := 0; i < len(intermediateResponses); i++ {
		uniqueCount := 0
		for _ = range intermediateResponses[i].Values {
			uniqueCount += 1
		}
		var response ComputeResponse
		response.Values = make(map[string]int)
		response.TimeStart = intermediateResponses[i].TimeStart
		response.IntervalSize = intervalSize
		response.FieldName = fieldName
		response.Filename = filepath

		response.Values["count"] = uniqueCount
		responses = append(responses, response)
		lineCount += 1
	}

	log.Printf("Processing done. \nTime: %v. \nTotal number of records processed: %v. \nAvg. time pr record: %v ns", time.Now().Sub(ts), lineCount, int64(lineCount)/(time.Now().Sub(ts).Nanoseconds()/int64(1000000)))
	return responses
}

func cacheResponse(reqString string, response ComputeResponse) {
	//If response could potentially change as data is generated, do not cache
	if response.TimeStart+response.IntervalSize > time.Now().Unix() {
		return
	}

	responseCache.Store(reqString, response)
}
