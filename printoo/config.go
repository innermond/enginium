package printoo

import (
	"encoding/json"
	"os"
)

type Conf struct {
	Env        string `json:"env"`
	ServerPem  string `json:"serverPem"`
	ServerKey  string `json:"serverKey"`
	ServerAddr string `json:"serverAddr"`
}

func ConfigFrom(path string) (Conf, error) {
	conf := Conf{}
	confFile, err := os.Open(path)
	if err != nil {
		return conf, err
	}
	json.NewDecoder(confFile).Decode(&conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}
