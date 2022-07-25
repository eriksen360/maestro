package lib

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

type HostModel struct {
	gorm.Model
	Fqdn   string
	IpAddr string
	PubKey string
	Groups []string
	Tasks  []Task // Each host has an array of tasks corresponding to it
}

type TaskModel struct {
	gorm.Model
	TaskName       string
	User           string
	ScheduledAt    time.Time
	PersistSession bool
	LogFile        *string // Should be an io_util.file object instead
	Instructions   []InstructionModel
}

type InstructionModel struct {
	gorm.Model
	Name            string
	InstructionType int32
	Dependencies    []DependencyModel
	FileSrc         string
	FileDst         string
	Command         string
	Retries         int
}

type DependencyModel struct {
	gorm.Model
	_Host           HostModel
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
		log.Fatal("Could not migrate models", err)
		return nil
	}

	return db
}

func SaveConfigurationFile(db *gorm.DB, hosts []Host, tasks []Task) {

	for _, host := range hosts {

		var hm HostModel
		result := db.Limit(1).Find(&hm, "ipAddr = ?", host.IpAddr)
		if result.RowsAffected != 0 { // No entries found
			continue
		}

		hostModel := HostModel{
			Fqdn:   host.Fqdn,
			IpAddr: host.IpAddr,
			PubKey: host.PubKey,
			Groups: host.Groups,
		}
		db.Create(&hostModel)
	}

	// Dependencies er bundet op på instruktioner, og Instruktioner afhænger af tasks
	for _, task := range tasks {

		var instructionModels []InstructionModel
		for _, inst := range task.Instructions {

			var dependencyModels []DependencyModel
			for _, dep := range inst.Dependencies {

				var hostModel HostModel
				db.Find(&hostModel, "IpAddr = ? ", dep.HostIpAddr)

				dependencyModel := DependencyModel{
					_Host:           hostModel,
					HostIpAddr:      dep.HostIpAddr,
					TaskName:        dep.TaskName,
					InstructionName: dep.InstructionName,
				}

				dependencyModels = append(dependencyModels, dependencyModel)
				db.Create(&dependencyModel)
			}

			// Du er mongol. De hænger ikke nødvendigvis sammen jo

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
			db.Create(&instructionModel)
		}

		taskModel := TaskModel{
			TaskName:       task.TaskName,
			User:           task.User,
			ScheduledAt:    task.ScheduledAt,
			PersistSession: task.PersistSession,
			LogFile:        task.LogFile,
			Instructions:   instructionModels,
		}

		for host := range task.Hosts {
			var hm HostModel
			result := db.Limit(1).Find(&hm, "ipAddr = ?", host.IpAddr)
			if result.RowsAffected == 0 {

				// Create host?
			} else {

				if task not in hm.Tasks {
				db.Update(&hm, append(hm.Tasks, task))
				} else {
					// task already found. Error
				}

			}

		}

		// Opret task med []InstructionModel
		// For hver host i task.Hosts -> tilføj task som FK hvis task.TaskName ikke allerede
		// findes for host

	}
}
