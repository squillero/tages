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
	"math/rand"
)

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// Learn the expected reward for cooperation/defection
// (Japanese: 無; Korean: 무) or Chinese wu (traditional Chinese: 無; simplified Chinese: 无)
//////////////////////////////////////////////////////////////////////////////

type VanillaRL struct {
	LastMove                           PPly
	rewardCooperation, rewardDefection float64
}

func (_ *VanillaRL) GetName() string {
	return "Vanilla RL"
}

func (p *VanillaRL) FirstPly() PPly {
	p.rewardCooperation, p.rewardDefection = 0.1, 0.1

	return p.move()
}

func (p *VanillaRL) RePly(oppPly PPly) PPly {
	if p.LastMove == COOPERATE {
		p.rewardCooperation += float64(Payoff[p.LastMove][oppPly])
	} else {
		p.rewardDefection += float64(Payoff[p.LastMove][oppPly])
	}

	return p.move()
}

func (p *VanillaRL) move() PPly {
	r := rand.Float64() * (p.rewardDefection + p.rewardCooperation)
	if r < p.rewardDefection {
		p.LastMove = DEFECT
	} else {
		p.LastMove = COOPERATE
	}

	return p.LastMove
}

/*
//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// Cooperate/defect with a probability equals to opponent's c/d frequency
//////////////////////////////////////////////////////////////////////////////

type FairStrategy struct {
	History [3]int
}

func (_ *FairStrategy) GetName() string {
	return "Fair Bob"
}

func (p *FairStrategy) FirstPly() PPly {
	p.History[DEFECT] = 1
	p.History[COOPERATE] = 1
	if rand.Intn(p.History[DEFECT]+p.History[COOPERATE]) < p.History[DEFECT] {
		return DEFECT
	} else {
		return COOPERATE
	}
}

func (p *FairStrategy) RePly(oppPly PPly) PPly {
	p.History[oppPly]++
	if rand.Intn(p.History[DEFECT]+p.History[COOPERATE]) < p.History[DEFECT] {
		return DEFECT
	} else {
		return COOPERATE
	}
}
*/

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// Unreliable TFT
//////////////////////////////////////////////////////////////////////////////

type UnreliableSparring struct {
}

func (_ *UnreliableSparring) GetName() string {
	return "Lunatic Jack (90% Tit for Tat)"
}

func (p *UnreliableSparring) FirstPly() PPly {
	if rand.Float64() < 0.5 {
		return COOPERATE
	} else {
		return DEFECT
	}
}

func (p *UnreliableSparring) RePly(oppPly PPly) PPly {
	if rand.Float64() < 0.90 {
		return oppPly
	} else {
		return OppositeMove(oppPly)
	}
}
