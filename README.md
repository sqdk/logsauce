logsauce
========

Simple but powerful self-hosted logging system built in Go. Based on the principle of recording the raw log data first, and analyzing it later. The system consists of a server responsible for sourcing all data, and a remote agent responsible for monitoring logfiles.

Analyzing/tokenizing data is done with predefined grok patterns and saved for visualization or report generation at a later point.
