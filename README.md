logsauce
========

Simple but powerful self-hosted logging system built in Go. Based on the principle of recording raw data first, and analyzing it later. The opposite of how logstash works. The system consists of a server responsible for sourcing all data and a remote agent responsible for monitoring data containers (text files, databases, etc.).

Depending on the format of the sourced data, grok can be used to tokenize the data to make mapreduce operations possible. This makes it possible to easily monitor logfiles with a specific pattern (apache, iptables, auth). If the data is pure JSON, tokenization will not be necessary.

A client API will give frontend developers a way to do simple mapreduce operations on the tokenized data to display it in different ways.
