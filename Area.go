package main

import (
	"fmt"
	"math/rand"
)

const REGIONSIZE uint = 32

type region struct {
	areas               []*area
	size                uint16
	sizeAreas           uint8
	inhabitantsPerHouse uint8
}

func (r *region) removeStaticObject(object interface{}) {
	j, ok := object.(staticObject)
	if !ok {
		return
	}

	theArea := r.getAreaFromCoords(j.getX(), j.getY())
	theArea.removeStaticObject(object)
}

func (r *region) getAreaFromAreaCoords(areax int, areay int) *area {
	if areax >= int(r.sizeAreas) || areax < 0 {
		return nil
	}
	if areay >= int(r.sizeAreas) || areay < 0 {
		return nil
	}

	return r.areas[areax+int(r.sizeAreas)*areay]
}

func (r *region) getAreaFromCoords(x float64, y float64) *area {
	return r.getAreaFromAreaCoords(int(x/float64(REGIONSIZE)), int(y/float64(REGIONSIZE)))
}

func constructRegion(properties *gameProperties) *region {

	humanCount := 0
	buildingCount := 0
	resourceCount := 0

	reg := new(region)
	reg.size = uint16(properties.mapSizeAreas) * uint16(REGIONSIZE)
	reg.sizeAreas = properties.mapSizeAreas
	reg.inhabitantsPerHouse = properties.maxPeoplePerHouse

	fullArea := reg.sizeAreas * reg.sizeAreas
	reg.areas = []*area{}
	for i := uint8(0); i < fullArea; i++ {
		newArea := createArea(uint16(uint(i)%uint(reg.sizeAreas)), uint16(uint(i)/uint(reg.sizeAreas)))
		reg.areas = append(reg.areas, &newArea)
	}

	//Buildings
	for i := uint8(0); i < properties.buildingCount; i++ {
		newBuilding := createBuilding()
		newBuilding.x = rand.Float64() * float64(reg.size)
		newBuilding.y = rand.Float64() * float64(reg.size)
		reg.putStaticObject(newBuilding)
		buildingCount++

		//Humans
		for j := uint8(0); j < reg.inhabitantsPerHouse; j++ {
			newHuman := createHuman()
			newHuman.x = rand.Float64() * float64(reg.size)
			newHuman.y = rand.Float64() * float64(reg.size)
			newHuman.setHome(newBuilding)
			newHuman.setAction(interactionMoveToObject{newHuman.home})

			reg.putHuman(newHuman)
			humanCount++
		}
	}

	//Resources
	woodMotherload := reg.placeResource(RES_WOOD, float64(properties.motherloadPercentage)*float64(properties.resourcePerArea)*float64(len(reg.areas)))
	stoneMotherload := reg.placeResource(RES_STONE, float64(properties.motherloadPercentage)*float64(properties.resourcePerArea)*float64(len(reg.areas)))
	foodMotherload := reg.placeResource(RES_FOOD, float64(properties.motherloadPercentage)*float64(properties.resourcePerArea)*float64(len(reg.areas)))
	woodMotherload = float64(properties.resourcePerArea)*float64(len(reg.areas)) - woodMotherload
	stoneMotherload = float64(properties.resourcePerArea)*float64(len(reg.areas)) - stoneMotherload
	foodMotherload = float64(properties.resourcePerArea)*float64(len(reg.areas)) - foodMotherload
	resourceCount += 3

	maxWood := woodMotherload / 100
	maxStone := stoneMotherload / 100
	maxFood := foodMotherload / 100

	for ; woodMotherload > 0; woodMotherload -= reg.placeResource(RES_WOOD, maxWood) {
		resourceCount++
	}
	for ; stoneMotherload > 0; stoneMotherload -= reg.placeResource(RES_STONE, maxStone) {
		resourceCount++
	}
	for ; foodMotherload > 0; foodMotherload -= reg.placeResource(RES_FOOD, maxFood) {
		resourceCount++
	}

	fmt.Println("created Region with", len(reg.areas), "areas,", humanCount, "humans,", buildingCount, "buildings and", resourceCount, "resources.")

	return reg
}

func (r *region) reCountRegion() {
	humanCount := 0
	buildingCount := 0
	resourceCount := 0

	for _, currArea := range r.areas {
		for _, currHuman := range currArea.people {
			if currHuman != nil {
				humanCount++
			}
		}

		for _, currObject := range currArea.staticObjects {
			if currObject != nil {
				if (*currObject).isType(BUILDING) {
					buildingCount++
				} else {
					resourceCount++
				}
			}
		}
	}

	fmt.Println("Region has", len(r.areas), "areas,", humanCount, "humans,", buildingCount, "buildings and", resourceCount, "resources.")
}

func (r *region) placeResource(ofType int, amount float64) float64 {
	actualAmount := amount * (0.75 + (0.5 * rand.Float64()))
	newRes := createResource(ofType, actualAmount)
	newRes.x = rand.Float64() * float64(r.size)
	newRes.y = rand.Float64() * float64(r.size)
	r.putStaticObject(newRes)
	return actualAmount
}

func (reg *region) update() uint {

	var operations uint = 0
	for _, element := range reg.areas {
		operations += element.update(reg)
	}
	return operations
}

func (reg *region) putHuman(person *human) {

	if person.x > float64(reg.size) {
		person.x = float64(reg.size)
	}
	if person.x < 0 {
		person.x = 0
	}
	if person.y > float64(reg.size) {
		person.y = float64(reg.size)
	}
	if person.y < 0 {
		person.y = 0
	}
	var newRegionX uint8 = uint8(person.x / float64(REGIONSIZE))
	var newRegionY uint8 = uint8(person.y / float64(REGIONSIZE))
	newArea := reg.areas[newRegionX+reg.sizeAreas*newRegionY]
	newArea.addHuman(person)
}

func (reg *region) putStaticObject(obj interface{}) {
	sobj, ok := obj.(staticObject)
	if !ok {
		fmt.Println("Could not put static object into region!")
		return
	}

	cArea := reg.getAreaFromCoords(sobj.getX(), sobj.getY())
	cArea.staticObjects = append(cArea.staticObjects, &sobj)
}

func (reg *region) move(person *human, toX float64, toY float64) {

	if toX > float64(reg.size) {
		toX = float64(reg.size)
	}
	if toX < 0 {
		toX = 0
	}
	if toY > float64(reg.size) {
		toY = float64(reg.size)
	}
	if toY < 0 {
		toY = 0
	}

	oldArea := reg.getAreaFromCoords(person.x, person.y)
	newArea := reg.getAreaFromCoords(toX, toY)

	if !newArea.isEqual(oldArea) {
		//fmt.Println(person.name, "enters other Area.")
		oldArea.removeHuman(person)
		newArea.addHuman(person)
	}

	person.x = toX
	person.y = toY
}

type area struct {
	x uint16
	y uint16

	staticObjects []*staticObject
	people        []*human
}

func (a area) String() string {
	var people string = "{"
	for _, person := range a.people {
		people = fmt.Sprint(people, person, ",")
	}
	people = fmt.Sprint(people, "}")
	return fmt.Sprint("{", a.x, a.y, people)
}

func (r *region) String() string {
	var data string = fmt.Sprint(r.size, ",", r.sizeAreas, ",")
	for _, currArea := range r.areas {
		data = fmt.Sprint(data, currArea.String(), ",")
	}
	return "{" + data + "}"
}

func (a *area) update(world *region) uint {
	var operations uint = 0
	for _, person := range a.people {
		if person != nil {

			if person.act(world) {
				operations++
			}
		}
	}

	return operations
}

func (r *region) getNearestStaticObject(x float64, y float64, objType int) *staticObject {

	//Fast'n Dirty
	for _, a := range r.areas {
		for _, staticObj := range a.staticObjects {
			if (*staticObj).isType(objType) {
				return staticObj
			}
		}
	}

	//Make faster Version with LIFO Container and circular search from starting point
	return nil
}

func createArea(x uint16, y uint16) area {
	a := area{}
	a.x = x
	a.y = y
	a.people = []*human{}
	a.staticObjects = []*staticObject{}
	return a
}

func (a area) isEqual(b *area) bool {

	if b == nil {
		return false
	}

	if a.x == b.x && a.y == b.y {
		return true
	}
	return false
}

func (a *area) removeHuman(person *human) {
	foundIndex := 0
	for i, found := range a.people {
		if found.isEqual(person) {
			foundIndex = i
			break
		}
	}

	copy(a.people[foundIndex:], a.people[foundIndex+1:])
	a.people[len(a.people)-1] = nil
	a.people = a.people[:len(a.people)-1]
}

func (a *area) removeStaticObject(object interface{}) {

	sobj, ok := object.(staticObject)
	if !ok {
		fmt.Println("Could not remove static object from region!")
		return
	}

	foundIndex := -1
	for i, found := range a.staticObjects {
		if *found == sobj {
			foundIndex = i
			break
		}
	}
	if foundIndex != -1 {
		copy(a.staticObjects[foundIndex:], a.staticObjects[foundIndex+1:])
		a.staticObjects[len(a.staticObjects)-1] = nil
		a.staticObjects = a.staticObjects[:len(a.staticObjects)-1]
	}
}

func (a *area) addHuman(person *human) {
	a.people = append(a.people, person)
}
