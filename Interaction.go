package main

import (
	"fmt"
	"math/rand"
)

type interaction interface {
	do(operator *human, landscape *region) bool //returns true, if finished
}

type interactionMoveToLocation struct {
	x float64
	y float64
}

func (i interactionMoveToLocation) do(operator *human, landscape *region) bool {

	maxDistance := operator.fitness * 4;
	if distanceBetween(operator.x, operator.y, i.x, i.y) < float64(maxDistance) {
		landscape.move(operator, i.x, i.y)
		fmt.Println(operator.name, "has arrived at Position (", i.x, ", ", i.y, ").");
		return true
	} else {
		diffX := i.x - operator.x
		diffY := i.y - operator.y
		diffX, diffY = limitToLength(diffX, diffY, float64(maxDistance))
		landscape.move(operator, operator.x+diffX, operator.y+diffY)
		return false
	}
}

type interactionMoveToBuildingSite struct {
	x float64
	y float64
}

func (i interactionMoveToBuildingSite) do(operator *human, landscape *region) bool {

	maxDistance := operator.fitness * 4;
	if distanceBetween(operator.x, operator.y, i.x, i.y) < float64(maxDistance) {
		landscape.move(operator, i.x, i.y)
		operator.setAction(interactionBuildNewHome{})
		return true
	} else {
		diffX := i.x - operator.x
		diffY := i.y - operator.y
		diffX, diffY = limitToLength(diffX, diffY, float64(maxDistance))
		landscape.move(operator, operator.x+diffX, operator.y+diffY)
		return false
	}
}

type interactionBuildNewHome struct {
}

func (b interactionBuildNewHome) do(operator *human, landscape *region) bool {
	nearestBuilding := landscape.getNearestStaticObject(operator.x, operator.x, BUILDING)
	if distanceBetween((*nearestBuilding).getX(), (*nearestBuilding).getY(), operator.x, operator.y) < 10 {
		operator.setAction(interactionMoveToBuildingSite{rand.Float64() * float64(landscape.size), rand.Float64() * float64(landscape.size)})
		fmt.Println(operator.name, "cannot build a Home here, another one is too close.")
		return true
	} else {
		newBuilding := createBuilding()
		newBuilding.x = operator.x
		newBuilding.y = operator.y
		landscape.putStaticObject(newBuilding)
		fmt.Println(operator.name, "has built a new Home")

		for j := uint8(0); j < landscape.inhabitantsPerHouse; j++ {
			newHuman := createHuman()
			newHuman.x = newBuilding.x
			newHuman.y = newBuilding.y
			newHuman.setHome(newBuilding)
			landscape.putHuman(newHuman)
			fmt.Println(newHuman.name, "is now alive; Welcome!")
			newBuilding.interact(newHuman, landscape)
		}
		operator.setAction(interactionMoveToObject{operator.home})
		return true
	}
}

type interactionMoveToObject struct {
	target *staticObject
}

func (i interactionMoveToObject) do(operator *human, landscape *region) bool {

	if i.target == nil {
		fmt.Println(operator.name, "has no Target, so he/she is confused now.");
		return true
	}

	maxDistance := operator.fitness * 4
	targetX := (*i.target).getX()
	targetY := (*i.target).getY()
	if distanceBetween(operator.x, operator.y, targetX, targetY) < float64(maxDistance) {
		landscape.move(operator, targetX, targetY)
		fmt.Println(operator.name, "has arrived at", (*i.target).getName()+".");
		(*i.target).interact(operator, landscape)
		return true
	} else {
		diffX := targetX - operator.x
		diffY := targetY - operator.y
		diffX, diffY = limitToLength(diffX, diffY, float64(maxDistance))
		landscape.move(operator, operator.x+diffX, operator.y+diffY)
		return false
	}
}
