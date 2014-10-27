package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		log.Println("Proper usage:")
		log.Println("ping-run command-to-execute address-to-ping [additional-addresses-to-ping...]")
		log.Println("Failure, exiting")
		return
	}
	command := os.Args[1]
	addresses := os.Args[2:]
	remaining := len(addresses)
	doneArr := make([]bool, remaining)
	done := make(chan int)
	ticker := time.Tick(time.Second * 15)

	for id, address := range addresses {
		go waitForPing(address, id, done)
	}

	for remaining > 0 {
		select {
		case id := <-done:
			remaining -= 1
			doneArr[id] = true
			log.Println("Recieved:", addresses[id])
		case <-ticker:
			log.Println("Currently waiting on:")
			for id := range addresses {
				if !doneArr[id] {
					log.Println("--", addresses[id])
				}
			}
		}
	}

	cmd := exec.Command(command)
	err := cmd.Start()
	if err != nil {
		log.Println("Error starting command: ", err.Error())
	}
}

func waitForPing(address string, id int, done chan int) {
	for {
		cmd := pinger(address)
		err := cmd.Start()
		if err == nil {
			err = cmd.Wait()
		}
		if err == nil {
			break
		}
		if err.Error() != "exit status 1" {
			log.Println("Error in pining", address, "ignoring.")
			break
		}
	}
	done <- id
}
