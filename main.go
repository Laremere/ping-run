package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

//Wait for a set of pings to return successfully, then run a command
func main() {
	//Check arguments
	if len(os.Args) < 3 {
		log.Println("Proper usage:")
		log.Println("ping-run command-to-execute address-to-ping [additional-addresses-to-ping...]")
		log.Println("Failure, exiting")
		return
	}
	//Read arguments
	command := os.Args[1]
	addresses := os.Args[2:]

	//How many addresses to ping
	remaining := len(addresses)
	//Has the address returned a successful ping yet?
	doneArr := make([]bool, remaining)
	//Communication channel back from pingers when done.
	done := make(chan int)
	//Every 15 seconds print who we're waiting on still
	ticker := time.Tick(time.Second * 15)

	//Run a go routine for each addres to ping
	for id, address := range addresses {
		go waitForPing(address, id, done)
	}

	//While we still have addresses left
	for remaining > 0 {
		select {
		case id := <-done:
			//Address successfully pinged, say so and update variables
			remaining -= 1
			doneArr[id] = true
			log.Println("Recieved:", addresses[id])
		case <-ticker:
			//Every 15 seconds print who we're waiting on still
			log.Println("Currently waiting on:")
			for id := range addresses {
				if !doneArr[id] {
					log.Println("--", addresses[id])
				}
			}
		}
	}

	//Run the command and exit, letting the command run without waiting on it
	cmd := exec.Command(command)
	err := cmd.Start()
	if err != nil {
		log.Println("Error starting command: ", err.Error())
	}
}

//Repeatedly ping the address until it responds, and then
//send id over done channel
func waitForPing(address string, id int, done chan int) {
	for {
		//Start OS specific ping command
		cmd := pinger(address)
		err := cmd.Start()
		if err == nil {
			err = cmd.Wait()
		}
		//Successful return
		if err == nil {
			break
		}
		//If the error is "exit status 1", then on linux the command has timed out
		//On windows there may have been some other error, but it's not easy to determine
		//which has happened.  So we just ignore "exit status 1" until ping is succcessful
		if err.Error() != "exit status 1" {
			log.Println("Error in pining", address, "ignoring.")
			break
		}
	}
	//Tell main that we're done pinging id.
	done <- id
}
