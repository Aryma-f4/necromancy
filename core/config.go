package core

type Config struct {
	Ports       string
	Interface   string
	Connect     string
	Maintain    int
	NoLog       bool
	NoUpgrade   bool
	OSCPSafe    bool
	ServeDir    string
}

var GlobalConfig *Config

func InitConfig() {
	if GlobalConfig == nil {
		GlobalConfig = &Config{}
	}
}
