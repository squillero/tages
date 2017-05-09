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

const (
	EXPLORATION = .1
	LEN         = 4
)

type PseudoRL struct {
	pastPayoff map[string]int
	lastMove   PPly
	history    string
}

func (_ *PseudoRL) GetName() string {
	return fmt.Sprintf("PseudoRL[%d]", LEN)
}

func (p *PseudoRL) FirstPly() PPly {
	p.pastPayoff = make(map[string]int)
	p.history = ""

	return p.doMove(COOPERATE)
}

func (p *PseudoRL) RePly(oppPly PPly) PPly {
	var po int
	var m PPly
	if p.lastMove == COOPERATE && oppPly == COOPERATE {
		po = 3
	} else if p.lastMove == COOPERATE && oppPly == DEFECT {
		po = 0
	} else if p.lastMove == DEFECT && oppPly == COOPERATE {
		po = 5
	} else if p.lastMove == DEFECT && oppPly == DEFECT {
		po = 1
	}
	p.pastPayoff[p.history] += po
	if len(p.history) > LEN {
		log.Fatal(p)
	}
	// log.Println(p.pastPayoff)

	if len(p.history) > LEN-1 {
		p.history = p.history[1:LEN]
	}
	if rand.Float32() < EXPLORATION {
		m = PPly(rand.Intn(2))
	} else {
		r := rand.Intn(p.pastPayoff[p.history+_D] + p.pastPayoff[p.history+_C] + 1)
		// log.Println(r, p.history+_D_, p.pastPayoff[p.history+_D_], p.history+_C_, p.pastPayoff[p.history+_C_])
		if r < p.pastPayoff[p.history+_D] {
			m = DEFECT
		} else {
			m = COOPERATE
		}
	}
	return p.doMove(m)
}

func (p *PseudoRL) doMove(m PPly) PPly {
	p.history += m.Mark()
	p.lastMove = m

	// log.Println(p.history)
	return m
}

func (p *PseudoRL) Describe() {
	log.Println(p)
}
