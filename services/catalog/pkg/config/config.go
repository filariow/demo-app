package config

import (
	"eshop-catalog/pkg/awsconfig"
	"eshop-catalog/pkg/persistence"
	"eshop-catalog/pkg/queue"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
)

const (
	EnvServiceBindingRoot = "SERVICE_BINDING_ROOT"
)

type Config struct {
	Aws      awsconfig.Config             `sbc-provider-env:"AWS_PROVIDER" sbc-provider:"aws"`
	Postgres persistence.PostgreSQLConfig `sbc-provider-env:"RDS_PROVIDER" sbc-provider:"postgres"`
	SQS      queue.SQSConfig              `sbc-provider-env:"SQS_PROVIDER" sbc-provider:"sqs"`
}

func NewConfigFromServiceBinding() Config {
	c := Config{}
	ReadConfig(os.Getenv(EnvServiceBindingRoot), &c)
	return c
}

func ReadConfig(basePath string, configPtr interface{}) {
	sv := reflect.ValueOf(configPtr).Elem()
	v := sv.Type()

	for i := 0; i < sv.NumField(); i++ {
		te := v.Field(i).Tag.Get("sbc-provider-env")
		if te != "" {
			k, ok := os.LookupEnv(te)
			if !ok {
				log.Printf("error fetching sbc-provider from env var '%s': env var not found", te)
			} else {
				c := reflect.New(v.Field(i).Type)
				readProviderConfig(basePath, k, c)
				sv.Field(i).Set(c.Elem())
				continue
			}
		}

		t := v.Field(i).Tag.Get("sbc-provider")
		if t != "" {
			c := reflect.New(v.Field(i).Type)
			readProviderConfig(basePath, t, c)
			sv.Field(i).Set(c.Elem())
		}
	}
}

func readProviderConfig(basePath string, provider string, cv reflect.Value) {
	sv := cv.Elem()
	v := sv.Type()

	for i := 0; i < sv.NumField(); i++ {
		t := v.Field(i).Tag.Get("sbc-key")
		if t != "" {
			k, err := readProviderKey(basePath, provider, t)
			if err != nil {
				log.Printf("error reading key '%s/%s/%s': %s", basePath, provider, t, err)
				continue
			}
			sv.Field(i).Set(reflect.ValueOf(*k))
		}
	}
}

func listProviderKeys(basePath, provider string) ([]string, error) {
	p := path.Join(basePath, provider)
	ii, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("error listing file in directory '%s': %w", p, err)
	}

	nn := []string{}
	for _, i := range ii {
		if !i.IsDir() {
			nn = append(nn, i.Name())
		}
	}
	return nn, err
}

func readProviderKey(basePath, provider, key string) (*string, error) {
	p := path.Join(basePath, provider, key)
	bb, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %w", p, err)
	}

	d := strings.TrimRight(string(bb), "\n")
	return &d, nil
}
