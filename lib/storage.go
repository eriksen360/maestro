package lib

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type HostModel struct {
	gorm.Model
	Fqdn   string
	IpAddr string
	PubKey string
	Groups string
	Tasks  []TaskModel `gorm:"foreignKey:TaskId"`
}

type TaskModel struct {
	gorm.Model
	TaskId         uint `gorm:"primarykey"`
	TaskName       string
	User           string
	ScheduledAt    time.Time
	PersistSession bool
	LogFile        *string // Should be an io_util.file object instead
	InstructionID  int
	Instructions   []InstructionModel `gorm:"foreignKey:InstructionId"`
}

type InstructionModel struct {
	gorm.Model
	InstructionId   uint `gorm:"primarykey"`
	Name            string
	InstructionType int32
	Dependencies    []DependencyModel `gorm:"foreignKey:DependencyId"`
	FileSrc         string
	FileDst         string
	Command         string
	Retries         int
}

type DependencyModel struct {
	gorm.Model
	DependencyId    uint `gorm:"primarykey"`
	HostIpAddr      string
	TaskName        string
	InstructionName string
}

func Connect() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&HostModel{},
		&TaskModel{},
		&InstructionModel{},
		&DependencyModel{},
	)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

func SaveConfigurationFile(db *gorm.DB, hosts []Host, tasks []Task) {

	var TaskExists bool = false
	for _, host := range hosts {

		var hm HostModel
		result := db.Find(&hm, "ip_addr = ?", host.IpAddr)
		if result.RowsAffected != 0 { // If existing host already configured
			continue
		}

		hostModel := HostModel{
			Fqdn:   host.Fqdn,
			IpAddr: host.IpAddr,
			PubKey: host.PubKey,
			Groups: strings.Join(host.Groups, ","),
		}
		db.Create(&hostModel)
	}

	for _, task := range tasks {

		fmt.Println("\nTask: ", task.TaskName)
		var instructionModels []InstructionModel
		for _, inst := range task.Instructions {

			fmt.Println("Instruction name: ", inst.Name)
			var dependencyModels []DependencyModel
			for i, dep := range inst.Dependencies {

				fmt.Println("Dependency number: ", i)
				dependencyModel := DependencyModel{
					HostIpAddr:      dep.HostIpAddr,
					TaskName:        dep.TaskName,
					InstructionName: dep.InstructionName,
				}

				dependencyModels = append(dependencyModels, dependencyModel)
			}

			// Dependencies for a specific instruction are bound to that instruction
			instructionModel := InstructionModel{
				Name:            inst.Name,
				InstructionType: inst.InstructionType,
				Dependencies:    dependencyModels,
				FileSrc:         inst.FileSrc,
				FileDst:         inst.FileDst,
				Command:         inst.Command,
				Retries:         inst.Retries,
			}
			instructionModels = append(instructionModels, instructionModel)
		}

		taskModel := TaskModel{
			TaskName:       task.TaskName,
			User:           task.User,
			ScheduledAt:    task.ScheduledAt,
			PersistSession: task.PersistSession,
			LogFile:        task.LogFile,
			Instructions:   instructionModels,
		}

		// Couple Task and Host to create unique pair
		for _, hostIpAddr := range task.Hosts {

			var hm HostModel
			result := db.Find(&hm, "ip_addr = ?", hostIpAddr)
			if result.RowsAffected != 0 {

				for _, existingTask := range hm.Tasks {
					if existingTask.TaskName == task.TaskName {
						TaskExists = true
					}
				}
				if TaskExists {
					fmt.Printf("Task %v already exists for host %v",
						task.TaskName, hostIpAddr)
				} else {
					hm.Tasks = append(hm.Tasks, taskModel)
					db.Save(&hm)
				}
			} else {
				fmt.Printf("Configuration for host %v not found", hostIpAddr)
			}
		}
	}
}
