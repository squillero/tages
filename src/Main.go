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
	"flag"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sort"
	"time"
)

// This useless comment makes goLint happy!
const (
	TagesVersion = "1.0.1"
)

// Payoff matrix
const (
	TEMPTATION = 5
	REWARD     = 3
	PUNISHMENT = 1
	SUCKER     = 0
)

// EA constants
const (
	MAX_GENERATIONS            = 300
	STEADY_STATE               = MAX_GENERATIONS / 2
	EXTINCTION_THRESHOLD       = STEADY_STATE * 3 / 4
	SEEK_LENGTH                = 10
	MU                         = 100
	LAMBDA                     = MU
	NU                         = MU * 5
	TOURNAMENT_SIZE            = 3
	INDIVIDUAL_INITIAL_DIM     = 5
	INDIVIDUAL_INITIAL_NON_DET = 0.50
	Θ                          = 0.01
	A_BIG_NUMBER               = 1e99
)

// Global payoff
var Payoff = PayoffMatrix{{-REWARD, -REWARD, -REWARD}, {0, PUNISHMENT, TEMPTATION}, {0, SUCKER, REWARD}}

// Dump best individuals as PNG
var OptionFirstMove = flag.String("1", "", "Force Tages' first move")

// Dump best individuals as PNG
var OptionDumpBest = flag.Bool("p", false, "Dump best individuals as PNG")

// Set a random seed, use 0 for current time
var OptionRandomSeed = flag.Int64("s", 42, "Random seed")

// Load single opponent from file
var OptionLoadFile = flag.String("f", "", "Load opponent")

// Show known strategies
var OptionList = flag.Bool("l", false, "Show known strategies")

// Plot graphs in color (graphwiz option)
var OptionColors = flag.Bool("c", false, "Plot graphs in color")

// Enable persistence for strategies, may help adaptive opponents.
// Notez bien: do not use it when running on a cluster!
var OptionPersistant = flag.Bool("P", false, "Use persistence")

// Boast a little bit
var OptionVerbose = flag.Bool("v", false, "Be slightly more verbose than usual")

func main() {
	log.SetFlags(0)
	log.Println("Tages v" + TagesVersion + " - Yet another adaptive EC-based player for the IPD")
	log.Println("(!) between 2014 and 2015 by the Cyber Lasas")
	log.Println("This is free software, and you are welcome to redistribute it under certain conditions")
	log.Println("")
	log.SetFlags(log.Lmicroseconds)

	flag.Parse()

	runtime.GOMAXPROCS(2)
	log.Printf("Running %d parallel processes on %d available processors\n", runtime.GOMAXPROCS(0), runtime.NumCPU())

	s := *OptionRandomSeed
	if s <= 0 {
		s = time.Now().UTC().UnixNano()
	}
	rand.Seed(s)

	log.Println("Random seed: ", s)
	log.Println("Payoff matrix:", &Payoff)

	Opponents := make(map[string]Player)
	OpponentName := ""
	var OppList1, OppList2 []Player

	// ND-FSMs
	if *OptionLoadFile != "" {
		log.Println("Reading nd-fsm opponent from file \"" + *OptionLoadFile + "\"")
		o := LoadNdFsm(*OptionLoadFile)
		ft := "FSM"
		for _, n := range o.Node {
			if len(n.OnCooperation)+len(n.OnDefection) > 2 {
				ft = "ND-FSM"
			}
		}
		o.Name = ft + "# " + o.Name
		Opponents[o.GetName()] = o
		OpponentName = o.GetName()
	} else {
		log.Println("Reading nd-fsm opponents from directory \"Opponents\"")
		for _, o := range LoadNdFsms("Opponents") {
			ft := "FSM"
			for _, n := range o.Node {
				if len(n.OnCooperation)+len(n.OnDefection) > 2 {
					ft = "ND-FSM"
				}
			}
			o.Name = ft + "# " + o.Name
			Opponents[o.GetName()] = o
			t1, t2 := o.Duplicate(), o.Duplicate()
			OppList1 = append(OppList1, &t1)
			OppList2 = append(OppList2, &t2)
		}
	}

	// NON-ND-FSMs
	for _, p := range StdStrategies() {
		Opponents[p.GetName()] = p
		OppList1 = append(OppList1, p)
	}
	for _, p := range StdStrategies() {
		OppList2 = append(OppList2, p)
	}

	// RANDOM-FSMs
	for _, p := range RandomStrategies() {
		Opponents[p.GetName()] = p
		OppList1 = append(OppList1, p)
	}
	for _, p := range RandomStrategies() {
		OppList2 = append(OppList2, p)
	}

	// SPICY TFTs
	for _, p := range SpicyTfTs() {
		Opponents[p.GetName()] = p
		OppList1 = append(OppList1, p)
	}
	for _, p := range SpicyTfTs() {
		OppList2 = append(OppList2, p)
	}

	// Standard i/o
	var sio TCIAIG_TestStrategy
	Opponents[sio.GetName()] = &sio

	// selected opponent
	if *OptionList {
		log.Println("Available Opponents:")
		var ol []string
		for _, o := range Opponents {
			ol = append(ol, o.GetName())
		}
		sort.Strings(ol)
		for _, o := range ol {
			fmt.Println(o)
		}
	} else if len(flag.Args()) == 1 {
		OpponentName = flag.Args()[0]
		if _, ok := Opponents[OpponentName]; !ok {
			log.Panicf("Can't play against \"%s\"\n", OpponentName)
		}
	}

	glen := []int{168, 359, 306, 622, 319}
	//glenda := []int{16, 35}
	//glenda := []int{3, 5}

	if OpponentName != "" {
		log.Printf("Playing against \"%s\"\n", Opponents[OpponentName].GetName())
		Tusna(Opponents[OpponentName], glen)
	}
}
