package main

import (
	"fmt"
	"time"
	"os"
	"bufio"
)

var landscape *region

func main() {
	go doLogic()
	fmt.Println("Press ENTER to exit Simulation.")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadByte()
}

func doLogic() {
	properties := SetupGameProperties()
	landscape = constructRegion(properties)

	for ; true; {
		time.Sleep(50 * time.Millisecond)
		if landscape.update() == 0 {
			break
		}
	}
	landscape.reCountRegion()
	fmt.Println("Press ENTER to exit Simulation.")
}
