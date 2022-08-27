package main

import (
	"maestro/cmd"
	"maestro/storage"
)

// Main program.
func main() {
	// Print the config options from the new conf struct instance.
	storage.InitalizeDatabase()
	cmd.Execute()
}
