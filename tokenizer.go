package logsauce

func tokenizeLines(hostId int64, filePath string) {
	lines := getLoglinesForFileForHost(hostId, filePath)
}
