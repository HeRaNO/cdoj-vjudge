package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Configure struct {
	RDB    *RDBConfig   `yaml:"RDB"`
	Redis  *RedisConfig `yaml:"Redis"`
	MQ     *MQConfig    `yaml:"MQ"`
	Server *SrvConfig   `yaml:"Server"`
}

var conf Configure

type RDBConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	DBName   string `yaml:"DBName"`
}

type RedisConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

type MQConfig struct {
	Host      string `yaml:"Host"`
	Port      int    `yaml:"Port"`
	UserName  string `yaml:"Username"`
	Password  string `yaml:"Password"`
	QueueName string `yaml:"QueueName"`
}

type SrvConfig struct {
	Port int `yaml:"Port"`
}

var SrvAddr string

func InitConfig(file *string) {
	fileBytes, err := os.ReadFile(*file)
	if err != nil {
		log.Println("[FAILED] read config file failed")
		panic(err)
	}
	if err = yaml.Unmarshal(fileBytes, &conf); err != nil {
		log.Println("[FAILED] unmarshal yaml file failed")
		panic(err)
	}
}

func Init(file *string) {
	InitConfig(file)

	SrvAddr = fmt.Sprintf(":%d", conf.Server.Port)

	InitDB()
	InitRedis()
	InitMQ()

	log.Println("[INFO] init finished successfully")
}
