package logsauce

import (
	"log"
	"time"
)

func calcDistributionOverTimeUnsortedInput(patternName, hostname, filepath string, resolution int, fieldName string, timeStart, timeEnd int64) []ComputeResponse {
	periodLength := timeEnd - timeStart
	intervalSize := periodLength / int64(resolution)

	log.Printf("Calculating distribution with resolution: %v", resolution)

	ts := time.Now()

	lineCount := 0

	timeCount := timeStart + intervalSize
	var responses []ComputeResponse
	for i := 0; i < resolution; i++ {
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
		response.Filename = loglines[0].Filename
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
	}

	log.Printf("Calculation done. \nTime: %v. \nTotal number of records processed: %v. \nAvg. time pr record: %v ns", time.Now().Sub(ts), lineCount, int64(lineCount)/(time.Now().Sub(ts).Nanoseconds()/int64(1000000)))
	return responses
}
