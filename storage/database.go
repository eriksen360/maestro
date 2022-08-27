package storage

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
)

func InitalizeDatabase() {

	db, err := gorm.Open(sqlite.Open("maestro.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(
		&HostModel{},
		&TaskModel{},
	)
}

func ValidateConfigFile(f string) (bool, error) {

	// TODO: Add functionality so random keys cannot be added

	hostConfigKeysMandatory := [2]string{"ipaddr", "pubkey"}
	hostConfigKeysOptional := [2]string{"fqdn", "groups"}

	fmt.Println(hostConfigKeysMandatory, hostConfigKeysOptional)

	// HostConfig
	/*
		KeysMustBeProvided: IpAddr, PubKey
		KeysCanBeProvided: Fqdn - defaults to "", Groups - defaults to []
	*/

	taskConfigKeysMandatory := [2]string{"user", "hosts"}
	taskConfigKeysOptional := [2]string{"instructions"}

	fmt.Println(taskConfigKeysMandatory, taskConfigKeysOptional)

	// TaskConfig
	/*
		KeysMustBeProvided: User, Hosts[0]
		KeysCanBeProvided: Instructions - defaults to []
	*/

	return true, nil
}

func SaveConfigFile(f string) (bool, error) {

	ok, err := ValidateConfigFile(f)
	if !ok {
		panic(err)
	}

	h, t := UnmarshalConfigFile(f)
	fmt.Println(h.Hosts)
	fmt.Println(t.Tasks)

	db, err := gorm.Open(sqlite.Open("maestro.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	for _, host := range h.Hosts {
		h := HostModel{}
		result := db.First(&h, "fqdn = ?", host.Fqdn)
		if result.RowsAffected == 0 {
			hostModel := HostModel{
				Fqdn:   host.Fqdn,
				IpAddr: host.IpAddr,
				PubKey: host.PubKey,
				Groups: strings.Join(host.Groups, ";"),
				Tasks:  nil,
			}
			db.Create(&hostModel)
			fmt.Println("hostModel created")
		} else {
			h.PubKey = host.PubKey                    // TODO: Use built in Updates methods
			h.Groups = strings.Join(host.Groups, ";") // TODO: Validate and do nothing if no fields changed
			h.IpAddr = host.IpAddr
			h.Fqdn = host.Fqdn
			db.Save(&h)
			fmt.Println("hostModel updated")
		}
	}

	for _, task := range t.Tasks {

		instructionSet := ParseInstructions(task)
		// applicableHosts := append(task.Hosts, task.Groups...)
		applicableHosts := task.Hosts
		for _, hostIpAddr := range applicableHosts {

			h := HostModel{}
			result := db.First(&h, "ip_addr = ?", hostIpAddr)
			if result.RowsAffected != 0 { // If host exists

				for _, t := range h.Tasks {

					if t.Name == task.Name {
						t.Name = task.Name
						t.InstructionSet = ""
						db.Save(&t)
						break
					} else {
						newTask := TaskModel{
							Name:           task.Name,
							InstructionSet: instructionSet,
						}
						db.Create(&newTask)
					}
				}
			} else {
				fmt.Printf("Host %v does not exist\n", h.IpAddr)
			}
		}
	}

	return true, nil
}

func GetAllTasks() []TaskModel {
	db, err := gorm.Open(sqlite.Open("maestro.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	tasks := []TaskModel{}
	db.Order("name desc").Find(&tasks)
	return tasks
}

func GetAllHosts() []HostModel {
	db, err := gorm.Open(sqlite.Open("maestro.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	hosts := []HostModel{}
	db.Order("fqdn desc").Find(&hosts)
	return hosts
}
