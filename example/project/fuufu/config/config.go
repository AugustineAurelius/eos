package config

import "strconv"

type Config struct {
	PostgresConfig  Postgres
	PostgresSlave   Postgres
	LoggerConfig    Logger
	CollectorConfig Collector
	SeverConfig     Server
}

// Collector
type Collector struct {
	Host string
	Port int
}

func (c Collector) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func (man Manager) LoadCollector() Collector {
	return man.LoadConfig().CollectorConfig
}

// Logger
type Logger struct {
	Debug bool `yaml:"debug"`
	JSON  bool `yaml:"json"`
}

func (man Manager) LoadLogging() Logger {
	return man.LoadConfig().LoggerConfig
}

// Server
type Server struct {
	Addr                  string
	AuthMiddlewareExclude []string
	GeoMiddlewareExclude  []string
}

func (man Manager) LoadServer() Server {
	return man.LoadConfig().SeverConfig
}
