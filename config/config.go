package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"

	yaml "gopkg.in/yaml.v2"
)

const (

	// file
	envConfigHomePath      = "ENV_CONFIG_HOME_PATH"
	ReadWritePerms         = os.FileMode(0755)
	configFolderName       = ".cli"
	configFileName         = "config.yaml"
	contextsFolderName     = "contexts"
	defaultContextFileName = "default.yaml"

	// key
	CurrentContext = "current-context"
	CliVersion     = "cli-version"
	ApiURL         = "api-URL"

	// default value
	DefaultContext = "default"
	DefaultAPIURL  = "http://localhost:8080"
)

// path
var configHomePath string
var configFilePath string
var contextsFolderPath string
var defaultContextFilePath string
var contextFilePath string

// config/context
type Config map[string]string

var config = &Config{
	CurrentContext: DefaultContext,
	CliVersion:     Version,
}
var context = &Config{
	ApiURL: DefaultAPIURL,
}

// Init is called by main.go, no need to lock
func Init() {
	var err error

	// check env
	configHomePath = os.Getenv(envConfigHomePath)
	if configHomePath == "" {
		userHomePath, err := homedir.Dir()
		if err != nil {
			fmt.Printf("error finding home %s\n", err)
			os.Exit(1)
		}
		configHomePath = filepath.Join(userHomePath, configFolderName)
	}

	// ~/.cli
	if _, err := os.Stat(configHomePath); os.IsNotExist(err) {
		if err = os.Mkdir(configHomePath, ReadWritePerms); err != nil {
			fmt.Printf("error creating %s\n", configHomePath)
			os.Exit(1)
		}
		fmt.Printf("config home=%s\n", configHomePath)
	}

	// ~/.cli/config.yaml
	configFilePath = filepath.Join(configHomePath, configFileName)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if _, err := os.Create(configFilePath); err != nil {
			fmt.Printf("error creating %s\n", configFilePath)
			os.Exit(1)
		}

		// write default
		if err = WriteYamlFile(configFilePath, config); err != nil {
			fmt.Printf("error writing %s\n", configFilePath)
			os.Exit(1)
		}
		fmt.Printf("config file=%s\n", configFilePath)
	}

	// read latest config
	config, err = ReadYamlFile(configFilePath)
	if err != nil {
		fmt.Printf("error reading %s\n", configFilePath)
		os.Exit(1)
	}

	// ~/.cli/contexts
	contextsFolderPath = filepath.Join(configHomePath, contextsFolderName)
	if _, err := os.Stat(contextsFolderPath); os.IsNotExist(err) {
		if err = os.Mkdir(contextsFolderPath, ReadWritePerms); err != nil {
			fmt.Printf("error creating %s\n", contextsFolderPath)
			os.Exit(1)
		}
	}

	// ~/.cli/contexts/default.yaml
	defaultContextFilePath = filepath.Join(contextsFolderPath, defaultContextFileName)
	if _, err := os.Stat(defaultContextFilePath); os.IsNotExist(err) {
		if _, err = os.Create(defaultContextFilePath); err != nil {
			fmt.Printf("error creating %s\n", defaultContextFilePath)
			os.Exit(1)
		}

		// write default
		if err = WriteYamlFile(defaultContextFilePath, context); err != nil {
			fmt.Printf("error writing %s\n", defaultContextFilePath)
			os.Exit(1)
		}
		fmt.Printf("default context file=%s\n", defaultContextFilePath)
	}

	// read context in use
	contextFilePath = filepath.Join(contextsFolderPath, (*config)[CurrentContext]+".yaml")
	if _, err := os.Stat(contextFilePath); os.IsNotExist(err) {
		fmt.Printf("cannot find context file=%s\n", contextFilePath)
		os.Exit(1)
	}
	context, err = ReadYamlFile(contextFilePath)
	if err != nil {
		fmt.Printf("error reading %s\n", contextFilePath)
		os.Exit(1)
	}
}

func PrintConfig() {
	fmt.Printf("\nconfig home=%s\n", configHomePath)
	fmt.Printf("current context=%s\n\n", Get(CurrentContext))
}

func Get(key string) string {
	return (*config)[key]
}

func GetFromContext(key string) string {
	return (*context)[key]
}

func UseContext(context string) {
	var err error

	// check
	contextFilePath = filepath.Join(contextsFolderPath, context+".yaml")
	if _, err := os.Stat(contextFilePath); os.IsNotExist(err) {
		fmt.Printf("cannot find context file=%s\n", contextFilePath)
		os.Exit(1)
	}
	// write
	(*config)[CurrentContext] = context
	if err = WriteYamlFile(configFilePath, config); err != nil {
		fmt.Printf("error writing %s\n", configFilePath)
		os.Exit(1)
	}
	PrintConfig()

	// reload is not necessary for cli
	// Init()

}

func ReadYamlFile(filename string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(bytes, config)
	return config, err
}

func WriteYamlFile(filename string, config *Config) error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, ReadWritePerms)
}
