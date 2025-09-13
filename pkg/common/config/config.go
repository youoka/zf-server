package config

import conf "github.com/YuanJey/goconf/pkg/config"

var Config config

type config struct {
	Server struct {
		Port    string `yaml:"port" env:"ZF_PULL_PORT"`
		Cluster string `yaml:"cluster" env:"ZF_PULL_CLUSTER"`
	} `yaml:"server"`
	Mysql struct {
		Addr           string   `yaml:"addr"  env:"MYSQL_ADDR"`
		DBAddress      []string `yaml:"dbMysqlAddress"`
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
		Addr          string   `yaml:"addr"  env:"REDIS_ADDR"`
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
