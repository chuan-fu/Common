package configor

import (
	"github.com/chuan-fu/Common/zlog"
	jConfigor "github.com/jinzhu/configor"
)

func Load(config interface{}, files ...string) error {
	return jConfigor.Load(config, files...)
}

func MustLoad(config interface{}, files ...string) {
	if err := jConfigor.Load(config, files...); err != nil {
		log.Fatal(err)
	}
}

/*

var Config = struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
}{}

func main() {
	configor.Load(&Config, "config.yml")
	fmt.Printf("config: %#v", Config)
}
*/

/*

appname: test

db:
name:     test
user:     test
password: test
port:     1234

contacts:
- name: i test
email: test@test.com
*/
