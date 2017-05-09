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
	"math/rand"
)

func RandomStrategies() []Player {
	players := make([]Player, 0)

	for t := 0; t <= 10; t++ {
		p := new(RandomStrategy)
		p.SetAlpha(float64(t) / 10.0)
		players = append(players, p)
	}

	return players
}

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// Cooperate/defect with a probability equals to opponent's c/d frequency
//////////////////////////////////////////////////////////////////////////////

type RandomStrategy struct {
	α float64
}

func (r *RandomStrategy) SetAlpha(α float64) {
	r.α = α
}

func (r *RandomStrategy) GetName() string {
	return fmt.Sprintf("RND# RandomStrategy [α=%0.2f]", r.α)
}

func (r *RandomStrategy) ply() PPly {
	if rand.Float64() < r.α {
		return DEFECT
	} else {
		return COOPERATE
	}
}

func (r *RandomStrategy) FirstPly() PPly {
	return r.ply()
}

func (r *RandomStrategy) RePly(_ PPly) PPly {
	return r.ply()
}
