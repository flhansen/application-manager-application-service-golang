package controller

import (
	"fmt"
)

type DbConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func (conf DbConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
}
