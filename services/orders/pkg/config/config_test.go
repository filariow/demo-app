package config_test

import (
	"eshop-orders/pkg/config"
	"os"
	"path"
	"testing"
)

func Test_CreateConfig(t *testing.T) {
	var c struct {
		Value1 string `env:"VALUE_1"`
	}

	// arrange
	if err := os.Setenv("VALUE_1", "test"); err != nil {
		t.Error(err)
		t.Fail()
	}

	// act
	config.CreateConfig(&c)

	// assert
	if c.Value1 != "test" {
		t.Error("invalid value found in configuration")
		t.Fail()
	}
}

func Test_ReadConfig(t *testing.T) {
	var c struct {
		Provider struct {
			Value string `sbc-key:"key1"`
		} `sbc-provider:"provider"`
	}

	// arrange
	p := path.Join(os.TempDir(), "demo-soa-tests", "provider")
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		if err := os.RemoveAll(p); err != nil {
			t.Fatalf("error removing folder '%s': %s", p, err)
		}
	}

	err := os.MkdirAll(p, 0755)
	if err != nil {
		t.Fatalf("error creating folder '%s': %s", p, err)
	}
	// defer os.RemoveAll(p)

	fp := path.Join(p, "key1")
	if err := os.WriteFile(fp, []byte("Value1"), 0644); err != nil {
		t.Fatalf("error writing file '%s': %s", fp, err)
	}

	// act
	config.ReadConfig(path.Join(os.TempDir(), "demo-soa-tests"), &c)

	// assert
	if c.Provider.Value != "Value1" {
		t.Error("invalid value found in configuration")
		t.Fail()
	}
}
