package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 配置结构
type Config struct {
	Project struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"project"`
	Build struct {
		Output string   `json:"output"`
		Flags  []string `json:"flags"`
	} `json:"build"`
	Test struct {
		Coverage bool     `json:"coverage"`
		Verbose  bool     `json:"verbose"`
		Flags    []string `json:"flags"`
	} `json:"test"`
}

const configFile = ".goxconfig.json"

// ConfigCommand 配置管理命令
func ConfigCommand(args []string) error {
	if len(args) == 0 {
		return showConfig()
	}

	action := args[0]

	switch action {
	case "init":
		return initConfig()
	case "get":
		if len(args) < 2 {
			return fmt.Errorf("请指定配置项")
		}
		return getConfig(args[1])
	case "set":
		if len(args) < 3 {
			return fmt.Errorf("请指定配置项和值")
		}
		return setConfig(args[1], args[2])
	case "list":
		return showConfig()
	default:
		return fmt.Errorf("未知操作: %s", action)
	}
}

func initConfig() error {
	fmt.Println("⚙️  初始化配置文件...")

	// 默认配置
	config := Config{}
	config.Project.Name = "myproject"
	config.Project.Version = "1.0.0"
	config.Build.Output = "bin/"
	config.Build.Flags = []string{"-v"}
	config.Test.Coverage = true
	config.Test.Verbose = false

	return saveConfig(&config)
}

func showConfig() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	fmt.Println("⚙️  当前配置:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func getConfig(key string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// 简单的key访问（可以扩展）
	switch key {
	case "project.name":
		fmt.Println(config.Project.Name)
	case "project.version":
		fmt.Println(config.Project.Version)
	case "build.output":
		fmt.Println(config.Build.Output)
	default:
		return fmt.Errorf("未知配置项: %s", key)
	}

	return nil
}

func setConfig(key, value string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// 简单的key设置（可以扩展）
	switch key {
	case "project.name":
		config.Project.Name = value
	case "project.version":
		config.Project.Version = value
	case "build.output":
		config.Build.Output = value
	default:
		return fmt.Errorf("未知配置项: %s", key)
	}

	return saveConfig(config)
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("配置文件不存在，请运行 'gox config init'")
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("配置文件格式错误: %w", err)
	}

	return &config, nil
}

func saveConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return err
	}

	absPath, _ := filepath.Abs(configFile)
	fmt.Printf("✅ 配置已保存: %s\n", absPath)
	return nil
}
