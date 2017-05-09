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
)

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// Read moves from TCIAIG and ouput to STDOUT
// Required by TCIAIG reviewers
//////////////////////////////////////////////////////////////////////////////

type TCIAIG_TestStrategy struct {
	n int
}

func (_ *TCIAIG_TestStrategy) GetName() string {
	return "TCIAIG_TestStrategy"
}

func (p *TCIAIG_TestStrategy) FirstPly() PPly {
	return p.ReadFromStdin()
}

func (p *TCIAIG_TestStrategy) RePly(oppPly PPly) PPly {
	return p.ReadFromStdin()
}

func (p *TCIAIG_TestStrategy) ReadFromStdin() PPly {
	var in string
	fmt.Println(TCIAIG_DIRTY_HACK.String())
	fmt.Scan(&in)
	if in == "COOPERATE" || in == "C" {
		return COOPERATE
	}
	if in == "DEFECT" || in == "D" {
		return DEFECT
	}
	log.Fatalf("Can't parse move \"%s\"", in)
	return UNSPECIFIED
}
