//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//
//                                                                          //
//    ,  /\  .       Tages - Yet another adaptive EC-based player for IPD   //
//   //`-||-'\\                                                             //
//  (| -AEGM- |)     Boldly crafted in Go between 2014 and 2015             //
//   \\,-||-.//      by Giovanni Squillero <giovanni.squillero@polito.it>   //
//    `  ||  '       and Alberto Tonda, Elio Piccolo & Marco Gaudesi        //
//       ||                                                                 //
//       ||          "You don't need to have a big dream: be µ-ambitious!"  //
//       ()                                           -- Tim Minchin        //
//                                                                          //
//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//\\//

// This file is part of Tages
// Copyright © 2015 Giovanni Squillero
// GitHub page: https://github.com/squillero/tages
//
// Tages is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Tages is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Tages.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"log"
	"math/rand"
)

//////////////////////////////////////////////////////////////////////////////
// Types, constants Interfaces

type Population struct {
	MUtations      MUtations
	MU             int
	Lamda          int
	NU             int
	Opponents      []Player
	Turns          []int
	Individual     []Individual
	_champion      Individual
	TotIndividuals int
	Paths          []GameTrace
	Traces         []GameTrace
	ζ              int
}

func (p *Population) String() string {
	min, max := p.Turns[0], p.Turns[0]
	for _, t := range p.Turns {
		if max < t {
			max = t
		}
		if min > t {
			min = t
		}
	}
	return fmt.Sprintf("(%d:%d+%d) population; train: %v matches [%v-%v] against %d sparring mates", p.NU, p.MU, p.Lamda, len(p.Turns), min, max, len(p.Opponents))
}

func CreatePopulation(mu, nu, lamda, mD int, mNd float64, o []Player, t []int) Population {
	var p Population

	p.MUtations = NewMUtations()
	p.Individual = make([]Individual, nu)
	p.MU = mu
	p.Lamda = lamda
	p.NU = nu
	p.Opponents = o
	p.Turns = t
	p.Paths = make([]GameTrace, 0)
	p.Traces = make([]GameTrace, 0)
	p.ζ = 1

	for u := 0; u < nu; u++ {
		p.TotIndividuals++
		p.Individual[u] = RandomIndividual(fmt.Sprintf("I_%d", u), mD, mNd)
	}
	//p.Champion = RandomIndividual("The Boss", mD, mNd)
	//p.Champion.Fit = Fitness{Valid: true, Competitiveness: 0, Coherence: 0, Compactness: 0}
	//p.Invalidate()
	//p.Evaluate()
	//p.Slaughter()

	return p
}

//////////////////////////////////////////////////////////////////////////////
// Simplistic tournament selection
// TOURNAMENT_SIZE must be adjusted @ compile time
func (p *Population) Check() {
	for _, i := range p.Individual {
		i.Check()
	}
}

//////////////////////////////////////////////////////////////////////////////
// Simplistic tournament selection
// TOURNAMENT_SIZE must be adjusted @ compile time

func (p *Population) Select() Individual {
	winner := p.Individual[rand.Int()%len(p.Individual)]
	for t := 1; t < TOURNAMENT_SIZE; t++ {
		opponent := p.Individual[rand.Int()%len(p.Individual)]
		if FitCompare(winner.Fit, opponent.Fit) < 0 {
			winner = opponent
		}
	}
	return winner
}

//////////////////////////////////////////////////////////////////////////////
// Grow & shrink population

func (p *Population) Begat() {
	offspring := make([]Individual, 0, p.Lamda)
	for o := 0; o < p.Lamda; o++ {
		parent := p.Select()
		p.TotIndividuals++
		i := parent.Duplicate()
		i.Name = fmt.Sprintf("I_%v", p.TotIndividuals)
		i.Parent.Fit = parent.Fit
		i.Fit.Valid = false

		var success bool
		if rand.Float32() < .05 {
			p := p.Select()
			XOverAddIndividual(&i, &p)
			i.Parent.MUt = &StaticXOver
			success = true
		} else {
			op := p.MUtations.Select()
			i.Parent.MUt = op
			success = op.F(&i)
		}
		if success {
			offspring = append(offspring, i)
		}
	}
	for i := range offspring {
		p.Individual = append(p.Individual, offspring[i])
	}
}

func (p *Population) Extinguish() {
	p.Individual = p.Individual[:p.MU/5]
	for len(p.Individual) < NU {
		p.TotIndividuals++
		new := RandomIndividual(fmt.Sprintf("E_%d", p.TotIndividuals), 5, .2)
		new.Fit.Valid = false
		p.Individual = append(p.Individual, new)
	}
}

func (p *Population) Slaughter() {
	p.Individual = p.Individual[:p.MU]
}

// invalidate all fitness's
func (p *Population) Invalidate() {
	for i := range p.Individual {
		p.Individual[i].Fit.Valid = false
	}
}

//////////////////////////////////////////////////////////////////////////////

func (p *Population) Dump() {
	for _, i := range p.Individual {
		log.Println(&i)
		i.DumpPNG()
	}
}

//////////////////////////////////////////////////////////////////////////////

func (p *Population) UpdateOperators() {
	// Cleanup
	for m := range p.MUtations.Available {
		p.MUtations.Available[m].W = 0
	}
	for _, i := range p.Individual {
		if i.Parent.MUt.F != nil {
			i.Parent.MUt.Inc()
		} else {
			for m := range p.MUtations.Available {
				p.MUtations.Available[m].W++
			}
		}
	}

	//log.Println("** NEW: **")
	//for i := range p.MUtations.Available {
	//	//p.MUtations.Available[i].Inc()
	//	log.Println(p.MUtations.Available[i])
	//}
	//log.Panicln("...")
}

//////////////////////////////////////////////////////////////////////////////

func (p *Population) AddStep(round int, pathBit, footstepBit PPly) {
	p.Paths[round] = append(p.Paths[round], pathBit)
	p.Traces[round] = append(p.Traces[round], footstepBit)
}

func (p *Population) VanillaRL() PPly {
	reward := [2]int{1, 1}

	for g := range p.Traces {
		for s := range p.Traces[g] {
			r := Payoff[p.Paths[g][s]][p.Traces[g][s]]
			if p.Paths[g][s] == COOPERATE {
				reward[0] += r
			} else {
				reward[1] += r
			}
		}
	}
	log.Println("VanillaRL™ current status:", reward)
	if rand.Intn(reward[0]+reward[1]) < reward[0] {
		return COOPERATE
	} else {
		return DEFECT
	}
}

//////////////////////////////////////////////////////////////////////////////
// Calculate ζ coefficient

func (p *Population) Calculateζ(hit, miss int) {
	for hit+miss < 20 {
		hit++
		miss++
	}
	p.ζ = int(100 * float64(hit) / float64(miss+hit))
}

func (p *Population) Tagζ(i int) string {
	if len(p.Individual[i].ActiveNodes()) > p.ζ {
		return " [ζ]"
	} else {
		return ""
	}
}

//////////////////////////////////////////////////////////////////////////////
// ????

//func (p *Population) Insert(i Individual) {
//	p.Individual = append(p.Individual, i)
//}

// wtf?
//func (p *Population) Contains(i Individual) int {
//	for b := 0; b < len(p.Individual); b++ {
//		if p.Individual[b].NdFsm.CompareTo(i.NdFsm) {
//			return b
//		}
//	}
//	return -1
//}

//////////////////////////////////////////////////////////////////////////////
// Log!
//////////////////////////////////////////////////////////////////////////////
func jazz() {
	log.Panicln("Dummy!")
}
