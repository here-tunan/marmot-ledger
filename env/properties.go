package env

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Properties struct {
	Port string `json:"port"`

	Mysql struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"mysql"`

	Redis struct {
		HOST     string `json:"host"`
		Port     string `json:"port"`
		Db       int    `json:"db"`
		Password string `json:"password"`
	} `json:"redis"`
}

var Prop *Properties

// 初始化配置文件
func init() {
	readEncProperties()
}

func readEncProperties() {

	envFilePath := "./env/dev.yaml"

	// 根据环境变量解析配置文件
	appEnv := os.Getenv("MARMOT_LEDGER_ENV")
	if appEnv == "dev" || appEnv == "" {
		fmt.Println("Start marmot-ledger app in development environment!")
		envFilePath = "./env/dev.yaml"
	} else if appEnv == "prod" {
		fmt.Println("Start marmot-ledger app in production environment!")
		envFilePath = "./env/prod.yaml"
	} else {
		fmt.Println("Start marmot-ledger app in unknown environment! Maybe cause errors!")
	}

	// 读取YAML文件内容
	envFile, err := os.ReadFile(envFilePath)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	prop := &Properties{}

	// 解析YAML文件
	err = yaml.Unmarshal(envFile, prop)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}
	// 打印解析后的数据
	// fmt.Printf("%+v\n", prop)

	Prop = prop
}
