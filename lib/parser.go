package lib

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	FileTransfer int32 = 0
	Command            = 1
)

func InitializeViper(file string) *viper.Viper {

	pwd, _ := os.Getwd()
	var sViper = viper.New()
	s := strings.Split(file, ".")

	sViper.SetConfigName(s[0])
	sViper.SetConfigType(s[1])
	sViper.AddConfigPath(pwd + "/conf")

	err := sViper.ReadInConfig()
	if err != nil {
		log.Fatal("failed reading in config file: %w", err)
	}

	return sViper
}

func ParseConfigurationFile(file string) (bool, error) {

	var sViper *viper.Viper = InitializeViper(file)
	ok, err := Validate(sViper)
	if !ok {
		log.Fatal(err)
	}

	var hosts []Host = FormatHostSettingsFromFile(sViper)
	var tasks []Task = FormatTaskSettingsFromFile(sViper)

	db := Connect()
	if db == nil {
		panic("Could not get database")
	}
	SaveConfigurationFile(db, hosts, tasks)

	return true, nil
}

func FormatHostSettingsFromFile(sViper *viper.Viper) []Host {

	var hosts []Host

	if !sViper.IsSet("hosts") {
		log.Println("hosts not found")
		return nil
	}

	for h, _ := range sViper.GetStringMapString("hosts") {
		fqdn := sViper.GetString("hosts." + h + ".fqdn")
		ipAddr := sViper.GetString("hosts." + h + ".ipAddr")
		pubKey := sViper.GetString("hosts." + h + ".pubKey")
		groups := sViper.GetStringSlice("hosts." + h + ".groups")

		host := Host{
			Fqdn:   fqdn,
			IpAddr: ipAddr,
			PubKey: pubKey,
			Groups: groups,
		}

		hosts = append(hosts, host)
	}

	return hosts
}
func FormatTaskSettingsFromFile(sViper *viper.Viper) []Task {

	var tasks []Task

	if !sViper.IsSet("tasks") {
		log.Println("tasks not found")
		return nil
	}

	for t, _ := range sViper.GetStringMapString("tasks") {

		var logFile *string = nil

		user := sViper.GetString("tasks." + t + ".user")
		hosts := sViper.GetStringSlice("tasks." + t + ".hosts")
		groups := sViper.GetStringSlice("tasks." + t + ".groups")
		schedule := sViper.GetTime("tasks." + t + ".schedule")
		persistSession := sViper.GetBool("tasks." + t + ".persistSession")

		var resolvedGroups []string = ResolveGroupToHosts(sViper, groups)
		hosts = append(hosts, resolvedGroups...) // Must unpack slice to append

		if sViper.GetBool("tasks." + t + ".logs.logging") {
			logFilePath := sViper.GetString("tasks." + t + ".logs.logging.logFile")
			logFile = &logFilePath
		}

		newTask := Task{
			TaskName:       t,
			User:           user,
			Hosts:          hosts,
			ScheduledAt:    schedule,
			PersistSession: persistSession,
			LogFile:        logFile,
		}

		newTask.Instructions = CreateInstructionList(sViper, &newTask, t)
		tasks = append(tasks, newTask)
	}

	return tasks
}

func Validate(sViper *viper.Viper) (bool, error) {

	// Validate file

	return true, nil
}

func ResolveGroupToHosts(sViper *viper.Viper, groups []string) []string {

	var hostAddr []string
	var hosts []Host = FormatHostSettingsFromFile(sViper)

	for _, host := range hosts {
		for _, group := range groups {
			for _, hostGroup := range host.Groups {
				if group == hostGroup {
					hostAddr = append(hostAddr, host.IpAddr)
					break
				}
			}
		}
	}

	return hostAddr
}

func CreateInstructionList(sViper *viper.Viper, taskObj *Task, t string) []Instruction {

	var i int = 0
	var instPath string = "tasks." + t + ".instructions."
	var instructions []Instruction

	for sViper.Get(instPath+strconv.Itoa(i)) != nil {

		var j int = 0
		var fileSrc string
		var fileDst string
		var command string
		var instStepPath string = instPath + strconv.Itoa(i)
		var instructionType int32
		var dependencies []Dependency

		name := sViper.GetString(instStepPath + ".name")
		instType := sViper.GetString(instStepPath + ".type")
		retries := sViper.GetInt(instStepPath + ".retries")

		switch instType {
		case "fileTransfer":
			instructionType = FileTransfer
			fileSrc = sViper.GetString(instStepPath + ".file_src")
			fileDst = sViper.GetString(instStepPath + ".file_dst")
			break
		case "command":
			instructionType = Command
			command = sViper.GetString(instStepPath + ".command")
		default:
			break
		}

		var depPath string = instStepPath + ".dependencies."
		for sViper.Get(depPath+strconv.Itoa(j)) != nil {

			host := sViper.GetString(depPath + strconv.Itoa(j) + ".host")
			taskName := sViper.GetString(depPath + strconv.Itoa(j) + ".task")
			instructionName := sViper.GetString(depPath + strconv.Itoa(j) + ".step_name")

			dependency := Dependency{
				HostIpAddr:      host,
				TaskName:        taskName,
				InstructionName: instructionName,
			}

			dependencies = append(dependencies, dependency)
			j++
		}

		instruction := Instruction{
			Name:            name,
			InstructionType: instructionType,
			Dependencies:    dependencies,
			FileSrc:         fileSrc,
			FileDst:         fileDst,
			Command:         command,
			Retries:         retries,
		}

		instructions = append(instructions, instruction)
		i++
	}

	return instructions
}

type Host struct {
	Fqdn   string
	IpAddr string
	PubKey string
	Groups []string
	Tasks  []Task // Each host has an array of tasks corresponding to it
}

type Task struct {
	TaskName       string
	User           string
	Hosts          []string
	ScheduledAt    time.Time
	PersistSession bool
	LogFile        *string // Should be an io_util.file object instead
	Instructions   []Instruction
}

type Instruction struct {
	Name            string
	InstructionType int32
	Dependencies    []Dependency
	FileSrc         string
	FileDst         string
	Command         string
	Retries         int
}

type Dependency struct {
	HostIpAddr      string
	TaskName        string
	InstructionName string
}
