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
	"os"
	"sort"
	"time"
)

//////////////////////////////////////////////////////////////////////////////

// Types & Interfaces: GameTrace & Player
const (
	UNSPECIFIED = -1
	INVALID     = 0
	DEFECT      = 1
	COOPERATE   = 2
	_U          = "?"
	_I          = "!"
	_D          = "."
	_C          = "*"
	_UU         = "?"
	_II         = "!"
	_DD         = "D"
	_CC         = "C"
	_UUU        = "UNSPECIFIED"
	_III        = "INVALID"
	_DDD        = "DEFECT"
	_CCC        = "COOPERATE"
)

type PPly int
type GameTrace []PPly
type PayoffMatrix [3][3]int

func (g GameTrace) Reverse() GameTrace {
	r := make(GameTrace, len(g))
	for i, s := range g {
		r[len(g)-1-i] = s
	}
	return r
}

type Player interface {
	GetName() string
	RePly(oppPly PPly) PPly
	FirstPly() PPly
}

func (pm *PayoffMatrix) String() string {
	return fmt.Sprintf("R:%d/S:%d/T:%d/P:%d", pm[COOPERATE][COOPERATE], pm[COOPERATE][DEFECT], pm[DEFECT][COOPERATE], pm[DEFECT][DEFECT])
}

func (p PPly) String() string {
	return [...]string{_UUU, _III, _DDD, _CCC}[p+1]
}
func (p PPly) Mark() string {
	return [...]string{_UU, _II, _DD, _CC}[p+1]
}

func (gt GameTrace) String() string {
	str := ""
	for _, s := range gt {
		str += s.Mark()
	}
	return str
}

func (p PPly) Check() bool {
	if p != DEFECT && p != COOPERATE && p != INVALID && p != UNSPECIFIED {
		log.Panicf("Invalid ply: %d\n", p)
	}
	return true
}

func (p PPly) Check2() bool {
	if p != DEFECT && p != COOPERATE {
		log.Panicf("Invalid ply: %d\n", p)
	}
	return true
}

func (fp1 GameTrace) Compatible(fp2 GameTrace) bool {
	if len(fp1) != len(fp2) {
		return false
	}
	for t := range fp1 {
		if fp1[t] != fp2[t] {
			return false
		}
	}
	return true
}

func OppositeMove(p PPly) PPly {
	if p == COOPERATE {
		return DEFECT
	} else if p == DEFECT {
		return COOPERATE
	} else {
		return UNSPECIFIED
	}
}

//////////////////////////////////////////////////////////////////////////////
// The original Axlrod's game length

func GameLength() int {
	len := 1
	for rand.Float64() < 0.99654 {
		len++
	}
	return len
}

//////////////////////////////////////////////////////////////////////////////
// Matches

func Match(p1, p2 Player, turns int) (int, int) {
	r1, r2 := 0, 0
	m1, m2 := p1.FirstPly(), p2.FirstPly()

	for t := 0; t < turns; t++ {
		r1 += Payoff[m1][m2]
		r2 += Payoff[m2][m1]
		m1, m2 = p1.RePly(m2), p2.RePly(m1)
	}
	return r1, r2
}

func MatchPlus(p1, p2 Player, turns int) (int, int, GameTrace, GameTrace) {
	fp1, fp2 := GameTrace{}, GameTrace{}
	r1, r2 := 0, 0
	m1, m2 := p1.FirstPly(), p2.FirstPly()

	for t := 0; t < turns; t++ {
		fp1, fp2 = append(fp1, m1), append(fp2, m2)
		r1 += Payoff[m1][m2]
		r2 += Payoff[m2][m1]
		m1, m2 = p1.RePly(m2), p2.RePly(m1)
	}
	return r1, r2, fp1, fp2
}

//////////////////////////////////////////////////////////////////////////////
// Simple enumerator: returns the n-th possible path of a given size

func Enumerate(size, n int) ([]PPly, bool) {
	var last bool
	aux := make([]int, size)

	aux[0] = 1
	for i := 1; i < size; i++ {
		aux[i] = aux[i-1] * 2
	}
	if n < aux[size-1]*2-1 {
		last = false
	} else {
		last = true
	}

	result := make([]PPly, size)
	for i := range result {
		result[i] = [2]PPly{COOPERATE, DEFECT}[(n/aux[i])%2]
	}

	return result, last
}

//////////////////////////////////////////////////////////////////////////////
// Complete showdown + a few interfaces to sort scores

type TableScore struct {
	n string
	s int
}

type Table []TableScore

func (t Table) Len() int {
	return len(t)
}
func (t Table) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t Table) Less(i, j int) bool {
	return t[i].s > t[j].s
}

func Showdown(opp1, opp2 []Player, turns []int) {
	reward := make(Table, len(opp1))
	score := make(Table, len(opp1))
	table := make([][]float64, len(opp1))

	log.Println("SHOWDOWN:", turns)

	for t := 0; t < len(opp1); t++ {
		table[t] = make([]float64, len(opp1))
		score[t].n = opp1[t].GetName()
		reward[t].n = opp1[t].GetName()
	}

	tturns := 0
	for _, t := range turns {
		log.Println("Current game length :", t)
		tturns += t * len(opp1) * 2
		for o1 := 0; o1 < len(opp1); o1++ {
			for o2 := o1; o2 < len(opp1); o2++ {
				r1, r2, fp1, fp2 := MatchPlus(opp1[o1], opp2[o2], t)
				reward[o1].s += r1
				reward[o2].s += r2
				if r1 == r2 {
					score[o1].s += 1
					score[o2].s += 1
				} else if r1 > r2 {
					score[o1].s += 3
				} else {
					score[o2].s += 3
				}
				table[o1][o2] += float64(r1)
				table[o2][o1] += float64(r1)
				msg := fmt.Sprintf("%v %v  vs.  %v %v ", opp1[o1].GetName(), r1, opp2[o2].GetName(), r2)
				log.Println(msg)
				Scoresheet(opp1[o1].GetName(), opp2[o2].GetName(), t, r1, r2, fp1, fp2)
			}
		}
	}

	//log.Println("TABLE:")
	//sort.Sort(Table(score))
	//for t := 0; t < len(opponent); t++ {
	//	log.Printf("%02d) %-30s %3d\n", t+1, score[t].n, score[t].s)
	//}
	log.Println("TOTAL PAYOFF:")
	sort.Sort(Table(reward))
	for t := 0; t < len(opp1); t++ {
		log.Printf("%02d) %-30s %10d     \tpoints: %6d\n", t+1, reward[t].n, reward[t].s, score[t].s)
	}
}

func Scoresheet(name1, name2 string, length, reward1, reward2 int, moves1, moves2 GameTrace) {
	if name1 > name2 {
		Scoresheet(name2, name1, length, reward2, reward1, moves2, moves1)
	} else {
		scoresheet := "(SCORESHEET) " + name1 + " vs. " + name2 + ".txt"
		var file *os.File
		if _, err := os.Stat(scoresheet); err != nil {
			if file, err = os.OpenFile(scoresheet, os.O_RDWR|os.O_CREATE, 0644); err != nil {
				log.Println("Can't write log: ", scoresheet)
				log.Panicln(err)
			}
			file.WriteString(fmt.Sprintf("OPPONENT 1:: %v\nOPPONENT 2:: %v\n\n", name1, name2))
		} else {
			if file, err = os.OpenFile(scoresheet, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
				log.Println("Can't write log: ", scoresheet)
				log.Panicln(err)
			}
		}
		file.WriteString(fmt.Sprintf("DATE: %v\n", time.Now()))
		hn, _ := os.Hostname()
		file.WriteString(fmt.Sprintf("HOST: %v\n", hn))
		file.WriteString(fmt.Sprintf("GAME LENGTH: %v\n", length))
		if reward1 > reward2 {
			file.WriteString(fmt.Sprintf("WINNER: %v (%+f per turn)\n", name1, (float64(reward1)-float64(reward2))/float64(length)))
		} else if reward1 < reward2 {
			file.WriteString(fmt.Sprintf("WINNER: %v (%+f per turn)\n", name2, (float64(reward2)-float64(reward1))/float64(length)))
		} else {
			file.WriteString(fmt.Sprintf("WINNER: - (tie)\n"))
		}
		file.WriteString(fmt.Sprintf("REWARD OPPONENT 1: %v (%.3f avg)\n", reward1, float64(reward1)/float64(length)))
		file.WriteString(fmt.Sprintf("REWARD OPPONENT 2: %v (%.3f avg)\n", reward2, float64(reward2)/float64(length)))
		file.WriteString(fmt.Sprintf("FULL GAME TRACE:\n%v\n%v\n\n", moves1, moves2))
		file.Close()
	}
}
