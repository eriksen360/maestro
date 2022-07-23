package main

import (
	"maestro/lib"
)

func Setup() {

	ok, err := lib.InitializeViper()
	if !ok {
		panic(err)
	}
}

func main() {

	Setup()
	/* out, err := exec.Command("ssh", "-i", "/home/mathias/.ssh/bfs_thinkpad",
		"root@144.76.69.3", "ls").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out)) */

	// lib.ExecuteSingleCommand("144.76.69.3:22")
	lib.FormatConfigurationFileSettings()

	// Vi skal have en smart måde at gemme de her Task lister på de er stateful

}

/* Plan

Kan udføre kommandoer på 144.76.69.3 over ssh / anden sikker protokol
Fil-overførsel med scp/sftp
Kan parse .yml filer

Kan køre flere processer parallelt
Udvid konfigurationssprog

Kan styre hosts på lokalt netværk
Kontekst -> statefulness

Triggers
Triggers i konfigurationssprog

Kan køre processer over flere hosts med triggers og steps (event-drevet?)
Forbered performance gennem parallelisering

ssh symlink?

Byg CLI og API (REST eller gRPC) til grænseflade

Optimer og design videre med henblik på performance og høj concurrency

*/
