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
	"io/ioutil"
	"log"
	"math/rand"
)

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// Types, constants Interfaces
//////////////////////////////////////////////////////////////////////////////

type Parent struct {
	Fit Fitness
	MUt *GeneticOperator
}

type Individual struct {
	NdFsm
	Inverse [3]NdFsm
	Fit     Fitness
	Parent  Parent
}

func (ind *Individual) String() string {
	t, a, ar, n := 0, 0, 0, 0
	for _, s := range ind.Node {
		t++
		if s.Active {
			a++
		}
		if s.Active && s.reachable {
			ar++
		}
		if len(s.OnCooperation)+len(s.OnDefection) > 2 {
			n++
		}
	}
	return fmt.Sprintf("\"%s\" LM#%v %d/%d nodes, nd: %.2f; fitness: %v", ind.Name, ind.Parent.MUt.N, a, ar, float64(n)/float64(t), &ind.Fit)
}

func (ind *Individual) DumpPNG() {
	fsm := ind.NdFsm.Duplicate()
	//log.Println(&tfsm)
	//fsm.Name += " - " + ind.Fit.String()
	ioutil.WriteFile(ind.Name+".gv", fsm.GvEncode(), 0644)
	ioutil.WriteFile(ind.Name+".json", fsm.JsonEncode(), 0644)
	go Gv2Png(ind.Name)
}

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// Spare utilites
//////////////////////////////////////////////////////////////////////////////

func (ind *Individual) Check() bool {
	ind.NdFsm.Check()
	if ind.Fit.Valid {
		ind.Inverse[COOPERATE].Check()
		ind.Inverse[DEFECT].Check()
	}
	return true
}

func (ind *Individual) Cleanup() {
	for t := range ind.Node {
		ind.Node[t].Active = ind.Node[t].reachable
	}
}

// Duplicate (aka. DeepCopy)
func (ind *Individual) Duplicate() Individual {
	new_ind := new(Individual)
	new_ind.Fit = ind.Fit
	new_ind.Parent.Fit = ind.Parent.Fit
	new_ind.Parent.MUt = ind.Parent.MUt
	new_ind.NdFsm = ind.NdFsm.Duplicate()
	new_ind.Inverse[COOPERATE] = ind.Inverse[COOPERATE].Duplicate()
	new_ind.Inverse[DEFECT] = ind.Inverse[DEFECT].Duplicate()
	return *new_ind
}

//////////////////////////////////////////////////////////////////////////////
// Create a random individual named "name", with "states" states
// and a non-determinism of nd
func RandomIndividual(name string, states int, nd float64) Individual {
	var i Individual
	var b [2]PPly = [...]PPly{COOPERATE, DEFECT}

	i.Fit.Valid = false
	i.Parent.MUt = &GeneticOperator{N: "Ancestor"}
	i.NdFsm = NewNdFSM(name)

	for t := 0; t < states; t++ {
		i.AddNode(b[rand.Intn(2)])
		i.Node[t].color = "blue"
	}
	nodes := i.ActiveNodes()
	i.SetInitialState(nodes[rand.Intn(len(nodes))])
	for _, t := range nodes {
		i.AddDTransition(t, COOPERATE, nodes[rand.Intn(len(nodes))])
		i.AddDTransition(t, DEFECT, nodes[rand.Intn(len(nodes))])
	}

	// spice up
	for r := rand.Float64(); r < nd; r = rand.Float64() {
		MUtationAddNdTransition(&i)
	}

	i.Check()
	return i
}

//////////////////////////////////////////////////////////////////////////////
// Creates the two inverse fsms
func (ind *Individual) CreateInverse() {
	// create inverse
	inverse := NewBareNdFSM("Jabberwocky")
	for _, n := range ind.Node {
		inverse.AddNode(n.Action)
	}
	for n := range ind.Node {
		if !ind.Node[n].Active {
			inverse.Node[n].Active = false
		}
		if !ind.Node[n].reachable && ind.Node[n].Action != INVALID {
			inverse.Node[n].Active = false
		}
	}
	for from, From := range ind.Node {
		if From.Active && From.reachable && From.Action != INVALID {
			for _, to := range ind.Node[from].OnCooperation {
				inverse.AddNdTransition(to, COOPERATE, NodeIndex(from))
			}
			for _, to := range ind.Node[from].OnDefection {
				inverse.AddNdTransition(to, DEFECT, NodeIndex(from))
			}
		}
	}
	for _, n := range inverse.ActiveNodes() {
		if len(inverse.Node[n].OnCooperation) < 1 {
			inverse.AddNdTransition(NodeIndex(n), COOPERATE, inverse.GetSink())
		}
		if len(inverse.Node[n].OnDefection) < 1 {
			inverse.AddNdTransition(NodeIndex(n), DEFECT, inverse.GetSink())
		}
	}

	for _, s := range []PPly{DEFECT, COOPERATE} {
		ind.Inverse[s] = inverse.Duplicate()
		ind.Inverse[s].Name = ind.Name + " + i_" + s.String()
		//ind.Inverse[s].InitialState = make([]NodeIndex, 0)
		for n := range ind.Node {
			ind.Inverse[s].Node[n].color = ind.Node[s].color
			if ind.Node[n].Active && ind.Node[n].reachable && ind.Node[n].Action == s {
				ind.Inverse[s].AddInitialState(NodeIndex(n))
			}
		}
		// patch!
		sink := inverse.GetSink()
		ind.Inverse[s].AddNdTransition(NodeIndex(sink), COOPERATE, NodeIndex(sink))
		ind.Inverse[s].AddNdTransition(NodeIndex(sink), DEFECT, NodeIndex(sink))
		if len(ind.Inverse[s].InitialState) < 1 {
			ind.Inverse[s].SetInitialState(sink)
		}
		ind.Inverse[s].UpdateInternals()
	}
}

func UselessLog() {
	log.Panicln("...")
}
