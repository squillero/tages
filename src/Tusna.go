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
	"math"
	"math/rand"
	"time"
)

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// Tusna (ie. the sacred swan of Tages)
//////////////////////////////////////////////////////////////////////////////

var TCIAIG_DIRTY_HACK PPly

func Tusna(opponent Player, gameLengths []int) {
	sparring_partner := make([]Player, 0)
	sparring_partner = append(sparring_partner, new(FairStrategy))
	sparring_partner = append(sparring_partner, new(UnreliableSparring))
	//sparring_partner = append(sparring_partner, new(VanillaRL))
	//for _, f := range SpicyTfTs() {
	//	sparring_partner = append(sparring_partner, f)
	//}

	if len(sparring_partner) == 1 {
		log.Println("Using only one sparring partner")
	} else {
		log.Println("Using", len(sparring_partner), "sparring partners")
	}
	for _, o := range sparring_partner {
		log.Printf("++ %s\n", o.GetName())
	}

	// setup training lengths
	trainingLengths := []int{100, 200, 300, 400, 500, 600, 700, 800, 900}
	//trainingLengths := []int{100, 300, 500, 700, 900}
	//trainingLengths := []int{50, 100, 150, 200, 250, 300, 350, 400, 450, 500, 550, 600, 650, 50, 100, 150, 200, 250, 300, 350, 400, 450, 500, 550, 600, 650}
	//trainingLengths := []int{84, 153, 158, 169, 178, 217, 318, 432, 613, 718, 84, 153, 158, 169, 178, 217, 318, 432, 613, 718}
	//trainingLengths := []int{84, 153, 158, 169, 178, 217, 318, 432, 613, 718}
	log.Printf("Game length: %v\n", gameLengths)
	log.Printf("Training game length: %v\n", trainingLengths)

	// boot := []PPly{DEFECT, DEFECT, COOPERATE}
	// log.Println(FindBestPath(boot, 3, sparring_partner[0], Payoff))
	// os.Exit(0)

	//////////////////////////////////////////////////////////////////////////
	// Create new population

	population := CreatePopulation(MU, NU, LAMBDA, INDIVIDUAL_INITIAL_DIM, INDIVIDUAL_INITIAL_NON_DET, sparring_partner, trainingLengths)
	population.MUtations.AddOperator(MUtationChangeInitialState, "ChgIS")
	population.MUtations.AddOperator(MUtationAddInitialState, "AddIS")
	population.MUtations.AddOperator(MUtationRemoveInitialState, "RmvIS")
	population.MUtations.AddOperator(MUtationAddNode, "AddND")
	population.MUtations.AddOperator(MUtationRemoveNode, "RmvND")
	population.MUtations.AddOperator(MUtationChangeTransition, "ChgTRAN")
	population.MUtations.AddOperator(MUtationChangeStateAction, "ChgACT")
	population.MUtations.AddOperator(MUtationAddNdTransition, "AddNdTRAN")
	population.MUtations.AddOperator(MUtationRemoveNdTransition, "RmvNdTRAN")
	log.Println(&population)

	// ReconciliatoryFactor
	ReconciliatoryFactor := 0.0
	ReconciliationSteps := 0
	ReconciliatoryAttempts := 0

	prophecyTry, prophecySkip, prophecyFail, prophecySuccess := 0, 0, 0, 0
	for game, gameLen := range gameLengths {
		tInit := time.Now()

		poTages, poOpponent := 0, 0

		population.Paths = append(population.Paths, GameTrace{})
		population.Traces = append(population.Traces, GameTrace{})
		for turn := 0; turn < gameLen; turn++ {
			log.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
			log.Printf("%% GAME #%v - TURN #%v\n", game+1, turn+1)
			log.Printf("Path : %v\n", population.Paths)
			log.Printf("Trace: %v\n", population.Traces)

			var moveTages, moveOpponent, forecast PPly

			// let's roll

			// turan's own ply && forested ply
			if game == 0 && len(population.Paths[0]) == 0 {
				if *OptionFirstMove != "" {
					if *OptionFirstMove == "COOPERATE" || *OptionFirstMove == "C" {
						moveTages = COOPERATE
					} else if *OptionFirstMove == "DEFECT" || *OptionFirstMove == "D" {
						moveTages = DEFECT
					} else {
						log.Fatalf("Can't parse move \"%s\"", *OptionFirstMove)
					}
				} else {
					moveTages = population.VanillaRL()
				}
			} else {
				moveTages = SmartMove(&population)
			}

			if ReconciliationSteps > 0 {
				log.Println("Reconciliatory step (", ReconciliationSteps, "to go)")
				moveTages = COOPERATE
				ReconciliationSteps--
			} else if rand.Float64() < ReconciliatoryFactor {
				ReconciliatoryAttempts++
				ReconciliationSteps = 4 + rand.Int()%7
				log.Printf("Starting a reconciliation attempt (num=%d, length=%d)...", ReconciliatoryAttempts, ReconciliationSteps)
				moveTages = COOPERATE
				ReconciliationSteps--
			}

			if population.Individual[0].Fit.Usable {
				population.Paths[game] = append(population.Paths[game], moveTages)
				prob, _ := population.Individual[0].Survey(population.Paths[game])
				population.Paths[game] = population.Paths[game][:len(population.Paths[game])-1]
				if prob[len(prob)-1][COOPERATE] >= .5 {
					forecast = COOPERATE
				} else {
					forecast = DEFECT
				}
			} else {
				forecast = UNSPECIFIED
			}

			// opponent real ply
			TCIAIG_DIRTY_HACK = moveTages
			if turn == 0 {
				moveOpponent = opponent.FirstPly()
			} else {
				path := population.Paths[game]
				moveOpponent = opponent.RePly(path[len(path)-1])
			}

			// RECONCILIATION?
			var moveTagesLast PPly
			if turn == 0 {
				moveTagesLast = UNSPECIFIED
			} else {
				moveTagesLast = population.Paths[game][len(population.Paths[game])-1]
			}
			if moveTagesLast == DEFECT && moveOpponent == DEFECT {
				ReconciliatoryFactor += math.Pow(0.1, float64(ReconciliatoryAttempts+1))
			} else {
				ReconciliatoryFactor = 0
			}

			log.Printf("Tages: %v #vs# %v: %v\n", moveTages, opponent.GetName(), moveOpponent)
			p := ""
			if forecast != UNSPECIFIED {
				prophecyTry++
				if forecast == moveOpponent {
					prophecySuccess++
					p = "Opponent's move was correctly predicted"
				} else {
					prophecyFail++
					p = "Opponent's move was not correctly predicted"
				}
			} else {
				prophecySkip++
				p = "Opponent's move was not predicted"
			}
			population.Calculateζ(prophecySuccess, prophecyTry+prophecySkip)
			log.Printf("%s. Prophecy rate %.2f%%, accuracy %.2f%% (ζ = %v)\n", p, 100*float64(prophecyTry)/float64(prophecyTry+prophecySkip), 100*float64(prophecySuccess)/float64(prophecySuccess+prophecyFail), population.ζ)
			log.Printf("Reconciliation factor ρ = %v\n", ReconciliatoryFactor)

			poOpponent += Payoff[moveOpponent][moveTages]
			poTages += Payoff[moveTages][moveOpponent]
			log.Printf("Reward: %v #vs# %v (last); %v #vs# %v (tot); %.3v #vs# %.3v (avg)\n", Payoff[moveTages][moveOpponent], Payoff[moveOpponent][moveTages], poTages, poOpponent, float64(poTages)/float64(turn+1), float64(poOpponent)/float64(turn+1))

			population.AddStep(game, moveTages, moveOpponent)
			//population.Paths = append(population.Paths, moveTages)
			//population.Traces = append(population.Traces, moveOpponent)

			log.Printf("Elapsed: %v (global)\n", time.Now().Sub(tInit))
		}
		// nice stats
		log.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
		log.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%")
		log.Printf("%% GAME #%v (%v turns)\n", game+1, gameLen)
		s1, s2 := "Tages v"+TagesVersion, opponent.GetName()
		for len(s1) != len(s2) {
			if len(s1) < len(s2) {
				s1 += " "
			} else {
				s2 += " "
			}
		}
		log.Printf("%v: %v\n", s1, population.Paths[len(population.Paths)-1])
		log.Printf("%v: %v\n", s2, population.Traces[len(population.Paths)-1])
		log.Printf("Reward: %v #vs# %v (tot); %.3g #vs# %.3g (avg)\n", poTages, poOpponent, float64(poTages)/float64(gameLen), float64(poOpponent)/float64(gameLen))
		Scoresheet("Tages v"+TagesVersion, opponent.GetName(), gameLen, poTages, poOpponent, population.Paths[game], population.Traces[game])
		name := fmt.Sprintf("Model for %s (game %d)", opponent.GetName(), game+1)
		ioutil.WriteFile(name+".gv", population.Individual[0].GvEncode(), 0644)
		ioutil.WriteFile(name+".json", population.Individual[0].GvEncode(), 0644)
		//ioutil.WriteFile(name+".oy.txt", population.Individual[0].GvEncode(), 0644)
		Gv2Png(name)
	}
}

var PreviousModel NdFsm

func SmartMove(p *Population) PPly {
	log.Println("Starting model optimization")
	tInit := time.Now()

	p.Check()
	p.Invalidate()
	p.Evaluate()
	//p.Dump()
	currentBest := p.Individual[0].Fit
	log.Printf("... (-) %v%s\n", &p.Individual[0], p.Tagζ(0))
	evolved := ""
	g := 0

	var lastImprovement int
	for evolved == "" {
		// evolutionary canon
		p.Begat()
		p.Evaluate()
		p.Slaughter()
		p.UpdateOperators()

		if FitCompare(p.Individual[0].Fit, currentBest) > 0 {
			// fucking awesome
			//p.Individual[0].DumpPNG()
			//p.Individual[0].Inverse[COOPERATE].DumpPNG()
			//p.Individual[0].Inverse[DEFECT].DumpPNG()
			currentBest = p.Individual[0].Fit
			if *OptionVerbose {
				log.Printf("... (%v) %v%s\n", g, &p.Individual[0], p.Tagζ(0))
			}
			lastImprovement = g
		}
		g++

		if g >= MAX_GENERATIONS {
			// that's all folks
			if !*OptionVerbose {
				log.Printf("... (%v) %v%s\n", g, &p.Individual[0], p.Tagζ(0))
			}
			evolved = fmt.Sprintf("Reached max generation (%d)", MAX_GENERATIONS)
		}
		if g-lastImprovement >= EXTINCTION_THRESHOLD && !p.Individual[0].Fit.Usable {
			if !*OptionVerbose {
				log.Printf("... (%v) %v%s\n", g, &p.Individual[0], p.Tagζ(0))
			}
			log.Printf("... (%v) Starting an extinction phase\n", g)
			p.Extinguish()
			p.Evaluate()
			p.Slaughter()
			p.UpdateOperators()
			p.Check()
			lastImprovement = g
		} else if g-lastImprovement >= STEADY_STATE && p.Individual[0].Fit.Usable {
			// gee. it's getting *really* boring out there
			if !*OptionVerbose {
				log.Printf("... (%v) %v%s\n", g, &p.Individual[0], p.Tagζ(0))
			}
			evolved = fmt.Sprintf("Reached steady state (%d)", g)
			//p.Dump()
		}
	}

	Model := p.Individual[0].Duplicate()
	if *OptionDumpBest {
		if !FsmEqual(&Model.NdFsm, &PreviousModel) {
			Model.DumpPNG()
			PreviousModel = Model.NdFsm.Duplicate()
		}
	}
	log.Printf("%v. Best model: \"%v\"\n", evolved, Model.Name)
	//log.Println(&Model.NdFsm)
	log.Printf("Current chroma: %v\n", ChromaStats)
	var move PPly
	if !p.Individual[0].Fit.Usable {
		log.Println("WARNING:: The model is not reliable")
		move = p.VanillaRL()
	} else {
		// we got a "reasonable" Model of the opponent. let's exploit it
		game := len(p.Traces) - 1
		ftp := p.Traces[game]
		ftp = append(ftp, UNSPECIFIED)
		_, lastState := Model.Scout(p.Paths[game], ftp)
		n := 0
		last := false
		var pat GameTrace
		best_pat := make(GameTrace, SEEK_LENGTH)
		worst_val, best_val := -1.0, -1.0
		for !last {
			pat, last = Enumerate(SEEK_LENGTH, n)
			val := Model.ExpectedPayoff(lastState, pat)
			if val > best_val {
				best_val = val
				copy(best_pat, pat)
			}
			if worst_val < 0 || val < worst_val {
				worst_val = val
			}
			n++
		}

		log.Printf("Starting from %v, analyzed %v possible paths against %v; rewards between %.2g and %.2g\n", lastState, n, Model.Name, worst_val, best_val)
		log.Printf("Best path: %v\n", best_pat)
		move = best_pat[0]

		if rand.Float64() < Θ {
			log.Println("Using Θ-exploration")
			if move == COOPERATE {
				move = DEFECT
			} else {
				move = COOPERATE
			}
		}
	}
	log.Printf("Elapsed: %v (model optimization)\n", time.Now().Sub(tInit))

	return move
}
