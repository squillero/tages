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

func StdStrategies() []Player {
	// Zdet
	var ext2 Extort2
	var fixed ZDGTFT2_fixed
	var extortion ZDGTFT2_extortion
	var generous ZDGTFT2_generous
	// std
	var allc std_allc
	var alld std_alld
	var rand std_rand
	var grim std_grim
	var sgrim std_sgrim
	var spavlov simplified_pavlov
	var gradual std_gradual
	var sm std_sm
	var hm std_hm
	var prober std_prober
	var fbf std_fbf
	// Variants of TFT
	var tft std_tft
	var stft std_stft
	var htft std_htft
	var htft_2 std_htft_2
	var gtfg std_gtft
	var gtft_0_33 std_gtft_0_33
	var gtft_0_1 std_gtft_0_1
	var rtfg std_rtft
	var tftt std_tftt
	var tftt2 std_tftt2
	var ttft std_ttft
	var atft std_atft
	var ctft std_ctft
	var np std_np
	var np_0_1 std_np_0_1
	var rp std_rp
	var rp_0_1 std_rp_0_1
	var pavlov std_pavlov
	var otft std_otft
	var adapt std_adapt
	var bob FairStrategy
	var npm std_npm

	return []Player{&ext2, &fixed, &extortion, &generous,
		&allc, &alld, &rand, &grim, &sgrim, &spavlov, &gradual, &sm, &hm, &prober, &fbf,
		&tft, &atft, &stft, &htft, &gtfg, &rtfg, &tftt, &ttft, &np, &rp,
		&pavlov, &otft, &adapt, &bob, &npm, &tftt2, &gtft_0_1, &gtft_0_33, &np_0_1,
		&rp_0_1, &htft_2, &ctft}
}

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
	if rand.Float32() < 0.5 {

		return DEFECT
	} else {

		return COOPERATE
	}
}

func (p *FairStrategy) RePly(oppPly PPly) PPly {
	p.History[oppPly]++

	c := p.History[COOPERATE]
	d := p.History[DEFECT]
	for c+d < 10 {
		c++
		d++
	}

	if rand.Intn(c+d) < d {

		return DEFECT
	} else {

		return COOPERATE
	}
}

//////////////////////////////////////////////////////////////////////////////
// Zero Determinant Strategy: EXTORT-2
// [A.J. Stewart, J.B. Plotkin - Extortion and cooperation in the Prisoner's Dilemma]

type Extort2 struct {
	LastMove int
}

func (_ *Extort2) GetName() string {
	return "ZDET# Extort2"
}

func (p *Extort2) FirstPly() PPly {
	if rand.Intn(2) < 1 {
		p.LastMove = DEFECT
		return DEFECT
	}
	p.LastMove = COOPERATE
	return COOPERATE
}

func (p *Extort2) RePly(oppPly PPly) PPly {
	if p.LastMove == COOPERATE && oppPly == COOPERATE {
		if rand.Intn(9) < 8 {
			return COOPERATE
		}
		p.LastMove = DEFECT
		return DEFECT
	}
	if p.LastMove == COOPERATE && oppPly == DEFECT {
		if rand.Intn(2) < 1 {
			return COOPERATE
		}
		p.LastMove = DEFECT
		return DEFECT
	}
	if p.LastMove == DEFECT && oppPly == COOPERATE {
		if rand.Intn(3) < 1 {
			p.LastMove = COOPERATE
			return COOPERATE
		}
		return DEFECT
	}
	return DEFECT
}

//////////////////////////////////////////////////////////////////////////////
// Zero Determinant Strategy: ZDGTFT-2 (Fixed Score)
// [W.H. Press, F.J. Dyson - Iterated Prisoner's Dilemma contains strategies
// that dominate any evolutionary opponent]
//
// http://s3.boskent.com/prisoners-dilemma/fixed.html
//
// In this game, I've decided I want your average score to be 2.
// You'll see that, however you play, after a few hundred moves your average score will be approximately 2.

type ZDGTFT2_fixed struct {
	LastMove int
}

func (_ *ZDGTFT2_fixed) GetName() string {
	return "ZDET# ZDGTFT-2 Fixed Score"
}

func (p *ZDGTFT2_fixed) FirstPly() PPly {
	if rand.Intn(2) < 1 {
		p.LastMove = DEFECT
		return DEFECT
	}
	p.LastMove = COOPERATE
	return COOPERATE
}

func (p *ZDGTFT2_fixed) RePly(oppPly PPly) PPly {
	if oppPly == COOPERATE {
		if rand.Intn(3) < 2 {
			p.LastMove = COOPERATE
			return COOPERATE
		}
		p.LastMove = DEFECT
		return DEFECT
	}
	if p.LastMove == COOPERATE && oppPly == DEFECT {
		p.LastMove = DEFECT
		return DEFECT
	}
	if rand.Intn(3) < 1 {
		p.LastMove = COOPERATE
		return COOPERATE
	}
	p.LastMove = DEFECT
	return DEFECT
}

//////////////////////////////////////////////////////////////////////////////
// Zero Determinant Strategy: ZDGTFT-2 (Extortion)
// [W.H. Press, F.J. Dyson - Iterated Prisoner's Dilemma contains strategies
// that dominate any evolutionary opponent]
//
// http://s3.boskent.com/prisoners-dilemma/extortion.html
//
// In this game I will extort you. Your best strategy, if you want to maximise your own score,
// is to cooperate all the time, but then I will occasionally defect and so always do better than you.
//
// The only way you can avoid being taken advantage of is to resign yourself to the meagre rewards
// of mutual defection, and to defect on every turn. If you do anything else then I will take advantage
// of your cooperation and I will do three times better than you, in the sense that
// on average over the long run (my score minus 1) will be thrice (your score minus 1).

type ZDGTFT2_extortion struct {
	LastMove int
}

func (_ *ZDGTFT2_extortion) GetName() string {
	return "ZDET# ZDGTFT-2 Extortion"
}

func (p *ZDGTFT2_extortion) FirstPly() PPly {
	if rand.Intn(2) < 1 {
		p.LastMove = DEFECT
		return DEFECT
	}
	p.LastMove = COOPERATE
	return COOPERATE
}

func (p *ZDGTFT2_extortion) RePly(oppPly PPly) PPly {
	if p.LastMove == COOPERATE && oppPly == COOPERATE {
		if rand.Intn(13) < 11 {
			return COOPERATE
		}
		p.LastMove = DEFECT
		return DEFECT
	}
	if p.LastMove == DEFECT && oppPly == COOPERATE {
		if rand.Intn(26) < 7 {
			p.LastMove = COOPERATE
			return COOPERATE
		}
		return DEFECT
	}
	if p.LastMove == COOPERATE && oppPly == DEFECT {
		if rand.Intn(2) < 1 {
			return COOPERATE
		}
		p.LastMove = DEFECT
		return DEFECT
	}
	p.LastMove = DEFECT
	return DEFECT
}

//////////////////////////////////////////////////////////////////////////////
// Zero Determinant Strategy: ZDGTFT-2 (Generous)
// [W.H. Press, F.J. Dyson - Iterated Prisoner's Dilemma contains strategies
// that dominate any evolutionary opponent]
//
// http://s3.boskent.com/prisoners-dilemma/titfer.html
//
// This is possibly the strategy that slightly bettered generous
// Tit for Tat in Stewart and Plotkin's simulation.
// It is fair in the sense that if you always cooperate then I will too.
// It is a self-sacrificing strategy: if you do not always cooperate I will
// do worse than you. In fact on average I lose precisely twice as much as you:
// if your average score is less than 3, mine will be less by as much again.

type ZDGTFT2_generous struct {
	LastMove int
}

func (_ *ZDGTFT2_generous) GetName() string {
	return "ZDET# ZDGTFT-2 Generous"
}

func (p *ZDGTFT2_generous) FirstPly() PPly {
	if rand.Intn(2) < 1 {
		return DEFECT
	}
	return COOPERATE
}

func (p *ZDGTFT2_generous) RePly(oppPly PPly) PPly {
	if p.LastMove == COOPERATE && oppPly == COOPERATE {
		return COOPERATE
	}
	if p.LastMove == DEFECT && oppPly == COOPERATE {
		if rand.Intn(10) < 8 {
			p.LastMove = COOPERATE
			return COOPERATE
		}
		return DEFECT
	}
	if p.LastMove == COOPERATE && oppPly == DEFECT {
		if rand.Intn(10) < 3 {
			return COOPERATE
		}
		p.LastMove = DEFECT
		return DEFECT
	}
	if rand.Intn(10) < 2 {
		p.LastMove = COOPERATE
		return COOPERATE
	}
	p.LastMove = DEFECT
	return DEFECT
}

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////
// Strategies from http://www.prisoners-dilemma.com/
//////////////////////////////////////////////////////////////////////////////

//////////////////////////////////////////////////////////////////////////////
// Cooperates on every move.

type std_allc struct {
}

func (_ *std_allc) GetName() string {
	return "STD# Always Cooperate"
}

func (p *std_allc) FirstPly() PPly {
	return COOPERATE
}

func (p *std_allc) RePly(oppPly PPly) PPly {
	return COOPERATE
}

//////////////////////////////////////////////////////////////////////////////
// Defect on every move.

type std_alld struct {
}

func (_ *std_alld) GetName() string {
	return "STD# Always Defect"
}

func (p *std_alld) FirstPly() PPly {
	return DEFECT
}

func (p *std_alld) RePly(oppPly PPly) PPly {
	return DEFECT
}

//////////////////////////////////////////////////////////////////////////////
// Makes a random move.

type std_rand struct {
}

func (_ *std_rand) GetName() string {
	return "STD# Random Player"
}

func (p *std_rand) FirstPly() PPly {
	return [2]PPly{COOPERATE, DEFECT}[rand.Int()%2]
}

func (p *std_rand) RePly(oppPly PPly) PPly {
	return [2]PPly{COOPERATE, DEFECT}[rand.Int()%2]
}

//////////////////////////////////////////////////////////////////////////////
// Cooperates, until the opponent defects, and thereafter always defects

type std_grim struct {
	grimming bool
}

func (_ *std_grim) GetName() string {
	return "STD# Grim Trigger"
}

func (p *std_grim) FirstPly() PPly {
	p.grimming = false
	return COOPERATE
}

func (p *std_grim) RePly(oppPly PPly) PPly {
	if oppPly == DEFECT {
		p.grimming = true
	}
	if p.grimming {
		return DEFECT
	} else {
		return COOPERATE
	}
}

//////////////////////////////////////////////////////////////////////////////
// Cooperates on the first move. If a reward or temptation payoff is received
// in the last round then repeats last choice, otherwise chooses the opposite choice

type simplified_pavlov struct {
	OppLastPly, MyLastPly PPly
}

func (_ *simplified_pavlov) GetName() string {
	return "STD# Simplified Pavlov"
}

func (p *simplified_pavlov) FirstPly() PPly {
	p.OppLastPly = UNSPECIFIED
	p.MyLastPly = COOPERATE
	return p.MyLastPly
}

func (p *simplified_pavlov) RePly(oppPly PPly) PPly {
	p.OppLastPly = oppPly
	if p.OppLastPly == COOPERATE && p.MyLastPly == COOPERATE {
		// reward: repeat last moves
	} else if p.OppLastPly == DEFECT && p.MyLastPly == COOPERATE {
		// temptation: repeat last moves
	} else {
		p.MyLastPly = OppositeMove(p.MyLastPly)
	}
	return p.MyLastPly
}

//////////////////////////////////////////////////////////////////////////////
// Cooperates on the first move, and cooperates as long as the opponent cooperates.
// After the first defection of the other player, it defects one time and cooperates two times;
// ... After the nth defection it reacts with n consecutive defections and then
// calms down its opponent with two cooperations

type std_gradual struct {
	OppDefections, MyStepNUmber int
}

func (_ *std_gradual) GetName() string {
	return "STD# Gradual"
}

func (p *std_gradual) FirstPly() PPly {
	p.OppDefections, p.MyStepNUmber = 0, 0

	return COOPERATE
}

func (p *std_gradual) RePly(oppPly PPly) PPly {
	var move PPly
	if p.MyStepNUmber == 0 {
		if oppPly == DEFECT {
			p.OppDefections++
			p.MyStepNUmber = 1
			move = DEFECT
		} else {
			move = COOPERATE
		}
	} else if p.MyStepNUmber <= p.OppDefections {
		p.MyStepNUmber++
		move = DEFECT
	} else if p.MyStepNUmber == p.OppDefections {
		p.MyStepNUmber++
		move = COOPERATE
	} else {
		p.MyStepNUmber = 0
		move = COOPERATE
	}

	return move
}

//////////////////////////////////////////////////////////////////////////////
// Cooperates on the first move, and cooperates as long as the number of times
// the opponent has cooperated is greater than or equal to the number of times
// it has defected, else it defects

type std_sm struct {
	Defections, Cooperations int
}

func (_ *std_sm) GetName() string {
	return "STD# Soft Majority"
}

func (p *std_sm) FirstPly() PPly {
	p.Defections, p.Cooperations = 0, 0

	return COOPERATE
}

func (p *std_sm) RePly(oppPly PPly) PPly {
	if oppPly == COOPERATE {
		p.Cooperations++
	} else {
		p.Defections++
	}

	if p.Cooperations >= p.Defections {
		return COOPERATE
	} else {
		return DEFECT
	}
}

//////////////////////////////////////////////////////////////////////////////
// Defects on the first move, and defects if the number of defections of the opponent
// is greater than or equal to the number of times it has cooperated, else cooperates

type std_hm struct {
	Defections, Cooperations int
}

func (_ *std_hm) GetName() string {
	return "STD# Hard Majority"
}

func (p *std_hm) FirstPly() PPly {
	p.Defections, p.Cooperations = 0, 0
	return DEFECT
}

func (p *std_hm) RePly(oppPly PPly) PPly {
	if oppPly == COOPERATE {
		p.Cooperations++
	} else {
		p.Defections++
	}

	if p.Defections >= p.Cooperations {
		return DEFECT
	} else {
		return COOPERATE
	}
}

//////////////////////////////////////////////////////////////////////////////
// Like GRIM except that the opponent is punished with D,D,D,D,C,C

type std_sgrim struct {
	Sequence    []PPly
	seqPosition int
}

func (_ *std_sgrim) GetName() string {
	return "STD# Soft Grudger"
}

func (p *std_sgrim) FirstPly() PPly {
	p.Sequence = []PPly{DEFECT, DEFECT, DEFECT, DEFECT, COOPERATE, COOPERATE}
	p.seqPosition = 0
	return COOPERATE
}

func (p *std_sgrim) RePly(oppPly PPly) PPly {
	if p.seqPosition > 0 {
		// next step in the sequence
		m := p.Sequence[p.seqPosition]
		p.seqPosition = (p.seqPosition + 1) % len(p.Sequence)
		return m
	} else {
		if oppPly == DEFECT {
			p.seqPosition = 1
			return p.Sequence[0]
		} else {
			return COOPERATE
		}
	}
}

//////////////////////////////////////////////////////////////////////////////
// Cooperates on the first move, then copies the opponent's last move

type std_tft struct {
}

func (_ *std_tft) GetName() string {
	return "TFT# Tit for Tat"
}

func (p *std_tft) FirstPly() PPly {
	return COOPERATE
}

func (p *std_tft) RePly(oppPly PPly) PPly {
	return oppPly
}

//////////////////////////////////////////////////////////////////////////////
// An adaption rate r is used to compute a continuous variable ‘world’
// according to the history moves of the opponent

type std_atft struct {
	world        float64
	adaptionRate float64
}

func (_ *std_atft) GetName() string {
	return "TFT# Adaptive Tit for Tat"
}

func (p *std_atft) FirstPly() PPly {
	p.world = 0.5
	p.adaptionRate = 0.1
	return COOPERATE
}

func (p *std_atft) RePly(oppPly PPly) PPly {
	if oppPly == COOPERATE {
		p.world = p.world + p.adaptionRate*(1-p.world)
	} else {
		p.world = p.world + p.adaptionRate*(0-p.world)
	}
	if p.world >= 0.5 {
		return COOPERATE
	} else {
		return DEFECT
	}
}

//////////////////////////////////////////////////////////////////////////////
// Contrite TFT: Same as TFT when no noise. In a noisy environment, once it
// receives T because of error, it will choose cooperate twice in order recover
// mutual cooperation

type std_ctft struct {
	good                  bool
	oppLastPly, myLastPly PPly
}

func (_ *std_ctft) GetName() string {
	return "TFT# Contrite Tit for Tat"
}

func (p *std_ctft) FirstPly() PPly {
	p.myLastPly = COOPERATE
	p.oppLastPly = UNSPECIFIED
	p.good = true
	return p.myLastPly
}

func (p *std_ctft) RePly(oppPly PPly) PPly {
	p.oppLastPly = oppPly

	if p.oppLastPly == COOPERATE {
		p.good = true
	} else if p.oppLastPly == DEFECT && p.myLastPly == DEFECT {
		p.good = true
	} else {
		p.good = false
	}

	if p.good {
		p.myLastPly = COOPERATE
	} else {
		p.myLastPly = DEFECT
	}
	return p.myLastPly
}

//////////////////////////////////////////////////////////////////////////////
// It does the reverse of TFT. It defects on the first move,
// then plays the reverse of the opponent’s last move

type std_rtft struct {
}

func (_ *std_rtft) GetName() string {
	return "TFT# Reverse Tit for Tat"
}

func (p *std_rtft) FirstPly() PPly {
	return DEFECT
}

func (p *std_rtft) RePly(oppPly PPly) PPly {
	return OppositeMove(oppPly)
}

//////////////////////////////////////////////////////////////////////////////
// Same as TFT, except that it defects on the first move.

type std_stft struct {
}

func (_ *std_stft) GetName() string {
	return "TFT# Suspicious Tit for Tat"
}

func (p *std_stft) FirstPly() PPly {
	return DEFECT
}

func (p *std_stft) RePly(oppPly PPly) PPly {
	return oppPly
}

//////////////////////////////////////////////////////////////////////////////
// Cooperates on the first move, and defects if the opponent has defects
// on any of the previous three moves, else cooperates

type std_htft struct {
	oppLastPly []PPly
	num        int
}

func (_ *std_htft) GetName() string {
	return "TFT# Hard Tit for Tat [w=3]"
}

func (p *std_htft) FirstPly() PPly {
	p.oppLastPly = []PPly{UNSPECIFIED, UNSPECIFIED, UNSPECIFIED}
	p.num = 0
	return COOPERATE
}

func (p *std_htft) RePly(oppPly PPly) PPly {
	p.oppLastPly[p.num] = oppPly
	p.num = (p.num + 1) % len(p.oppLastPly)
	for _, o := range p.oppLastPly {
		if o == DEFECT {
			return DEFECT
		}
	}
	return COOPERATE
}

//////////////////////////////////////////////////////////////////////////////
// Cooperates on the first move, and defects if the opponent has defects
// on any of the previous two moves, else cooperates

type std_htft_2 struct {
	oppLastPly []PPly
	num        int
}

func (_ *std_htft_2) GetName() string {
	return "TFT# Hard Tit for Tat [w=2]"
}

func (p *std_htft_2) FirstPly() PPly {
	p.oppLastPly = []PPly{UNSPECIFIED, UNSPECIFIED}
	p.num = 0
	return COOPERATE
}

func (p *std_htft_2) RePly(oppPly PPly) PPly {
	p.oppLastPly[p.num] = oppPly
	p.num = (p.num + 1) % len(p.oppLastPly)
	for _, o := range p.oppLastPly {
		if o == DEFECT {
			return DEFECT
		}
	}
	return COOPERATE
}

//////////////////////////////////////////////////////////////////////////////
// Like Tit for Tat, but occasionally defects with a small probability.

type std_np_0_1 struct {
	ε float64
}

func (_ *std_np_0_1) GetName() string {
	return "TFT# Naïve Prober [ε=0.01]"
}

func (p *std_np_0_1) FirstPly() PPly {
	p.ε = 0.01
	return COOPERATE
}

func (p *std_np_0_1) RePly(oppPly PPly) PPly {
	if rand.Float64() > p.ε {
		return oppPly
	} else {
		return DEFECT
	}
}

//////////////////////////////////////////////////////////////////////////////
// Like Tit for Tat, but occasionally defects with a small probability.

type std_np struct {
	ε float64
}

func (_ *std_np) GetName() string {
	return "TFT# Naïve Prober [ε=0.1]"
}

func (p *std_np) FirstPly() PPly {
	p.ε = 0.1
	return COOPERATE
}

func (p *std_np) RePly(oppPly PPly) PPly {
	if rand.Float64() > p.ε {
		return oppPly
	} else {
		return DEFECT
	}
}

//////////////////////////////////////////////////////////////////////////////
//  Like Naive Prober, but it tries to break the series of mutual defections after defecting.

type std_rp struct {
	unfair bool
	ε      float64
}

func (_ *std_rp) GetName() string {
	return "TFT# Remorseful Prober [ε=0.01]"
}

func (p *std_rp) FirstPly() PPly {
	p.ε = 0.01
	p.unfair = false
	return COOPERATE
}

func (p *std_rp) RePly(oppPly PPly) PPly {
	if rand.Float64() > p.ε {
		if oppPly == DEFECT && p.unfair == true {
			p.unfair = false
			return COOPERATE
		} else {
			p.unfair = false
			return oppPly
		}
	} else {
		p.unfair = true
		return DEFECT
	}
}

//////////////////////////////////////////////////////////////////////////////
//  Like Naive Prober, but it tries to break the series of mutual defections after defecting.

type std_rp_0_1 struct {
	unfair bool
	ε      float64
}

func (_ *std_rp_0_1) GetName() string {
	return "TFT# Remorseful Prober [ε=0.1]"
}

func (p *std_rp_0_1) FirstPly() PPly {
	p.ε = 0.1
	p.unfair = false
	return COOPERATE
}

func (p *std_rp_0_1) RePly(oppPly PPly) PPly {
	if rand.Float64() > p.ε {
		if oppPly == DEFECT && p.unfair == true {
			p.unfair = false
			return COOPERATE
		} else {
			p.unfair = false
			return oppPly
		}
	} else {
		p.unfair = true
		return DEFECT
	}
}

//////////////////////////////////////////////////////////////////////////////
// Cooperates on the first move, and defects only when the opponent defects two times

type std_tftt struct {
	oppLastPly [2]PPly
}

func (_ *std_tftt) GetName() string {
	return "TFT# Tit for Two Tats"
}

func (p *std_tftt) FirstPly() PPly {
	p.oppLastPly[0], p.oppLastPly[1] = UNSPECIFIED, UNSPECIFIED
	return COOPERATE
}

func (p *std_tftt) RePly(oppPly PPly) PPly {
	p.oppLastPly[1], p.oppLastPly[0] = p.oppLastPly[0], oppPly

	if p.oppLastPly[0] == DEFECT && p.oppLastPly[1] == DEFECT {
		return DEFECT
	} else {
		return COOPERATE
	}
}

//////////////////////////////////////////////////////////////////////////////
// Same as Tit for Tat except that it defects twice when the opponent defects

type std_ttft struct {
	oppLastPly [2]PPly
}

func (_ *std_ttft) GetName() string {
	return "TFT# Two Tits for Tat"
}

func (p *std_ttft) FirstPly() PPly {
	p.oppLastPly[0], p.oppLastPly[1] = UNSPECIFIED, UNSPECIFIED
	return COOPERATE
}

func (p *std_ttft) RePly(oppPly PPly) PPly {
	p.oppLastPly[1], p.oppLastPly[0] = p.oppLastPly[0], oppPly

	if p.oppLastPly[0] == DEFECT || p.oppLastPly[1] == DEFECT {
		return DEFECT
	} else {
		return COOPERATE
	}
}

//////////////////////////////////////////////////////////////////////////////
// Same as TFT, except that it cooperates with a probability q when the opponent defects

type std_gtft struct {
	ε float64
}

func (_ *std_gtft) GetName() string {
	return "TFT# Generous Tit for Tat [ε=0.05]"
}

func (p *std_gtft) FirstPly() PPly {
	p.ε = 0.05
	return COOPERATE
}

func (p *std_gtft) RePly(oppPly PPly) PPly {
	if rand.Float64() < p.ε {
		return COOPERATE
	} else {
		return oppPly
	}
}

//////////////////////////////////////////////////////////////////////////////
// Same as TFT, except that it cooperates with a probability q when the opponent defects

type std_gtft_0_1 struct {
	ε float64
}

func (_ *std_gtft_0_1) GetName() string {
	return "TFT# Generous Tit for Tat [ε=0.1]"
}

func (p *std_gtft_0_1) FirstPly() PPly {
	p.ε = 0.1
	return COOPERATE
}

func (p *std_gtft_0_1) RePly(oppPly PPly) PPly {
	if rand.Float64() < p.ε {
		return COOPERATE
	} else {
		return oppPly
	}
}

//////////////////////////////////////////////////////////////////////////////
// Same as TFT, except that it cooperates with a probability q when the opponent defects

type std_gtft_0_33 struct {
	ε float64
}

func (_ *std_gtft_0_33) GetName() string {
	return "TFT# Generous Tit for Tat [ε=0.33]"
}

func (p *std_gtft_0_33) FirstPly() PPly {
	p.ε = 0.33
	return COOPERATE
}

func (p *std_gtft_0_33) RePly(oppPly PPly) PPly {
	if rand.Float64() < p.ε {
		return COOPERATE
	} else {
		return oppPly
	}
}

//////////////////////////////////////////////////////////////////////////////
// Starts with D,C,C and then defects if the opponent has cooperated in
// the second and third move; otherwise, it plays TFT.

type std_prober struct {
	Sequence, OppPlys []PPly
	seqPosition       int
	oppLastPly        PPly
}

func (_ *std_prober) GetName() string {
	return "STD# Prober"
}

func (p *std_prober) FirstPly() PPly {
	p.Sequence = []PPly{DEFECT, COOPERATE, COOPERATE}
	p.OppPlys = make([]PPly, len(p.Sequence))
	p.seqPosition = 1
	p.oppLastPly = UNSPECIFIED
	return p.Sequence[0]
}

func (p *std_prober) RePly(oppPly PPly) PPly {
	t := p.oppLastPly
	p.oppLastPly = oppPly
	if p.seqPosition > 0 {
		p.OppPlys[p.seqPosition-1] = oppPly
		// next step in the sequence
		m := p.Sequence[p.seqPosition]
		p.seqPosition = (p.seqPosition + 1) % len(p.Sequence)
		return m
	} else if p.OppPlys[1] == COOPERATE && p.OppPlys[2] == COOPERATE {
		return DEFECT
	} else {
		return t
	}
}

//////////////////////////////////////////////////////////////////////////////
// Cooperates on the first move, and cooperates except after receiving a sucker payoff

type std_fbf struct {
	oppLastPly, myLastPly PPly
}

func (_ *std_fbf) GetName() string {
	return "STD# Firm But Fair"
}

func (p *std_fbf) FirstPly() PPly {
	p.oppLastPly = COOPERATE
	p.myLastPly = COOPERATE
	return p.myLastPly
}

func (p *std_fbf) RePly(oppPly PPly) PPly {
	p.oppLastPly = oppPly
	if p.oppLastPly == DEFECT && p.myLastPly == COOPERATE {
		p.myLastPly = DEFECT
	} else {
		p.myLastPly = COOPERATE
	}
	return p.myLastPly
}

//////////////////////////////////////////////////////////////////////////////
//

type std_pavlov struct {
	oppLastPly_4_TFTT                      [2]PPly
	OppLastPly, MyLastPly                  PPly
	totPayoff, seqPosition, nDef, behavior int
}

func (_ *std_pavlov) GetName() string {
	return "ADAPT# Pavlov"
}

func (p *std_pavlov) FirstPly() PPly {
	p.seqPosition = 0
	p.behavior = 0
	p.nDef = 0

	p.oppLastPly_4_TFTT[0], p.oppLastPly_4_TFTT[1] = UNSPECIFIED, UNSPECIFIED
	p.OppLastPly = UNSPECIFIED
	p.MyLastPly = COOPERATE
	return p.MyLastPly
}

func (p *std_pavlov) RePly(oppPly PPly) PPly {
	p.OppLastPly = oppPly
	p.oppLastPly_4_TFTT[1], p.oppLastPly_4_TFTT[0] = p.oppLastPly_4_TFTT[0], oppPly
	p.seqPosition++

	p.totPayoff += Payoff[p.MyLastPly][p.OppLastPly]

	if p.seqPosition == 6 {
		// computation of average payoff
		if p.nDef == 0 {
			p.behavior = 0
		} else if p.nDef == 3 {
			p.behavior = 1
		} else {
			p.behavior = 2
		}
		p.nDef = 0
		p.seqPosition = 0
		p.totPayoff = 0
	} else if p.seqPosition%6 == 0 {
		if p.totPayoff < REWARD*6 {
			p.seqPosition = 0
			p.nDef = 0
			p.behavior = 0
		}
		p.totPayoff = 0
	}
	if p.OppLastPly == DEFECT {
		p.nDef++
	}
	if p.behavior == 0 {
		// Plays as a TFT
		p.MyLastPly = p.OppLastPly
	} else if p.behavior == 1 {
		// Plays as a TFTT
		if p.oppLastPly_4_TFTT[0] == DEFECT && p.oppLastPly_4_TFTT[1] == DEFECT {
			p.MyLastPly = DEFECT
		} else {
			p.MyLastPly = COOPERATE
		}
	} else if p.behavior == 2 {
		// Plays as a Always Defect
		p.MyLastPly = DEFECT
	}

	return p.MyLastPly
}

//////////////////////////////////////////////////////////////////////////////
// ADAPT# Adaptive

type std_adapt struct {
	OppLastPly, MyLastPly      PPly
	totPayoff, nMyCoop, nMyDef int
	avgCoop, avgDef            float64
}

func (_ *std_adapt) GetName() string {
	return "ADAPT# Adaptive"
}

func (p *std_adapt) FirstPly() PPly {
	p.nMyCoop = 0
	p.avgCoop = 0
	p.avgDef = 0
	p.nMyDef = 0
	p.OppLastPly = UNSPECIFIED
	p.MyLastPly = COOPERATE
	return p.MyLastPly
}

func (p *std_adapt) RePly(oppPly PPly) PPly {
	p.OppLastPly = oppPly

	// update statistics
	if p.MyLastPly == COOPERATE {
		p.avgCoop = ((p.avgCoop * float64(p.nMyCoop)) + float64(Payoff[p.MyLastPly][p.OppLastPly])) / float64(p.nMyCoop+1)
		p.nMyCoop++
	} else {
		p.avgDef = ((p.avgDef * float64(p.nMyDef)) + float64(Payoff[p.MyLastPly][p.OppLastPly])) / float64(p.nMyDef+1)
		p.nMyDef++
	}

	if (p.nMyCoop + p.nMyDef) < 10 {
		// starting phase
		if (p.nMyCoop + p.nMyDef) < 5 {
			p.MyLastPly = COOPERATE
		} else {
			p.MyLastPly = DEFECT
		}
	} else {
		if p.avgDef > p.avgCoop {
			p.MyLastPly = DEFECT
		} else {
			p.MyLastPly = COOPERATE
		}
	}

	return p.MyLastPly
}

//////////////////////////////////////////////////////////////////////////////
// Omega TFT

type std_otft struct {
	OppLastPly, MyLastPly                                    PPly
	Deadlock_th, Randomness_th, Deadlock_cnt, Randomness_cnt int
}

func (_ *std_otft) GetName() string {
	return "TFT# Omega Tit for Tat"
}

func (p *std_otft) FirstPly() PPly {
	p.OppLastPly = UNSPECIFIED
	p.MyLastPly = COOPERATE
	p.Deadlock_th = 3
	p.Randomness_th = 8
	p.Deadlock_cnt = 0
	p.Randomness_cnt = -1
	return p.MyLastPly
}

func (p *std_otft) RePly(oppPly PPly) PPly {
	if p.Deadlock_cnt >= p.Deadlock_th {
		p.MyLastPly = COOPERATE
		if p.Deadlock_cnt == p.Deadlock_th {
			p.Deadlock_cnt = p.Deadlock_th + 1
		} else {
			p.Deadlock_cnt = 0
		}
	} else {
		if oppPly == COOPERATE && p.OppLastPly == COOPERATE {
			p.Randomness_cnt--
		}
		if p.OppLastPly != oppPly {
			p.Randomness_cnt++
		}
		if p.MyLastPly != oppPly {
			p.Randomness_cnt++
		}
		if p.Randomness_cnt >= p.Randomness_th {
			p.MyLastPly = DEFECT
		} else {
			p.MyLastPly = oppPly
			if oppPly != p.OppLastPly {
				p.Deadlock_cnt++
			} else {
				p.Deadlock_cnt = 0
			}
		}
	}
	p.OppLastPly = oppPly

	return p.MyLastPly
}

//////////////////////////////////////////////////////////////////////////////
//  Naive Peace Maker (Tit For Tat with Random Co-operation)
//  Repeat opponent's last choice (ie Tit For Tat), but sometimes make peace by co-operating in lieu of defecting

type std_npm struct {
	ε float64
}

func (_ *std_npm) GetName() string {
	return "TFT# Naïve Peace Maker"
}

func (p *std_npm) FirstPly() PPly {
	p.ε = 0.01
	return COOPERATE
}

func (p *std_npm) RePly(oppPly PPly) PPly {
	if rand.Float64() > p.ε {
		return oppPly
	} else {
		return COOPERATE
	}
}

type std_tftt2 struct {
	oppLastPly [2]PPly
	lastMove   PPly
}

func (_ *std_tftt2) GetName() string {
	return "TFT# Hard Tit for Two Tats"
}

func (p *std_tftt2) FirstPly() PPly {
	p.oppLastPly[0], p.oppLastPly[1] = UNSPECIFIED, UNSPECIFIED
	p.lastMove = COOPERATE
	return p.lastMove
}

func (p *std_tftt2) RePly(oppPly PPly) PPly {
	p.oppLastPly[1], p.oppLastPly[0] = p.oppLastPly[0], oppPly

	if p.lastMove == COOPERATE {
		if p.oppLastPly[0] == DEFECT && p.oppLastPly[1] == DEFECT {
			p.lastMove = DEFECT
		}
	} else {
		if p.oppLastPly[0] == COOPERATE && p.oppLastPly[1] == COOPERATE {
			p.lastMove = COOPERATE
		}
	}
	return p.lastMove
}
