package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Postgres struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	PasswordPath    string `yaml:"password_path"`
	Database        string `yaml:"database"`
	SSLMode         string `yaml:"ssl_mode"`
	SSLCert         string `yaml:"ssl_cert"`
	SSLKey          string `yaml:"ssl_key"`
	SSLRootCert     string `yaml:"ssl_root_cert"`
	SSLServerName   string `yaml:"ssl_server_name"`
	MaxConns        int    `yaml:"max_conns"`
	MinConns        int    `yaml:"min_conns"`
	MaxConnLifetime int    `yaml:"max_conn_lifetime"`
}

func (c Postgres) DNS() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.Username, c.Password, c.Database)
}

func (c Postgres) URL() string {
	password := c.Password
	if password == "" && c.PasswordPath != "" {
		data, err := os.ReadFile(c.PasswordPath)
		if err != nil {
			panic(fmt.Errorf("failed to read password file: %w", err))
		}
		password = strings.TrimSpace(string(data))
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.Username, password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   c.Database,
	}

	q := u.Query()

	if c.SSLMode != "" {
		q.Set("sslmode", c.SSLMode)
	}
	if c.SSLCert != "" {
		q.Set("sslcert", c.SSLCert)
	}
	if c.SSLKey != "" {
		q.Set("sslkey", c.SSLKey)
	}
	if c.SSLRootCert != "" {
		q.Set("sslrootcert", c.SSLRootCert)
	}
	if c.SSLServerName != "" {
		q.Set("sslservername", c.SSLServerName)
	}

	if c.MaxConns > 0 {
		q.Set("pool_max_conns", strconv.Itoa(c.MaxConns))
	}
	if c.MinConns > 0 {
		q.Set("pool_min_conns", strconv.Itoa(c.MinConns))
	}
	if c.MaxConnLifetime > 0 {
		q.Set("pool_max_conn_lifetime", fmt.Sprintf("%ds", c.MaxConnLifetime))
	}

	u.RawQuery = q.Encode()

	return u.String()
}

func (man Manager) LoadPostgres() Postgres {
	return man.LoadConfig().PostgresConfig
}

func (man Manager) LoadPostgresSlave() Postgres {
	return man.LoadConfig().PostgresSlave
}
