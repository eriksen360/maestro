package storage

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func UnmarshalConfigFile(f string) (HostConfig, TaskConfig) {

	path := strings.Split(f, "/")
	file := strings.Split(path[len(path)-1], ".")
	path = path[:len(path)-1]

	sViper := viper.New()
	sViper.AddConfigPath(strings.Join(path, "/"))
	sViper.SetConfigName(file[0])
	sViper.SetConfigType(file[1])
	sViper.ReadInConfig()

	h := UnmarshalHostsFromConfigFile(sViper)
	t := UnmarshalTasksFromConfigFile(sViper)

	return h, t
}

func UnmarshalHostsFromConfigFile(sViper *viper.Viper) HostConfig {

	var hosts []Host
	hostsFromFile := sViper.GetStringMapStringSlice("hosts")
	for h, _ := range hostsFromFile {
		var host Host
		sViper.UnmarshalKey("hosts."+h, &host)
		hosts = append(hosts, host)
	}

	hostsConfig := HostConfig{
		Hosts: hosts,
	}
	return hostsConfig
}

func UnmarshalTasksFromConfigFile(sViper *viper.Viper) TaskConfig {

	var tasks []Task
	tasksFromFile := sViper.GetStringMapStringSlice("tasks")
	for t, _ := range tasksFromFile {
		var task Task
		sViper.UnmarshalKey("tasks."+t, &task)
		task.Name = t
		tasks = append(tasks, task)
	}
	taskConfig := TaskConfig{
		Tasks: tasks,
	}
	return taskConfig
}

func ParseInstructions(task Task) string {

	var instructionStr string
	for _, instruction := range task.Instructions {
		instructionStr += instruction.Command
	}
	fmt.Printf("Instruction %v for task %v", instructionStr, task.Name)
	return instructionStr
}
