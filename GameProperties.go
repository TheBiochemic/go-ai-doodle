package main

import (
	"os"
	"fmt"
	"strings"
	"strconv"
)

type gameProperties struct {
	maxPeoplePerHouse    uint8
	mapSizeAreas         uint8
	resourcePerArea      uint16
	motherloadCount      uint8
	motherloadPercentage float32
	buildingCount        uint8
}

/**
	Creates a pre-filled game Properties struct
	If there are Program arguments submitted, they are used for
	replacing the default values
 */
func SetupGameProperties() *gameProperties {
	gameProp := new(gameProperties)

	gameProp.maxPeoplePerHouse = 1
	gameProp.mapSizeAreas = 4
	gameProp.resourcePerArea = 100
	gameProp.motherloadCount = 1
	gameProp.motherloadPercentage = 0.5
	gameProp.buildingCount = 1

	gameProp.updateFromArguments()

	return gameProp
}

func (this *gameProperties) updateFromArguments() {
	if (len(os.Args) > 1) {
		for _, arg := range os.Args[1:] {
			currArg := strings.Split(arg, "=")
			if (len(currArg) > 1) {
				this.updateArgument(currArg[0], currArg[1])
			}
		}
	}
}

func (this *gameProperties) updateArgument(name string, parameter string) {
	switch strings.ToLower(name) {

	case "maxpeopleperhouse":
		newint, _ := strconv.Atoi(parameter)
		this.maxPeoplePerHouse = uint8(newint)
		break

	case "mapsizeareas":
		newint, _ := strconv.Atoi(parameter)
		this.mapSizeAreas = uint8(newint)
		break

	case "resourceperarea":
		newint, _ := strconv.Atoi(parameter)
		this.resourcePerArea = uint16(newint)
		break

	case "motherloadcount":
		newint, _ := strconv.Atoi(parameter)
		this.motherloadCount = uint8(newint)
		break

	case "motherloadpercentage":
		newfloat, _ := strconv.ParseFloat(parameter, 32)
		this.motherloadPercentage = float32(newfloat)
		break

	case "buildingcount":
		newint, _ := strconv.Atoi(parameter)
		this.buildingCount = uint8(newint)
		break

	default:
		fmt.Println("unrecognized parameter \"" + name + "\"")
	}
}
