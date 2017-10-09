package main

import (
	"math"
	"math/rand"
	"fmt"
)

const BUILDING int = 0
const RES_WOOD int = 1
const RES_STONE int = 2
const RES_FOOD int = 3

type staticObject interface {
	getX() float64
	getY() float64
	interact(interactor *human, world *region)
	isType(int) bool
	getName() string
}

type building struct {
	x     float64
	y     float64
	wood  float64
	stone float64
	food  float64
}

func createBuilding() *building {
	return &building{0, 0, 0, 0, 0}
}

func (b *building) getX() float64 {
	return b.x
}

func (b *building) getY() float64 {
	return b.y
}

func (b *building) interact(interactor *human, world *region) {

	if interactor.carryType == RES_WOOD {
		b.wood += interactor.carryAmount
	}

	if interactor.carryType == RES_STONE {
		b.stone += interactor.carryAmount
	}

	if interactor.carryType == RES_FOOD {
		b.food += interactor.carryAmount
	}

	if b.wood >= 15 && b.stone >= 15 && b.food >= 15 {
		interactor.setAction(interactionMoveToBuildingSite{
			rand.Float64() * float64(world.size),
			rand.Float64() * float64(world.size)})
		b.wood -= 15
		b.stone -= 15
		b.food -= 15
		return
	}

	amount := b.wood
	rtype := RES_WOOD

	if b.food < amount {
		amount = b.food
		rtype = RES_FOOD
	}

	if b.stone < amount {
		amount = b.stone
		rtype = RES_STONE
	}

	nearest := world.getNearestStaticObject(interactor.x, interactor.y, rtype)
	interactor.setAction(interactionMoveToObject{nearest})

}

func (b *building) isType(thetype int) bool {
	if thetype == BUILDING {
		return true
	}
	return false
}

func (b *building) getName() string {
	return "House"
}

type resource struct {
	x            float64
	y            float64
	resourceType int
	amount       float64
}

func createResource(ofType int, amount float64) *resource {
	return &resource{0, 0, ofType, amount}
}

func (r *resource) getX() float64 {
	return r.x
}

func (r *resource) getY() float64 {
	return r.y
}

func (r *resource) interact(interactor *human, world *region) {

	interactor.carryAmount = math.Min(float64(interactor.fitness)*4, r.amount)
	interactor.carryType = r.resourceType

	if interactor.carryAmount >= r.amount {
		r.amount = 0
		world.removeStaticObject(r)
	} else {
		r.amount -= interactor.carryAmount
	}

	if interactor.carryAmount > 0 {
		fmt.Println(interactor.name, "has collected around", int(interactor.carryAmount), r.getName())
		interactor.setAction(interactionMoveToObject{interactor.home})
	} else {
		fmt.Println(interactor.name, "has nothing found here. Searching for another", r.getName())
		nearest := world.getNearestStaticObject(interactor.x, interactor.y, r.resourceType)
		interactor.setAction(interactionMoveToObject{nearest})
	}
}

func (r *resource) isType(thetype int) bool {
	if thetype == r.resourceType {
		return true
	}
	return false
}

func (r *resource) getName() string {
	switch r.resourceType {
	case RES_FOOD:
		return "Berry Bush"
	case RES_WOOD:
		return "Tree"
	case RES_STONE:
		return "Rock"
	}

	return "Unknown Resource"
}
