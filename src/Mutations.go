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
// a bit of self tuning

const (
	MAX_W = 50
	MIN_W = 0
)

type GeneticOperator struct {
	F func(i *Individual) bool
	W int
	N string
}

type MUtations struct {
	Available []GeneticOperator
}

func (m *GeneticOperator) Inc() {
	if m.W < MAX_W {
		m.W++
	}
}
func (m *GeneticOperator) Dec() {
	if m.W > MIN_W+1 {
		m.W--
	}
}

func (o *MUtations) AddOperator(f func(i *Individual) bool, n string) {
	o.Available = append(o.Available, GeneticOperator{F: f, W: MIN_W, N: n})

}

func (o *MUtations) Select() *GeneticOperator {
	var i int

	if rand.Float32() < 0.25 {
		i = rand.Intn(len(o.Available))
	} else {
		totW := 1
		for _, m := range o.Available {
			totW += m.W
		}
		i = 0
		for r := rand.Intn(totW) - o.Available[i].W; r > 0; r -= o.Available[i].W {
			i++
		}
	}
	return &o.Available[i]
}

func NewMUtations() MUtations {
	m := new(MUtations)
	m.Available = make([]GeneticOperator, 0)
	return *m
}

//////////////////////////////////////////////////////////////////////////////
// interfaces

func (m *GeneticOperator) String() string {
	return fmt.Sprintf("%s: %d", m.N, m.W)
}

func (M *MUtations) String() string {
	s := ""
	for _, m := range M.Available {
		s += fmt.Sprintf(" %v", &m)
	}
	return s
}

//////////////////////////////////////////////////////////////////////////////
// names almost self explanatory

func MUtationChangeInitialState(i *Individual) bool {
	// log.Printf("mutationChangeInitialState: %v\n", i)
	pt := i.ActiveNodes()
	i.SetInitialState(pt[rand.Intn(len(pt))])
	return true
}

func MUtationAddInitialState(i *Individual) bool {
	// log.Printf("mutationAddInitialState: %v\n", i)
	pt := make([]NodeIndex, 0, len(i.Node))
	for ni, nn := range i.Node {
		if i.IsInitialState(NodeIndex(ni)) == false && nn.Active == true && nn.Action != INVALID {
			pt = append(pt, NodeIndex(ni))
		}
	}
	if len(pt) >= 1 {
		i.AddInitialState(pt[rand.Intn(len(pt))])
	}
	return true
}

func MUtationRemoveInitialState(i *Individual) bool {
	if len(i.InitialState) > 1 {
		i.RemoveInitialState(i.InitialState[rand.Intn(len(i.InitialState))])
		return true
	} else {
		return false
	}
}

func MUtationAddNode(i *Individual) bool {
	// log.Printf("mutationAddNode: %v\n", i)
	pt := i.ActiveNodes()

	nn := i.AddNode([2]PPly{COOPERATE, DEFECT}[rand.Intn(2)])
	i.Node[nn].Active = true

	// link it (non deterministic)
	i.AddNdTransition(pt[rand.Intn(len(pt))], COOPERATE, nn)

	pt = i.ActiveNodes()
	i.AddNdTransition(nn, COOPERATE, pt[rand.Intn(len(pt))])
	i.AddNdTransition(nn, DEFECT, pt[rand.Intn(len(pt))])

	return true
}

// caveats: transitions towards the deleted node are randomized
func MUtationRemoveNode(i *Individual) bool {
	//log.Printf("mutationRemoveNode: %v\n", i)
	pt := i.ActiveNodes()
	if len(pt) == 1 {
		return false
	}
	n := pt[rand.Intn(len(pt))]

	//log.Printf("Deactivating N%v:%v\n", n, &i.Node[n])
	i.Node[n].Active = false
	i.Node[n].OnDefection = make([]NodeIndex, 0)
	i.Node[n].OnCooperation = make([]NodeIndex, 0)

	// remove from transitions
	for _, t := range i.ActiveNodes() {
		var tmp []NodeIndex
		tmp = make([]NodeIndex, 0, len(i.Node[t].OnCooperation))
		for _, u := range i.Node[t].OnCooperation {
			if u != n {
				tmp = append(tmp, u)
			}
		}
		if len(tmp) == 0 {
			tmp = append(tmp, t)
		}
		i.Node[t].OnCooperation = tmp

		tmp = make([]NodeIndex, 0, len(i.Node[t].OnDefection))
		for _, u := range i.Node[t].OnDefection {
			if u != n {
				tmp = append(tmp, u)
			}
		}
		if len(tmp) == 0 {
			tmp = append(tmp, t)
		}
		i.Node[t].OnDefection = tmp
	}

	// patching initial state (remove inactive)
	i.RemoveInitialState(n)
	if len(i.InitialState) == 0 {
		MUtationAddInitialState(i)
	}

	return true
}

func MUtationChangeTransition(i *Individual) bool {
	// log.Printf("mutationChangeTransition: %v\n", i)
	i.Check()
	pt := i.ActiveNodes()
	node := pt[rand.Intn(len(pt))]
	if rand.Float64() < 0.5 {
		tran := rand.Int() % len(i.Node[node].OnDefection)
		newt := pt[rand.Intn(len(pt))]
		if !i.CheckTransition(node, DEFECT, newt) {
			// log.Printf("Changing transition %d to %d in %v\n", i.Node[node].OnDefection[tran], newt, node)
			i.Node[node].OnDefection[tran] = newt
		}
	} else {
		tran := rand.Int() % len(i.Node[node].OnCooperation)
		newt := pt[rand.Intn(len(pt))]
		if !i.CheckTransition(node, COOPERATE, newt) {
			i.Node[node].OnCooperation[tran] = newt
			// log.Printf("Changing transition %d to %d in %v\n", i.Node[node].OnDefection[tran], newt, node)
		}
	}
	i.Check()
	return true
}

func MUtationChangeStateAction(i *Individual) bool {
	pt := i.ActiveNodes()
	n := pt[rand.Intn(len(pt))]
	if i.Node[n].Action == COOPERATE {
		i.Node[n].Action = DEFECT
	} else {
		i.Node[n].Action = COOPERATE
	}
	return true
}

// spice up fsm by adding some non-determinism
func MUtationAddNdTransition(i *Individual) bool {
	// log.Printf("mutationAddNdTransition: %v\n", i)
	b := [2]PPly{COOPERATE, DEFECT}
	pt := i.ActiveNodes()
	if i.AddDTransition(pt[rand.Intn(len(pt))], b[rand.Intn(2)], pt[rand.Intn(len(pt))]) == false {
		ni := RandomIndividual("R"+i.Name, 5, .2)
		i = &ni
		return true
	} else {
		return false
	}
}

// dull down fsm by removing non-determinism
func MUtationRemoveNdTransition(i *Individual) bool {
	// select node
	nodeWeight := make([]int, len(i.Node))
	totWeight := 0
	for _, ni := range i.ActiveNodes() {
		nodeWeight[ni] = len(i.Node[ni].OnCooperation) + len(i.Node[ni].OnDefection) - 2
		totWeight += nodeWeight[ni]
	}

	// remove transition
	if totWeight > 0 {
		n := 0
		for nodeWeight[n] == 0 {
			n++
		}
		r := rand.Intn(totWeight) - nodeWeight[n]
		for r > 0 {
			n++
			r -= nodeWeight[n]
		}
		//log.Println("Current r:", r, "n:", n)
		//log.Printf("** Before: [Selected node %d]", n)
		//i.Debug()
		if len(i.Node[n].OnCooperation) == 1 || (len(i.Node[n].OnDefection) > 1 && rand.Float64() > 0.5) {
			t := rand.Int() % len(i.Node[n].OnDefection)
			i.Node[n].OnDefection[t] = i.Node[n].OnDefection[len(i.Node[n].OnDefection)-1]
			i.Node[n].OnDefection = i.Node[n].OnDefection[0 : len(i.Node[n].OnDefection)-1]
		} else {
			t := rand.Int() % len(i.Node[n].OnCooperation)
			i.Node[n].OnCooperation[t] = i.Node[n].OnCooperation[len(i.Node[n].OnCooperation)-1]
			i.Node[n].OnCooperation = i.Node[n].OnCooperation[0 : len(i.Node[n].OnCooperation)-1]
		}
		//log.Printf("** After: ")
		//i.Debug()
		//log.Printf("\n\n")
		return true
	} else {
		return false
	}
}

func zop() {
	log.Println("zop")
}
