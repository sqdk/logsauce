package logsauce

type Configuration struct {
	ListenPort    int64  `json:"listen_port"`
	ServerAddress string `json:"server_address"`
	ClientMode    bool   `json:"client_mode"`
	ServerMode    bool   `json:"server_mode"`
	Relaymode     bool   `json:"relay_mode"`

	ServerConfiguration struct {
		DbAddress            string `json:"db_address"`
		DbUsername           string `json:"db_username"`
		DbPassword           string `json:"db_password"`
		AccessToken          string `json:"access_token"`
		ServerCertificate    string `json:"server_cert"`
		ServerCertificateKey string `json:"server_key"`
	} `json:"server_config"`

	ClientConfiguration struct {
		FilesToWatch      []string `json:"files_to_watch"`
		ServerCertificate string   `json:"server_cert"`
	} `json:"client_config"`
}
