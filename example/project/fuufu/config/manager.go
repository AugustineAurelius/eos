package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func flagNameFromConfigKey(key string) string {
	return strings.ReplaceAll(key, ".", "_")
}

func getFlagUsage(key string, usage string) string {
	return fmt.Sprintf("Env: %s\n\t\t%s", envNameFromConfigKey(key), usage)
}
func envNameFromConfigKey(key string) string {
	return "FUUFU_" + strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}

type Manager struct {
	viper    *viper.Viper
	command  *cobra.Command
	defaults map[string]interface{}
}

func NewManager(command *cobra.Command) *Manager {
	man := &Manager{
		viper:    viper.New(),
		command:  command,
		defaults: map[string]interface{}{},
	}
	man.addConfigs()
	return man
}

func (man Manager) addConfigs() {
	addPgxPoolConfig := func(prefix, defaultAddr, usageSuffix string) {
		var host string
		var port int

		parts := strings.Split(defaultAddr, ":")
		host = parts[0]
		p, _ := strconv.Atoi(parts[1])
		port = p

		man.addConfigString(prefix+".host", host, "PostgreSQL server host"+usageSuffix)
		man.addConfigInt(prefix+".port", port, "PostgreSQL server port"+usageSuffix)
		man.addConfigString(prefix+".username", "postgres", "PostgreSQL server username"+usageSuffix)
		man.addConfigString(prefix+".password", "postgres",
			"PostgreSQL server password (prefer env variable for security)"+usageSuffix)
		man.addConfigString(prefix+".password_path", "", "Path to file containing PostgreSQL server password"+usageSuffix)
		man.addConfigString(prefix+".database", "postgres", "PostgreSQL database name"+usageSuffix)
		man.addConfigString(prefix+".ssl_mode", "disable",
			"PostgreSQL SSL mode (disable|allow|prefer|require|verify-ca|verify-full)"+usageSuffix)
		man.addConfigString(prefix+".ssl_cert", "", "Path to PostgreSQL SSL client certificate"+usageSuffix)
		man.addConfigString(prefix+".ssl_key", "", "Path to PostgreSQL SSL client key"+usageSuffix)
		man.addConfigString(prefix+".ssl_root_cert", "", "Path to PostgreSQL root CA certificate"+usageSuffix)
		man.addConfigString(prefix+".ssl_server_name", "", "PostgreSQL TLS server name for SNI"+usageSuffix)
		man.addConfigInt(prefix+".max_conns", 50, "Maximum number of database connections in pool"+usageSuffix)
		man.addConfigInt(prefix+".min_conns", 5, "Minimum number of idle database connections"+usageSuffix)
		man.addConfigInt(prefix+".max_conn_lifetime", 10, "Maximum lifetime of database connection (seconds)"+usageSuffix)
	}

	addPgxPoolConfig("postgres", "localhost:5432", ". Applies to PostgreSQL connection.")
	addPgxPoolConfig("postgres.slave", "localhost:5433", ". Applies to PostgreSQL read-only connection.")

	man.addConfigBool("logging.debug", true, "Enable debug logging")
	man.addConfigBool("logging.json", false, "Prints all logs in json format")

	man.addConfigString("collector.host", "localhost", "Host for otel collector")
	man.addConfigInt("collector.port", 4317, "Port for otel collector")

	man.addConfigString("server.address", "0.0.0.0:7070", "default server address")
}

func (man Manager) LoadConfig() Config {
	man.loadConfigFile()

	loadPgxPoolConfig := func(prefix string) Postgres {
		return Postgres{
			Host:            man.getConfigString(prefix + ".host"),
			Port:            man.getConfigInt(prefix + ".port"),
			Username:        man.getConfigString(prefix + ".username"),
			Password:        man.getConfigString(prefix + ".password"),
			PasswordPath:    man.getConfigString(prefix + ".password_path"),
			Database:        man.getConfigString(prefix + ".database"),
			SSLMode:         man.getConfigString(prefix + ".ssl_mode"),
			SSLCert:         man.getConfigString(prefix + ".ssl_cert"),
			SSLKey:          man.getConfigString(prefix + ".ssl_key"),
			SSLRootCert:     man.getConfigString(prefix + ".ssl_root_cert"),
			SSLServerName:   man.getConfigString(prefix + ".ssl_server_name"),
			MaxConns:        man.getConfigInt(prefix + ".max_conns"),
			MinConns:        man.getConfigInt(prefix + ".min_conns"),
			MaxConnLifetime: man.getConfigInt(prefix + ".max_conn_lifetime"),
		}
	}
	cfg := Config{
		PostgresConfig: loadPgxPoolConfig("postgres"),
		PostgresSlave:  loadPgxPoolConfig("postgres.slave"),
		LoggerConfig: Logger{
			Debug: man.getConfigBool("logging.debug"),
			JSON:  man.getConfigBool("logging.json"),
		},
		CollectorConfig: Collector{
			Host: man.getConfigString("collector.host"),
			Port: man.getConfigInt("collector.port"),
		},
		SeverConfig: Server{
			Addr: man.getConfigString("server.address"),
		},
	}
	return cfg
}

func (man Manager) loadConfigFile() {
	man.viper.SetConfigType("yaml")

	configFile := man.command.PersistentFlags().Lookup("config").Value.String()

	if configFile == "" {
		return
	}

	man.viper.SetConfigFile(configFile)
	err := man.viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error loading config file:", err)
		os.Exit(1)
	}

	fmt.Println("Using config file:", man.viper.ConfigFileUsed())
}

func (man Manager) addConfigString(key, defVal, usage string) {
	man.command.PersistentFlags().String(flagNameFromConfigKey(key), defVal, getFlagUsage(key, usage))
	man.viper.BindPFlag(key, man.command.PersistentFlags().Lookup(flagNameFromConfigKey(key))) //nolint:errcheck
	man.viper.BindEnv(key, envNameFromConfigKey(key))                                          //nolint:errcheck

	man.addDefault(key, defVal)
}

func (man Manager) getConfigString(key string) string {
	interfaceVal := man.getInterfaceVal(key)
	stringVal, err := cast.ToStringE(interfaceVal)
	if err != nil {
		panic("Unable to cast to string for key " + key + ": " + err.Error())
	}

	return stringVal
}

func (man Manager) addConfigStringArray(key string, defVal []string, usage string) {
	man.command.PersistentFlags().StringArray(flagNameFromConfigKey(key), defVal, getFlagUsage(key, usage))
	man.viper.BindPFlag(key, man.command.PersistentFlags().Lookup(flagNameFromConfigKey(key))) //nolint:errcheck
	man.viper.BindEnv(key, envNameFromConfigKey(key))                                          //nolint:errcheck

	man.addDefault(key, defVal)
}

func (man Manager) getConfigStringArray(key string) []string {
	interfaceVal := man.getInterfaceVal(key)
	stringVal, err := cast.ToStringSliceE(interfaceVal)
	if err != nil {
		panic("Unable to cast to string array for key " + key + ": " + err.Error())
	}

	return stringVal
}

func (man Manager) addConfigInt(key string, defVal int, usage string) {
	man.command.PersistentFlags().Int(flagNameFromConfigKey(key), defVal, getFlagUsage(key, usage))
	man.viper.BindPFlag(key, man.command.PersistentFlags().Lookup(flagNameFromConfigKey(key))) //nolint:errcheck
	man.viper.BindEnv(key, envNameFromConfigKey(key))                                          //nolint:errcheck

	man.addDefault(key, defVal)
}

func (man Manager) getConfigInt(key string) int {
	interfaceVal := man.getInterfaceVal(key)
	intVal, err := cast.ToIntE(interfaceVal)
	if err != nil {
		panic("Unable to cast to int for key " + key + ": " + err.Error())
	}

	return intVal
}

func (man Manager) addConfigBool(key string, defVal bool, usage string) {
	man.command.PersistentFlags().Bool(flagNameFromConfigKey(key), defVal, getFlagUsage(key, usage))
	man.viper.BindPFlag(key, man.command.PersistentFlags().Lookup(flagNameFromConfigKey(key))) //nolint:errcheck
	man.viper.BindEnv(key, envNameFromConfigKey(key))                                          //nolint:errcheck

	man.addDefault(key, defVal)
}

func (man Manager) getConfigBool(key string) bool {
	interfaceVal := man.getInterfaceVal(key)
	boolVal, err := cast.ToBoolE(interfaceVal)
	if err != nil {
		panic("Unable to cast to string for key " + key + ": " + err.Error())
	}

	return boolVal
}

func (man Manager) getInterfaceVal(key string) interface{} {
	interfaceVal := man.viper.Get(key)
	if interfaceVal == nil {
		var ok bool
		interfaceVal, ok = man.defaults[key]
		if !ok {
			panic("Tried to look up default value for nonexistent config option: " + key)
		}
	}
	return interfaceVal
}

func (man Manager) addDefault(key string, defVal interface{}) {
	if _, exists := man.defaults[key]; exists {
		panic("Trying to add duplicate config for key " + key)
	}

	man.defaults[key] = defVal
}
