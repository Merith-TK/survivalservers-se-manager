package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"
	// toml "github.com/BurntSushi/toml"
)

var (
	conf       config
	configfile = "config.json"
)

type config struct {
	Server struct {
		Address  string `json:"address"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"server"`
	Files []struct {
		Local  string `json:"local"`
		Remote string `json:"remote"`
	} `json:"files"`
}

func setDefaultConfig() {
	conf.Server.Address = "127.0.0.1"
	conf.Server.Port = "21"
	conf.Server.User = "username"
	conf.Server.Password = "password"
	conf.Files = append(conf.Files, struct {
		Local  string `json:"local"`
		Remote string `json:"remote"`
	}{Local: "local", Remote: "remote"})
}

func setupConfig() error {
	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		log.Println("[Config] No config file found, creating default config")
		f, err := os.Create(configfile)
		if err != nil {
			return err
		}
		setDefaultConfig()

		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		enc.Encode(conf)

		f.Close()
	} else {
		log.Println("[Config] Found config, loading config")
		f, err := os.Open(configfile)
		if err != nil {
			return err
		}
		dec := json.NewDecoder(f)
		err = dec.Decode(&conf)
		if err != nil {
			return err
		}
	}

	setupEnvironment()
	return nil
}

func setupEnvironment() {
	log.Println("[Config] Loading Environment Variables")
	// Replace Environment Variables
	for k, v := range conf.Files {
		// Define the regular expression pattern
		pattern := `\{([^}]+)\}`

		// // Compile the regular expression
		// regExp, err := regexp.Compile(pattern)
		// if err != nil {
		// 	fmt.Println("Error compiling regular expression:", err)
		// 	return
		// }

		// matches := regExp.FindAllStringSubmatch(v, -1)
		// for _, match := range matches {
		// 	// placeholder := match[1]
		// 	// newValue := os.Getenv(placeholder)
		// 	// v = strings.Replace(v, match[0], newValue, -1)
		// }
		// conf.Files[k] = v

		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(v.Local, -1)
		for _, match := range matches {
			placeholder := match[1]
			newValue := os.Getenv(placeholder)
			v.Local = strings.Replace(v.Local, match[0], newValue, -1)
		}
		conf.Files[k] = v
	}
}
