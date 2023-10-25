package main

// this is just for testing.
import "log"

func LogInfo(message string) {
	log.Printf("INFO - %v", message)
}

func LogWarning(message string) {
	log.Printf("WARN - %v", message)
}

func LogError(message string) {
	log.Printf("ERROR - %v", message)
}

func main() {
	log.Print("Running main.go....")
	log.Print("Finished!")
}
