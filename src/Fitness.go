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
	"math"
	"math/rand"
	"sort"
)

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// Types, constants Interfaces
//////////////////////////////////////////////////////////////////////////////

var ChromaStats [4]float64

const (
	ξ                 = 1.0001
	STATES_HARD_LIMIT = 100
	HARD
)

type Fitness struct {
	Valid           bool
	Usable          bool
	Likelyhood      float64
	Plausibility    float64
	Competitiveness float64
	Compactness     float64
}

func (f *Fitness) Check() bool {
	if !f.Valid {
		log.Panicln("Invalid fitness: ", f)
	}
	if f.Competitiveness < 0 || math.IsNaN(f.Competitiveness) || math.IsInf(f.Competitiveness, +1) {
		log.Panicln("Illegal competitiveness: ", *f)
	}
	if f.Likelyhood > 0 {
		log.Panicln("Illegal Likelyhood: ", *f)
	}
	if f.Plausibility < 0 || math.IsNaN(f.Plausibility) || math.IsInf(f.Plausibility, +1) {
		log.Panicln("Illegal plausibility: ", *f)
	}
	if f.Compactness < 0 || math.IsNaN(f.Compactness) {
		log.Panicln("Illegal compactness: ", *f)
	}
	return true
}

func (f *Fitness) String() string {
	var Usable string
	if f.Usable {
		Usable = "+"
	} else {
		Usable = "-"
	}
	if f.Valid {
		return Usable + fmt.Sprintf("P=%.2g/L=%.2g:C=%g:X=%g", f.Plausibility, f.Likelyhood, f.Competitiveness, f.Compactness)
	} else {
		return "?P-/C-:$-:<-"
	}
}

//////////////////////////////////////////////////////////////////////////////
// Compare two individuals, reflect fitness comparison
// compare f1 vs. f2. returns +1 if >, 0 if ==, -1 if <
func (i *Individual) CompareTo(o *Individual) bool {
	if FitCompare(i.Fit, o.Fit) > 0 {
		return true
	} else {
		return false
	}
}

//////////////////////////////////////////////////////////////////////////////
// population sort
//////////////////////////////////////////////////////////////////////////////

func (s *Population) Len() int {
	return len(s.Individual)
}

func (s *Population) Swap(i, j int) {
	s.Individual[i], s.Individual[j] = s.Individual[j], s.Individual[i]
}

type ByFitness struct {
	*Population
}

func (s ByFitness) Less(j, i int) bool {
	if FitCompare(s.Population.Individual[i].Fit, s.Population.Individual[j].Fit) <= 0 {
		return true
	} else {
		return false
	}
}

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// compare fitnesses
//////////////////////////////////////////////////////////////////////////////

func FitCompare(fit1, fit2 Fitness) int {
	fit1.Check()
	fit2.Check()

	// if only one is usable
	if fit1.Usable && !fit2.Usable {
		return +1
	} else if !fit1.Usable && fit2.Usable {
		return -1
	}
	// patch Likelyhood
	minLikelyhood := 0.0
	if fit1.Likelyhood > fit2.Likelyhood {
		minLikelyhood = fit1.Likelyhood
	} else {
		minLikelyhood = fit2.Likelyhood
	}

	// let's go chromatic
	ind1 := []float64{fit1.Plausibility, fit1.Likelyhood - minLikelyhood, fit1.Competitiveness, -fit1.Compactness}
	ind2 := []float64{fit2.Plausibility, fit2.Likelyhood - minLikelyhood, fit2.Competitiveness, -fit2.Compactness}
	c := ChromaticCompare(ind1, ind2)

	//}
	return c
}

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// fitness evaluators
//////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////////////////////////
// Evaluate the fitness of a whole population
func (p *Population) Evaluate() {
	//p.Champion.Evaluate(p.Opponents, p.Turns)

	for i := range p.Individual {
		if p.Individual[i].Fit.Valid == false {
			// Update reachability & create inverses
			p.Individual[i].UpdateInternals()
			p.Individual[i].CreateInverse()

			// evaluate plausibility & usability
			p.Individual[i].EvaluatePlausibility(p.Paths, p.Traces)
			// evaluate Likelyhood
			p.Individual[i].EvaluateLikelyhood(p.Paths, p.Traces)
			// evaluate compactness
			p.Individual[i].EvaluateCompactness(p.Paths, p.Traces)
			// evaluate competitiveness
			p.Individual[i].EvaluateCompetitiveness(p.Opponents, p.Turns)

			// let's penalize bloating elements
			if len(p.Individual[i].ActiveNodes()) > p.ζ {
				p.Individual[i].Fit.Usable = false
				p.Individual[i].Fit.Likelyhood = 0
			}
			if len(p.Individual[i].ActiveNodes()) > int(1.5*float64(p.ζ)) {
				p.Individual[i].Fit.Competitiveness = 0
				p.Individual[i].Fit.Plausibility = 0
			}

			// that's all
			p.Individual[i].Fit.Valid = true

		}
	}
	sort.Sort(ByFitness{p})
}

func scale(trace []GameTrace) float64 {
	tl := 0.0
	for _, t := range trace {
		tl += float64(len(t))
	}
	return tl
}

//////////////////////////////////////////////////////////////////////////////
// Evaluate single fitness components

func (ind *Individual) EvaluatePlausibility(Paths, Footprints []GameTrace) {
	if len(Paths) == 0 {
		ind.Fit.Plausibility = 1
		return
	}
	ind.Fit.Plausibility = 0.0
	k := 1.0
	ind.Fit.Usable = true
	for i := range Paths {
		plausibility, valid := ind.EvaluatePlausibility_Single(Paths[i], Footprints[i])
		ind.Fit.Plausibility += k * plausibility
		if !valid {
			ind.Fit.Usable = false
		}
		k *= ξ
	}
}

// Plausibility ∈ [0, 1]
func (ind *Individual) EvaluatePlausibility_Single(Paths, Footprints GameTrace) (float64, bool) {
	if len(Paths) == 0 {
		return 1, true
	}

	// select machine
	lastState := Footprints[len(Footprints)-1]
	// drop last bit
	ftp := Footprints[0:len(Footprints)]
	pat := Paths[0 : len(Paths)-1]
	// invert
	revFtp := ftp.Reverse()
	revPat := pat.Reverse()
	l, s := ind.Inverse[lastState].Scout(revPat, revFtp)

	FullyValid := false
	for _, i := range ind.InitialState {
		if s[i] > 0 {
			// gotcha! we traced back the fsm up to an initial state
			if !FullyValid {
				l++
			}
			FullyValid = true
		}
	}

	return float64(l) / float64(len(Paths)), FullyValid
}

// Likelyhood ∈ [0, +Inf]
func (ind *Individual) EvaluateLikelyhood(Paths, Footprints []GameTrace) {
	Likelyhood := 0.0
	k := 1.0

	for i := range Paths {
		evidence, _ := ind.Survey(Paths[i])
		for e := 0; e < len(Paths[i])-1; e++ {
			Likelyhood += k * math.Log(evidence[e][Footprints[i][e]])
		}
		k *= ξ
	}
	ind.Fit.Likelyhood = Likelyhood
}

// Compactness ∈ [0, +Inf]
func (ind *Individual) EvaluateCompactness(Paths, Footprints []GameTrace) {
	Compactness := 0.0
	if len(Paths) == 0 || len(Paths[0]) == 0 {
		Compactness = float64(len(ind.ActiveNodes()))
		ind.Fit.Compactness = Compactness
		return
	}
	if len(ind.Node) > STATES_HARD_LIMIT {
		Compactness = math.Inf(+1)
		ind.Fit.Compactness = Compactness
		return
	}
	lastSeen := make([]int, len(ind.Node))

	pathLength := 0
	totStates := 0.0
	for i := range Paths {
		if len(Paths[i]) > 0 {
			_, t := ind.Survey(Paths[i])
			totStates += t - 1
			for j := range lastSeen {
				if ind.Node[j].lastSeen > 0 {
					lastSeen[j] = ind.Node[j].lastSeen + pathLength
				}
			}
			pathLength += len(Paths[i])
		}
	}

	Threshold_RECENT := pathLength * 4 / 5
	Threshold_AVERAGE := pathLength / 2

	for _, a := range ind.ActiveNodes() {
		//ind.Node[a].color = float64(lastSeen[a]) / float64(pathLength)
		if lastSeen[a] >= Threshold_RECENT {
			ind.Node[a].color = "lightskyblue"
			Compactness += 1
		} else if lastSeen[a] >= Threshold_AVERAGE {
			ind.Node[a].color = "lightgoldenrod1"
			Compactness += 5
		} else {
			ind.Node[a].color = "orangered"
			Compactness += 10
		}
		//ind.Node[a].color = fmt.Sprintf("gray%d", int(90-50*float64(ind.Node[a].lastSeen)/float64(pathLength)))
		if !ind.Node[a].reachable {
			ind.Node[a].color = "silver"
			Compactness += 20
		}
	}
	Compactness *= math.Exp(float64(len(ind.ActiveNodes())))
	Compactness *= 1 + 100*totStates
	ind.Fit.Compactness = Compactness
}

// Competitiveness ∈ [0, +∞]
func (i *Individual) EvaluateCompetitiveness(opponents []Player, turns []int) {
	var totPayoff float64
	for t := 0; t < len(opponents); t++ {
		payoff := 0
		for _, tt := range turns {
			r1, _ := Match(i, opponents[t], tt)
			payoff += r1
		}
		totPayoff += float64(payoff)
	}
	i.Fit.Competitiveness = math.Max(0, totPayoff)
}

// min & max is missing from go?
func min(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}
func max(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func ExperimentalChromaticCompare(ind1, ind2 []float64) int {
	chroma := make([]float64, len(ind1))
	rainbow := 0.0
	for i := range ind1 {
		chroma[i] = max(ind1[i], ind2[i]) - min(ind1[i], ind2[i])
		rainbow += chroma[i]
	}
	if rainbow == 0 {
		return 0
	}
	rainbow *= rand.Float64()
	c := -1
	for rainbow >= 0 {
		c++
		rainbow -= chroma[c]
	}
	if c < 0 || c >= len(ChromaStats) {
		log.Panicln("PANIC: CC overflowed", c, chroma)
	}
	ChromaStats[c]++
	//log.Println(ChromaStats, "#", chroma)
	if ind1[c] > ind2[c] {
		//log.Println("Checking ", c, ind1[c], "vs.", ind2[c], "1 is better!")
		return +1
	} else if ind1[c] < ind2[c] {
		//log.Println("Checking ", c, ind1[c], "vs.", ind2[c], "2 is better!")
		return -1
	} else {
		return 0
	}
}

func valid(v float64) bool {
	if math.IsNaN(v) {
		return false
	}
	if math.IsInf(v, +1) {
		return false
	}
	if math.IsInf(v, -1) {
		return false
	}
	return true
}

func ChromaticCompare(ind1, ind2 []float64) int {
	chroma := make([]float64, len(ind1))
	rainbow := 0.0

	for i := range ind1 {
		if ind1[i] == 0 && ind2[i] == 0 {
			chroma[i] = 0
		} else if !valid(ind1[i]) || !valid(ind2[i]) {
			chroma[i] = math.Inf(+1)
		} else {
			delta := math.Max(ind1[i], ind2[i]) - math.Min(ind1[i], ind2[i])
			baseline := math.Max(math.Abs(ind1[i]), math.Abs(ind2[i]))
			chroma[i] = delta / baseline
		}
		rainbow += chroma[i]
	}
	var color int
	if valid(rainbow) {
		if rainbow == 0 {
			return 0
		}
		rainbow *= rand.Float64()
		color = -1
		for rainbow >= 0 {
			color++
			rainbow -= chroma[color]
		}
	} else {
		color = 0
		for valid(chroma[color]) {
			color++
		}
	}

	// paranoia check!
	if color < 0 || color >= len(ChromaStats) {
		log.Panicln("PANIC: CC overflowed", color, chroma)
	}

	ChromaStats[color]++
	if ind1[color] > ind2[color] {
		return +1
	} else if ind1[color] < ind2[color] {
		return -1
	} else {
		return 0
	}
}
