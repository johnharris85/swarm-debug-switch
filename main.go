package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"fmt"
	"strconv"
	"syscall"
	"strings"
)

func main() {
	dockerConfigPath := getEnv("DOCKER_CONFIG_PATH", "/etc/docker/daemon.json")
	dockerDebugValue := getEnv("DOCKER_DEBUG_VALUE", "true")

	if err := updateJson(dockerConfigPath, dockerDebugValue); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := sighupDocker(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func getEnv(key, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func updateJson(configPath, debug string) error {
	fileReader, err := os.Open(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		dockerConfig := make(map[string]interface{})
		dockerConfig["debug"] = true

		fileWriter, err := os.Create(configPath)
		if err != nil {
			return err
		}
		defer fileWriter.Close()
		json.NewEncoder(fileWriter).Encode(dockerConfig)

		return nil
	}

	var debugBool bool
	switch strings.ToLower(debug) {
	case "true":
		debugBool = true
	case "false":
		debugBool = false
	}

	var dockerConfig map[string]interface{}
	json.NewDecoder(fileReader).Decode(&dockerConfig)

	dockerConfig["debug"] = debugBool

	fileWriter, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer fileWriter.Close()
	json.NewEncoder(fileWriter).Encode(dockerConfig)

	return nil
}

func sighupDocker() error {
	dockerPid, err := exec.Command("pidof", "dockerd").Output()
	if err != nil {
		return err
	}

	dockerPidInt, err := strconv.Atoi(strings.Trim(string(dockerPid), "\n"))
	if err != nil {
		return err
	}

	dockerProcess, err := os.FindProcess(dockerPidInt)
	if err != nil {
		return err
	}

	err = dockerProcess.Signal(syscall.SIGHUP)
	if err != nil {
		return err
	}

	return nil
}