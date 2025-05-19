package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/term"
)

var configs = mysql.Config{
	User:                 "root",
	Passwd:               "jeno1234",
	Net:                  "tcp",
	Addr:                 fmt.Sprintf("%s:%s", "127.0.0.1", "3306"),
	DBName:               "terminal_task",
	AllowNativePasswords: true,
	ParseTime:            true,
}

type DBConfig struct {
	User                 string `json:"user"`
	Passwd               string `json:"passwd"`
	Net                  string `json:"net"`
	Addr                 string `json:"addr"`
	DBName               string `json:"dbName"`
	AllowNativePasswords bool   `json:"allowNativePasswords"`
	ParseTime            bool   `json:"parseTime"`
}

func prompt(reader *bufio.Reader, label, defaultVal string, promptType string) (string, error) {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	var input string = ""
	var err error = nil
	if promptType == "password" {
		byteInput, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		fmt.Println()
		input = string(byteInput)
	} else {
		input, err = reader.ReadString('\n')
		if err != nil {
			return "", err
		}
	}

	input = strings.TrimSpace(input)
	if input == "" && (promptType == "password" || promptType == "userName") {
		fmt.Println("Required field! Please enter a value.")
		return prompt(reader, label, defaultVal, promptType)
	}
	return input, nil
}

func loadConfig() (mysql.Config, error) {
	// Determine config path
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return mysql.Config{}, err
	}
	cfgDir = filepath.Join(cfgDir, "terminal_task")
	cfgFile := filepath.Join(cfgDir, "config.json")

	// First run: prompt & write
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		if err := os.MkdirAll(cfgDir, 0o755); err != nil {
			return mysql.Config{}, err
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Welcome to terminal_task! Let's do one time set up your database connection:")

		user, err := prompt(reader, "DB user", "root", "userName")
		if err != nil {
			return mysql.Config{}, err
		}
		passwd, err := prompt(reader, "DB password", "", "password")
		if err != nil {
			return mysql.Config{}, err
		}
		netProto, err := prompt(reader, "Network protocol", "default : tcp", "networdProtocol")
		if err != nil {
			return mysql.Config{}, err
		}
		if netProto == "" {
			netProto = "tcp"
		}
		addr, err := prompt(reader, "DB address (host:port)", "default: 127.0.0.1:3306", "address")
		if err != nil {
			return mysql.Config{}, err
		}
		if addr == "" {
			addr = "127.0.0.1:3306"
		}
		dbName := "terminal_task"

		cfg := DBConfig{
			User:                 user,
			Passwd:               passwd,
			Net:                  netProto,
			Addr:                 addr,
			DBName:               dbName,
			AllowNativePasswords: true,
			ParseTime:            true,
		}

		f, err := os.Create(cfgFile)
		if err != nil {
			return mysql.Config{}, err
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode(cfg); err != nil {
			return mysql.Config{}, err
		}

		fmt.Printf("Configuration saved to %s\n\n", cfgFile)
		// proceed with the just-entered values rather than exiting
		return mysql.Config{
			User:                 cfg.User,
			Passwd:               cfg.Passwd,
			Net:                  cfg.Net,
			Addr:                 cfg.Addr,
			DBName:               cfg.DBName,
			AllowNativePasswords: cfg.AllowNativePasswords,
			ParseTime:            cfg.ParseTime,
		}, nil
	}

	// Subsequent runs: load from file
	f, err := os.Open(cfgFile)
	if err != nil {
		return mysql.Config{}, err
	}
	defer f.Close()

	var c DBConfig
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return mysql.Config{}, err
	}

	return mysql.Config{
		User:                 c.User,
		Passwd:               c.Passwd,
		Net:                  c.Net,
		Addr:                 c.Addr,
		DBName:               c.DBName,
		AllowNativePasswords: c.AllowNativePasswords,
		ParseTime:            c.ParseTime,
	}, nil
}
