package main

import (
	"fmt"
	"math/rand"
)

type human struct {
	name       string
	maxHealth  uint16
	currHealth uint16
	fitness    float32

	carryAmount float64
	carryType   int

	home   *staticObject
	action interaction
	x      float64
	y      float64
}

var humanNames []string = []string{
	"Abe", "Adalberto", "Adan", "Agustin", "Allan", "Alva", "Alvaro", "Amos",
	"Arnulfo", "Arron", "Arthur", "Aubrey", "Barry", "Bertram", "Blaine",
	"Booker", "Brady", "Brain", "Bret", "Carmine", "Carol", "Carter", "Chi",
	"Christoper", "Chuck", "Chung", "Clarence", "Clement", "Conrad", "Cory",
	"Craig", "Cristopher", "Daniel", "Dannie", "Darin", "Darrell", "Demetrius",
	"Derick", "Deshawn", "Dion", "Dominic", "Don", "Donovan", "Dorian",
	"Doyle", "Dustin", "Dwight", "Earnest", "Edmond", "Eduardo", "Efren",
	"Elton", "Emil", "Emmett", "Erick", "Erik", "Erin", "Faustino", "Fausto",
	"Foster", "Frances", "Frank", "Frederick", "Gail", "Gale", "Garry",
	"Gavin", "Genaro", "Geraldo", "Gino", "Gregory", "Gustavo", "Hans",
	"Harrison", "Henry", "Hobert", "Hollis", "Hong", "Horacio", "Hoyt",
	"Huey", "Ian", "Ignacio", "Isaiah", "Isaias", "Isidro", "Isreal", "Ivan",
	"Jamaal", "Jamey", "Jamison", "Jarrett", "Jarvis", "Jefferey", "Jess",
	"Jewel", "Jimmy", "Joey", "Jonathan", "Joshua", "Juan", "Julio", "Junior",
	"Justin", "Kelley", "Kendall", "Kevin", "Kim", "Kurt", "Kyle", "Lawerence",
	"Lenny", "Leonardo", "Leonel", "Leroy", "Leslie", "Levi", "Louie", "Louis",
	"Loyd", "Luigi", "Manuel", "Marco", "Marcus", "Marlon", "Matthew", "Mauricio",
	"Michael", "Milton", "Mohamed", "Monte", "Myron", "Neil", "Nicky", "Noble",
	"Noel", "Norman", "Norris", "Olin", "Oscar", "Pablo", "Pasquale", "Peter",
	"Quinn", "Quinton", "Rafael", "Randal", "Raphael", "Reed", "Refugio", "Reggie",
	"Reginald", "Reuben", "Ricardo", "Ricky", "Rico", "Riley", "Robbie", "Rod",
	"Roderick", "Rodney", "Rodolfo", "Rolando", "Ronny", "Rory", "Rosendo",
	"Royal", "Salvatore", "Scott", "Shannon", "Sherwood", "Shirley", "Stacy",
	"Stan", "Stanley", "Stephan", "Ted", "Terence", "Theron", "Trenton",
	"Tyson", "Vincent", "Von", "Wally", "Walter", "Waylon", "Wiley", "Williams",
	"Willie", "Zachariah", "Zachary", "Zane",
}

func createHuman() *human {
	humanName := humanNames[rand.Uint64()%uint64(len(humanNames))]
	fitnessMod := rand.Float64()
	return &human{humanName, 100, 100, float32(0.5 + fitnessMod), 0, -1, nil, nil, 0, 0}
}

func (a human) isEqual(b *human) bool {
	if a.x == b.x &&
		a.y == b.y &&
		a.name == b.name &&
		a.maxHealth == b.maxHealth &&
		a.fitness == b.fitness {
		return true
	}

	return false
}

func (p *human) setAction(action interface{}) {

	if (action == nil) {
		p.action = nil
		return
	}

	a, ok := action.(interaction)
	if !ok {
		fmt.Println("Could not set new interaction for", p.name+"!")
		return
	}
	p.action = a
}

func (p *human) setHome(obj interface{}) {
	home, ok := obj.(staticObject)
	if !ok {
		fmt.Println("Could not set Home for", p.name+"!")
		return
	}

	p.home = &home
}

func (p *human) act(landscape *region) bool {
	if p.action != nil {
		action := p.action
		finished := p.action.do(p, landscape)
		if action == p.action && finished {
			p.action = nil
		}
		return true
	}

	return false
}
