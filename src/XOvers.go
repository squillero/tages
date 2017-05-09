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

var StaticXOver GeneticOperator = GeneticOperator{F: nil, W: 0, N: "XOver"}

func XOverAddIndividual(dst, src *Individual) {
	lat := make(map[NodeIndex]NodeIndex, len(dst.Node))

	src.UpdateInternals()

	// add nodes
	for i, n := range src.Node {
		if n.Active && n.reachable {
			ni := dst.AddNode(n.Action)
			dst.Node[ni].color = n.color
			lat[NodeIndex(i)] = ni
		}
	}

	// and link them
	for s, n := range src.Node {
		if n.Active && n.reachable {
			for _, d := range n.OnCooperation {
				dst.AddNdTransition(lat[NodeIndex(s)], COOPERATE, lat[d])
			}
			for _, d := range n.OnDefection {
				dst.AddNdTransition(lat[NodeIndex(s)], DEFECT, lat[d])
			}
		}
	}
	for _, t := range src.InitialState {
		dst.AddInitialState(lat[t])
	}
}
