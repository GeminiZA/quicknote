package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

type Config struct {
	NotePath string `json:"notepath"`
}

func main() {
	curTime := time.Now()
	fmt.Println("Enter file name (blank for time):")
	var fileName string
	fmt.Scanln(&fileName)
	if fileName == "" {
		fileName = fmt.Sprintf("%s.md", curTime.Format("2006-12-30_15:04:05"))
	} else {
		i := strings.Index(fileName, ".md")
		if i == -1 {
			fileName += ".md"
		}
	}
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	var filePath string
	configFile, err := os.OpenFile(fmt.Sprintf("%s/.config/scriptconfigs/quicknote.json", user.HomeDir), os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
		filePath = fmt.Sprintf("%s/%s", user.HomeDir, fileName)
	} else {
		fileBytes, err := io.ReadAll(configFile)
		if err != nil {
			panic(err)
		}
		var config Config
		json.Unmarshal(fileBytes, &config)
		if config.NotePath == "" {
			panic(errors.New("no notepath in config"))
		}
		filePath = fmt.Sprintf("%s/%s/%s", user.HomeDir, config.NotePath, fileName)
	}
	cmd := exec.Command("nvim", filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
