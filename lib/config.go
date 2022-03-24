package lib

type Config struct {
	token    string
	debug    bool
	readOnly bool
}

var config Config

func EnableDebug() {
	config.debug = true
}

func EnableDryRunMode() {
	config.readOnly = true
}

func SetToken(t string) {
	config.token = t
}
