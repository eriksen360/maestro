package main

import "maestro/lib"

func main() {
	lib.Parse()
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