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
	"log"
)

func Pop(g GameTrace) (PPly, GameTrace) {
	e := g[len(g)-1]
	return e, g[:len(g)-1]

}

func zap() {
	log.Panicln("ZAP")
}

// check if step and footstep are compatible
func Compatible(a, b PPly) bool {
	if a == INVALID || b == INVALID {
		return false
	} else if a == b {
		return true
	} else if a == UNSPECIFIED || b == UNSPECIFIED {
		return true
	} else {
		return false
	}
}

//////////////////////////////////////////////////////////////////////////////
// tries to follow backward the footprints on a path
func (fsm *NdFsm) Scout(Path, Footprints GameTrace) (numSteps int, lastState map[NodeIndex]float64) {
	CurrentState := make(map[NodeIndex]float64)

	// startup
	numSteps = 0
	firstMark := Footprints[0]
	Footprints = Footprints[1:]
	for _, s := range fsm.InitialState {
		if Compatible(fsm.Node[s].Action, firstMark) {
			CurrentState[s]++
		}
	}

	for i, step := range Path {
		if len(CurrentState) == 0 {
			return numSteps, CurrentState
		}

		numSteps++

		// next state
		NextState := make(map[NodeIndex]float64)
		for s, m := range CurrentState {
			next := fsm.Node[s].NextState(step)
			for _, n := range next {
				NextState[n] += m
			}
		}

		// strip out invalid or non-coherent states
		for s := range NextState {
			if !Compatible(fsm.Node[s].Action, Footprints[i]) {
				delete(NextState, s)
			}
		}

		CurrentState = NextState
	}
	return numSteps, CurrentState
}

//////////////////////////////////////////////////////////////////////////////
// simulate the game on a given path. returns the order-0 probability to
// cooperate in each step.
func (fsm *NdFsm) Survey(Path GameTrace) ([][3]float64, float64) {
	survey := make([][3]float64, len(Path))
	var totStates float64

	// init
	for _, s := range fsm.Node {
		s.lastSeen = 0
	}
	CurrentState := make(map[NodeIndex]float64)
	for _, s := range fsm.InitialState {
		CurrentState[s]++
	}

	// let's roll
	time := 0
	for i, step := range Path {
		time++

		// order-0 probs
		totStates = 0.0
		prob := new([3]float64)
		for s, m := range CurrentState {
			totStates += m
			prob[fsm.Node[s].Action] += m
			fsm.Node[s].lastSeen = time
		}
		for p := range prob {
			survey[i][p] = prob[p] / totStates
		}

		// next state
		NextState := make(map[NodeIndex]float64)
		for s, m := range CurrentState {
			next := fsm.Node[s].NextState(step)
			for _, n := range next {
				NextState[n] += m
			}
		}
		CurrentState = NextState
	}
	return survey, totStates
}
