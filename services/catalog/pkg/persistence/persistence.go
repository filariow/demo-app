package persistence

import "fmt"

type PostgreSQLConfig struct {
	Host     string `sbc-key:"host"`
	Username string `sbc-key:"username"`
	Password string `sbc-key:"password"`
	Port     string `sbc-key:"port"`
	Database string `sbc-key:"database"`
}

func (c PostgreSQLConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.Username, c.Password, c.Host, c.Port, c.Database)
}
