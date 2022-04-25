package email

import "gopkg.in/gomail.v2"

type EmailConf struct {
	Host     string `required:"true" json:"host" yaml:"host"`
	Port     int    `required:"true" json:"port" yaml:"port"`
	Username string `required:"true" json:"username" yaml:"username"`
	Password string `required:"true" json:"password" yaml:"password"`
}

var globalClient *gomail.Dialer

func NewClient(conf EmailConf) error {
	globalClient = gomail.NewDialer(conf.Host, conf.Port, conf.Username, conf.Password)
	return nil
}

func ReloadClient(conf EmailConf) error {
	return NewClient(conf)
}

func GetClient() *gomail.Dialer {
	return globalClient
}

func NewLocalClient(host string, port int, username, password string) *gomail.Dialer {
	return gomail.NewDialer(host, port, username, password)
}
