package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	dockerConfigPath := getEnv("DOCKER_CONFIG_PATH", "/etc/docker/daemon.json")
	dockerDebugValue := getEnv("DOCKER_DEBUG_VALUE", "true")

	if err := updateJSON(dockerConfigPath, dockerDebugValue); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := sighupDocker(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}

func getEnv(key, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func updateJSON(configPath, debug string) error {
	dockerConfig := make(map[string]interface{})
	var configFileMissing bool

	dockerConfigFromFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		configFileMissing = true
	}

	if !configFileMissing {
		err = json.Unmarshal(dockerConfigFromFile, &dockerConfig)
		if err != nil {
			return err
		}
	}

	var debugBool bool
	switch strings.ToLower(debug) {
	case "true":
		debugBool = true
	case "false":
		debugBool = false
	}
	dockerConfig["debug"] = debugBool

	dockerConfigToFile, err := json.Marshal(dockerConfig)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, dockerConfigToFile, 0644)
	if err != nil {
		return err
	}

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
