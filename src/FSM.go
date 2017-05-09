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

//////////////////////////////////////////////////////////////////////////////
// Basic structures (just to get some nice static checking)

type NodeIndex int
type NodeIndexList []NodeIndex

// sort interface
func (s NodeIndexList) Len() int           { return len(s) }
func (s NodeIndexList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s NodeIndexList) Less(i, j int) bool { return s[i] < s[j] }

type FSM interface {
	// pretty obvious stuff
	AddNode(action PPly) NodeIndex
	SetTransition(nodeCurrent NodeIndex, action PPly, nodeTarget NodeIndex) bool
	// probability [0, 1] that the recorded footprints were generated by the FSM on the given path (opponent moves)
	Likelihood(path, footprints GameTrace) float64
}
