package lib

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"net/url"
	"os"
	"path"
	"reflect"
)

type env struct {
	DONAR_PORT,
	PROXY_PORT,
	PROXY_URL,
	PLACE string
	ProxyUrl *url.URL
}

func (i *env) Init() error {

	pwd, err := os.Getwd()

	fmt.Println(pwd)

	if err != nil {
		return err
	}

	envPath := path.Join(pwd, ".env")

	if _, err := os.Stat(envPath); err == nil {
		fmt.Println("Env load from file")
		err := godotenv.Load(envPath)

		if err != nil {
			return err
		}
	}

	props := map[string]bool{
		"DONAR_PORT": true,
		"PROXY_PORT": true,
		"PROXY_URL":  true,
		"PLACE":      true,
	}

	for prop, isRequired := range props {

		v := os.Getenv(prop)

		if isRequired == true && v == "" {
			return errors.New("Bad " + prop)
		}

		reflect.ValueOf(i).Elem().FieldByName(prop).SetString(v)
	}

	if i.ProxyUrl, err = url.Parse(i.PROXY_URL); err != nil {
		return err
	}

	if !FileExists("logs") {
		err = os.Mkdir("logs", 0777)

		if err != nil {
			return err
		}
	}

	return nil
}

var ENV env
