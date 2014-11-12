package logsauce

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Configuration struct {
	ListenPort int  `json:"listen_port"`
	ClientMode bool `json:"client_mode"`
	ServerMode bool `json:"server_mode"`
	Relaymode  bool `json:"relay_mode"`

	ServerConfiguration struct {
		DbAddress            string `json:"db_address"`
		DbName               string `json:"db_name"`
		DbUsername           string `json:"db_username"`
		DbPassword           string `json:"db_password"`
		ServerCertificate    string `json:"server_cert"`
		ServerCertificateKey string `json:"server_key"`
		Hosts                []Host `json:"hosts"`
	} `json:"server_config"`

	ClientConfiguration struct {
		ServerAddress     string   `json:"server_address"`
		FilesToWatch      []string `json:"files_to_watch"`
		ServerCertificate string   `json:"server_cert"`
		ClientToken       string   `json:"token"`
	} `json:"client_config"`
}

func ReadConfig(path string) (Configuration, error) {
	var conf Configuration
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}

	err = json.Unmarshal(bytes, &conf)

	return conf, err
}

func WriteExampleServerConfig() {
	var conf Configuration
	conf.ListenPort = 6214
	conf.ServerMode = true

	conf.ServerConfiguration.DbAddress = "127.0.0.1:3306"
	conf.ServerConfiguration.DbName = "exampleDb"
	conf.ServerConfiguration.DbUsername = "user"
	conf.ServerConfiguration.DbPassword = "password"

	bytes, _ := json.MarshalIndent(&conf, "", "    ")

	ioutil.WriteFile("server.conf.example", bytes, os.ModeAppend)
}
