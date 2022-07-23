package lib

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
)

/* Flow

New config file passed to CLI
Read, validate and parse file
Convert to instruction set
Save instruction set and hosts in config file
"Discard" configuration file

*/

type InstructionType int32

const (
	FileTransfer int32 = 0
	Command            = 1
)

func InitializeViper() (bool, error) {

	pwd, _ := os.Getwd()
	var base string = pwd + "/base.yaml"
	_, err := os.Create(base)
	if err != nil {
		fmt.Errorf("failed creating initial configuraiton file %v with error %w", base, err)
		return false, err
	}

	viper.SetConfigName("base")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(pwd + "/conf")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("failed reading in config file: %w", err)
		return false, err
	}

	return true, nil
}

// TODO: test
func ParseConfigurationFile(fileName string) (bool, error) {

	s := strings.Split(fileName, ".")
	viper.SetConfigName(s[0])
	err := viper.MergeInConfig()
	if err != nil {
		fmt.Errorf("failed in adding configuration file to base %w", err)
		return false, err
	}

	return true, nil
}

func AddConfigurationFileSettings(fileName string) (bool, error) {

	// Would in principle require a CLI interface to pass config file from anywhere, but for now we will
	// only use the conf/ directory
	ParseConfigurationFile(fileName)

	// Validate
	// Convert to instruction set

	// Add to file system for use by executor. NB only instructions sets are saved. Config files are
	// evaluated only once when passed to CLI

	return true, nil
}

func Validate(settings map[string]interface{}) (bool, error) {

	// Validate file

	return true, nil
}

// FormatConfigurationFileSettings This function translates the map of settings into a
// more explicit list, which can be further formatted to provide a concise instruction set
// for the executor algorithm
func FormatConfigurationFileSettings() {

	var tasks []Task

	// Should only validate single file passed!
	ok, err := Validate(viper.AllSettings())
	if !ok {
		panic(fmt.Errorf("invalid configuration file: %w", err))
	}

	// Do the same for hosts (logic already exists, just needs refactoring)

	// if viper.IsSet("tasks")
	for key, _ := range viper.GetStringMapString("tasks") {

		var keyPath string = "tasks." + key
		var logFile *string = nil

		user := viper.GetString(keyPath + ".user")
		hosts := viper.GetStringSlice(keyPath + ".hosts")
		groups := viper.GetStringSlice(keyPath + ".groups")

		var resolvedGroups []string = ResolveGroupToHosts(groups)
		hosts = append(hosts, resolvedGroups...) // Must unpack slice to append

		schedule := viper.GetTime(keyPath + ".schedule")
		persistSession := viper.GetBool(keyPath + ".persistSession")

		// Iterate over taskInfo to get user, hosts, resolve groups, get schedule and persistSession

		if viper.GetBool(keyPath + ".logs.logging") {
			logFilePath := viper.GetString(keyPath + ".logs.logging.logFile")
			logFile = &logFilePath
		}

		newTask := Task{
			user:           user,
			hosts:          hosts,
			scheduledAt:    schedule,
			persistSession: persistSession,
			logFile:        logFile,
		}

		newTask.instructions = CreateInstructionList(keyPath + ".instructions")
		tasks = append(tasks, newTask)
	}
}

func GetHostSettings() []Host {

	var hosts []Host

	AddConfigurationFileSettings("hosts")
	for key, _ := range viper.GetStringMapString("hosts") {

		keyPath := "hosts." + key

		host := Host{
			fqdn:   viper.GetString(keyPath + ".fqdn"),
			ipAddr: viper.GetString(keyPath + ".ipAddr"),
			pubKey: viper.GetString(keyPath + ".pubkey"),
			groups: viper.GetStringSlice(keyPath + ".groups"),
		}

		hosts = append(hosts, host)
	}

	return hosts
}

func ResolveGroupToHosts(groups []string) []string {

	var hostAddr []string
	var hosts []Host = GetHostSettings()

	for _, host := range hosts {
		for _, group := range groups {
			for _, hostGroup := range host.groups {
				if group == hostGroup {
					hostAddr = append(hostAddr, host.ipAddr)
					break
				}
			}
		}
	}

	return hostAddr
}

func CreateInstructionList(path string) []Instruction {

	var i int = 0
	var instructions []Instruction

	for viper.Get(path+"."+strconv.Itoa(i)) != nil {

		var j int = 0
		var fileSrc string
		var fileDst string
		var command string
		var instructionTypeInt int32
		var dependencies []Dependency

		var extPath string = path + "." + strconv.Itoa(i)
		name := viper.GetString(extPath + ".name")
		instructionType := viper.GetString(extPath + ".type")

		switch instructionType {
		case "fileTransfer":
			instructionTypeInt = FileTransfer
			fileSrc = "123"
			fileSrc = viper.GetString(extPath + ".filesrc")
			fileDst = viper.GetString(extPath + ".filedst")
			break
		default:
			instructionTypeInt = Command
			command = viper.GetString(extPath + ".command")
		}

		for viper.Get(extPath+".dependencies."+strconv.Itoa(j)) != nil {

			var depPath string = extPath + ".dependencies." + strconv.Itoa(j)

			host := viper.GetString(depPath + ".host")
			task := viper.GetString(depPath + ".task")
			nameOfStep := viper.GetString(depPath + ".stepname") // is empty

			dependency := Dependency{
				host:       host,
				task:       task,
				nameOfStep: nameOfStep,
			}

			dependencies = append(dependencies, dependency)
			j++
		}

		retries := viper.GetInt(extPath + ".retries")

		/* for _, key := range keys {
			if viper.IsSet(extPath + key) {
				fmt.Println(viper.Get(extPath + key))
			}
		} */

		// Use keys to build instruction set

		instruction := Instruction{
			id:              uuid.New(),
			name:            name,
			instructionType: instructionTypeInt,
			dependencies:    dependencies,
			fileSrc:         fileSrc,
			fileDst:         fileDst,
			command:         command,
			retries:         retries,
		}

		instructions = append(instructions, instruction)
		i++
	}

	//in := instructions[0]
	//fmt.Printf("id:%v\ntype:%v\ndep:%v\nsrc:%v\ndst:%v\ncmd:%v\nretries:%v\n\n\n", in.id, in.instructionType, in.dependencies, in.fileSrc, in.fileDst, in.command, in.retries)

	/* Uses the formatted conf file to create a list of instructions specific to a host, and passes
	   them to a context where the orchestrator can execute them */

	/* 	The steps are formatted in an ordered list of exeutions



	[
	1: {
		type: enum(1, 2, 3)     The type specifies which action path should be taken. Each
								action path expects some specific arguments. 'fileTransfer'
								expects fileSrc and fileDst while 'comand' expects command

		dependencies: {         This is a list of events that must be satisfied before the step
								can be executed.

		command: "abc | grep xyz"
		retries: N   			Number of times command is tried if return code is not zero
		desiredState: ?

	},
	2: {...},
	]
	*/

	return instructions
}

type Dependency struct {
	// instructionUUID int  Should be used in the future
	host       string
	task       string
	nameOfStep string
}

type Instruction struct {
	id              uuid.UUID
	name            string
	instructionType int32
	dependencies    []Dependency
	fileSrc         string
	fileDst         string
	command         string
	retries         int
}

type Task struct {
	user           string
	hosts          []string
	scheduledAt    time.Time
	persistSession bool
	logFile        *string // Should be a ioutil.file object instead
	instructions   []Instruction
}

type Host struct {
	fqdn   string
	ipAddr string
	pubKey string // change to x/ssh type for security
	groups []string
	// ...
}
