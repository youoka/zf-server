package config

import conf "github.com/YuanJey/goconf/pkg/config"

var Config config

type config struct {
	Server struct {
		Port string `yaml:"port" env:"ZF_PULL_PORT"`
	} `yaml:"server"`
	Mysql struct {
		DBAddress      []string `yaml:"dbMysqlAddress" env:"MYSQL_ADDRESS"`
		DBUserName     string   `yaml:"dbMysqlUserName" env:"MYSQL_USERNAME"`
		DBPassword     string   `yaml:"dbMysqlPassword" env:"MYSQL_PASSWORD"`
		DBDatabaseName string   `yaml:"dbMysqlDatabaseName"`
		DBMaxOpenConns int      `yaml:"dbMaxOpenConns"`
		DBMaxIdleConns int      `yaml:"dbMaxIdleConns"`
		DBMaxLifeTime  int      `yaml:"dbMaxLifeTime"`
		LogLevel       int      `yaml:"logLevel"`
		SlowThreshold  int      `yaml:"slowThreshold"`
	} `yaml:"mysql"`
	Redis struct {
		DBAddress     []string `yaml:"dbAddress"`
		DBMaxIdle     int      `yaml:"dbMaxIdle"`
		DBMaxActive   int      `yaml:"dbMaxActive"`
		DBIdleTimeout int      `yaml:"dbIdleTimeout"`
		DBUserName    string   `yaml:"dbUserName"`
		DBPassWord    string   `yaml:"dbPassWord"`
		EnableCluster bool     `yaml:"enableCluster"`
	} `yaml:"redis"`
}

func init() {
	conf.UnmarshalConfig(&Config, "config.yaml")
}
