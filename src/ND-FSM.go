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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

//////////////////////////////////////////////////////////////////////////////
// Basic structures
// Pointers to DFSM do satisfy the "Player" interface (*DFSM, not DFSM)

type NdFsmNode struct {
	Active        bool
	Action        PPly
	OnDefection   NodeIndexList
	OnCooperation NodeIndexList
	canonic       int
	reachable     bool
	lastSeen      int
	color         string
}

type NdFsm struct {
	Name         string
	InitialState NodeIndexList
	Node         []NdFsmNode
	currentState NodeIndex
}

//////////////////////////////////////////////////////////////////////////////
// Stringer interfaces

func (n *NdFsmNode) String() string {
	var str string
	if n.Active {
		str = fmt.Sprintf("{ %s D:%v C:%v @%v }", n.Action.Mark(), n.OnDefection, n.OnCooperation, n.color)
	} else {
		str = fmt.Sprintf("{ %s D:%v C:%v @%v }", "*", n.OnDefection, n.OnCooperation, n.color)
	}
	return str
}
func (f *NdFsm) String() string {
	str := fmt.Sprintf("{ %s %p:%v [ ", f.Name, &f.InitialState, f.InitialState)
	for i, n := range f.Node {
		str += fmt.Sprintf("%d:%v ", i, &n)
	}
	str += "]}"
	return str
}

//////////////////////////////////////////////////////////////////////////////
// simplistic compare

func FsmEqual(fsm1, fsm2 *NdFsm) bool {
	if fsm1.Signature() == fsm2.Signature() {
		return true
	} else {
		return false
	}
}

//////////////////////////////////////////////////////////////////////////////
// Duplicate (aka. DeepCopy)

func (f *NdFsm) Duplicate() NdFsm {
	nf := new(NdFsm)
	nf.Name = f.Name
	// copy initial states
	nf.InitialState = make([]NodeIndex, len(f.InitialState))
	copy(nf.InitialState, f.InitialState)
	// copy nodes
	nf.Node = make([]NdFsmNode, len(f.Node))
	for t, n := range f.Node {
		nf.Node[t] = n.Duplicate()
	}
	return *nf
}

func (n *NdFsmNode) Duplicate() NdFsmNode {
	nn := new(NdFsmNode)
	nn.Active = n.Active
	nn.Action = n.Action
	nn.lastSeen = n.lastSeen
	nn.color = n.color
	nn.canonic = n.canonic
	nn.OnDefection = make([]NodeIndex, len(n.OnDefection))
	copy(nn.OnDefection, n.OnDefection)
	nn.OnCooperation = make([]NodeIndex, len(n.OnCooperation))
	copy(nn.OnCooperation, n.OnCooperation)
	nn.reachable = n.reachable
	return *nn
}

//////////////////////////////////////////////////////////////////////////////
// Cool debug

func (fsm *NdFsm) Debug() {
	log.Println("@", fsm)
	fsm.DumpPNG()
}

func (fsm *NdFsm) DumpPNG() {
	tfsm := fsm.Duplicate()
	//log.Println(&tfsm)
	ioutil.WriteFile(fsm.Name+".gv", tfsm.GvEncode(), 0644)
	ioutil.WriteFile(fsm.Name+".json", tfsm.JsonEncode(), 0644)
	go Gv2Png(tfsm.Name)
}

func Gv2Png(name string) {
	var png bytes.Buffer
	cmd := exec.Command("dot", "-Tpng", name+".gv")
	cmd.Stdout = &png
	cmd.Run()
	ioutil.WriteFile(name+".png", png.Bytes(), 0644)
}

//////////////////////////////////////////////////////////////////////////////
// Self explanatory

func NewBareNdFSM(name string) NdFsm {
	fsm := *new(NdFsm)
	fsm.Name = name
	fsm.Node = make([]NdFsmNode, 0)
	fsm.InitialState = make([]NodeIndex, 0)
	return fsm
}
func NewNdFSM(name string) NdFsm {
	fsm := NewBareNdFSM(name)
	fsm.currentState = -1
	fsm.AddNode(INVALID)
	fsm.AddNdTransition(0, COOPERATE, 0)
	fsm.AddNdTransition(0, DEFECT, 0)
	fsm.SetInitialState(0)
	return fsm
}

func (fsm *NdFsm) addTransition(nodeCurrent NodeIndex, action PPly, nodeTarget NodeIndex) bool {
	action.Check()
	if action == COOPERATE {
		fsm.Node[nodeCurrent].OnCooperation = append(fsm.Node[nodeCurrent].OnCooperation, nodeTarget)
	} else {
		fsm.Node[nodeCurrent].OnDefection = append(fsm.Node[nodeCurrent].OnDefection, nodeTarget)
	}
	return true
}

func (fsm *NdFsm) AddDTransition(nodeCurrent NodeIndex, action PPly, nodeTarget NodeIndex) bool {
	if fsm.CheckTransition(nodeCurrent, action, nodeTarget) {
		return false
	}
	fsm.addTransition(nodeCurrent, action, nodeTarget)
	return true
}

func (fsm *NdFsm) AddNdTransition(nodeCurrent NodeIndex, action PPly, nodeTarget NodeIndex) bool {
	if !fsm.CheckTransition(nodeCurrent, action, nodeTarget) {
		fsm.addTransition(nodeCurrent, action, nodeTarget)
		return true
	} else {
		log.Println("**************** ADD ND TRANSITION FAILED!")
		return false
	}
}

func (fsm *NdFsm) CheckTransition(nodeCurrent NodeIndex, action PPly, nodeTarget NodeIndex) bool {
	var slice []NodeIndex
	if action == COOPERATE {
		slice = fsm.Node[nodeCurrent].OnCooperation
	} else {
		slice = fsm.Node[nodeCurrent].OnDefection
	}
	for _, t := range slice {
		if t == nodeTarget {
			return true
		}
	}
	return false
}

func (fsm *NdFsm) AddNode(action PPly) NodeIndex {
	found := -1
	for t, n := range fsm.Node {
		if !n.Active {
			found = t
			break
		}
	}
	if found == -1 {
		n := new(NdFsmNode)
		n.OnCooperation = make([]NodeIndex, 0, 1)
		n.OnDefection = make([]NodeIndex, 0, 1)
		fsm.Node = append(fsm.Node, *n)
		found = len(fsm.Node) - 1
	}

	fsm.Node[found].Active = true
	fsm.Node[found].Action = action
	fsm.Node[found].color = "cyan"

	return NodeIndex(found)
}

func (fsm *NdFsm) IsInitialState(node NodeIndex) bool {
	for _, n := range fsm.InitialState {
		if node == n {
			return true
		}
	}
	return false
}

func (fsm *NdFsm) SetInitialState(i NodeIndex) {
	fsm.InitialState = make([]NodeIndex, 1)
	fsm.InitialState[0] = i
}

func (fsm *NdFsm) AddInitialState(i NodeIndex) {
	fsm.InitialState = append(fsm.InitialState, i)
}

func (fsm *NdFsm) RemoveInitialState(i NodeIndex) {
	s := -1
	for t, n := range fsm.InitialState {
		if i == n {
			s = t
		}
	}
	if s != -1 {
		fsm.InitialState[s] = fsm.InitialState[len(fsm.InitialState)-1]
		fsm.InitialState = fsm.InitialState[:len(fsm.InitialState)-1]
	}
}

func (fsm *NdFsm) ActiveNodes() []NodeIndex {
	an := make([]NodeIndex, 0, 1)
	for t := 0; t < len(fsm.Node); t++ {
		if fsm.Node[t].Active && fsm.Node[t].Action != INVALID {
			an = append(an, NodeIndex(t))
		}
	}
	return an
}

var _CanonicCount int

func (fsm *NdFsm) UpdateInternals() {
	//sort.Sort(fsm.InitialState)
	for t := range fsm.Node {
		fsm.Node[t].reachable = false
		fsm.Node[t].canonic = -1
	}
	_CanonicCount = 0
	for _, i := range fsm.InitialState {
		fsm._UpdateInternals(i)
	}
	for t := range fsm.Node {
		if !fsm.Node[t].reachable {
			fsm.Node[t].color = "gray60"
		}
	}
}

func (fsm *NdFsm) _UpdateInternals(n NodeIndex) {
	if fsm.Node[n].Active != true {
		fsm.DumpPNG()
		log.Panicf("Reaching inactive state: %v\n# %v\n", fsm.Node[n], fsm)
	}
	if fsm.Node[n].reachable == true {
		return
	}
	fsm.Node[n].reachable = true
	fsm.Node[n].canonic = _CanonicCount
	_CanonicCount++
	for t := range fsm.Node[n].OnCooperation {
		fsm._UpdateInternals(fsm.Node[n].OnCooperation[t])
	}
	for t := range fsm.Node[n].OnDefection {
		fsm._UpdateInternals(fsm.Node[n].OnDefection[t])
	}
}

func (fsm *NdFsm) GetSink() NodeIndex {
	for i, n := range fsm.Node {
		if n.Action == INVALID {
			return NodeIndex(i)
		}
	}
	n := fsm.AddNode(INVALID)
	fsm.Node[n].color = "black"
	fsm.AddDTransition(n, COOPERATE, n)
	fsm.AddDTransition(n, DEFECT, n)
	return n
}

func (n *NdFsmNode) NextState(m PPly) []NodeIndex {
	if m == COOPERATE {
		return n.OnCooperation
	} else if m == DEFECT {
		return n.OnDefection
	} else {
		log.Panicln("Wrong step: ", m)
	}
	return nil
}

func (fsmNode *NdFsmNode) CompareTo(n NdFsmNode) bool {
	if fsmNode.Active != n.Active {
		return false
	}
	if fsmNode.Action != n.Action {
		return false
	}
	if fsmNode.reachable != n.reachable {
		return false
	}
	if len(fsmNode.OnCooperation) != len(n.OnCooperation) {
		return false
	}
	if len(fsmNode.OnDefection) != len(n.OnDefection) {
		return false
	}
	for b := 0; b < len(fsmNode.OnCooperation); b++ {
		if fsmNode.OnCooperation[b] != n.OnCooperation[b] {
			return false
		}
	}
	for b := 0; b < len(fsmNode.OnDefection); b++ {
		if fsmNode.OnDefection[b] != n.OnDefection[b] {
			return false
		}
	}
	return true
}

//////////////////////////////////////////////////////////////////////////////
// Check if valid
// Rationale: quit asap, don't propagate errors

func (fsm *NdFsm) Check() bool {
	verdict := make([]string, 0)

	if len(fsm.InitialState) == 0 {
		verdict = append(verdict, "Empty initial state")
	}
	for _, i := range fsm.InitialState {
		if i < 0 || int(i) >= len(fsm.Node) || fsm.Node[i].Active == false {
			verdict = append(verdict, "Invalid initial state")
		}
	}
	// Check colors
	if *OptionColors {
		for _, n := range fsm.Node {
			if n.color == "" {
				verdict = append(verdict, fmt.Sprintf("Missing color: %v", &n))
			}
		}
	}
	// Check transitions
	for _, n := range fsm.Node {
		if n.Active {
			if len(n.OnCooperation) == 0 && n.Action != INVALID {
				verdict = append(verdict, fmt.Sprintf("Missing action on cooperation: %v", &n))
			}
			for _, u := range n.OnCooperation {
				if u < 0 || int(u) >= len(fsm.Node) || !fsm.Node[u].Active {
					verdict = append(verdict, fmt.Sprintf("Illegal on cooperation action from node: %v", &n))
				}
			}
			if len(n.OnDefection) == 0 && n.Action != INVALID {
				verdict = append(verdict, fmt.Sprintf("Missing action on defection: %v", &n))
			}
			for _, u := range n.OnDefection {
				if u < 0 || int(u) >= len(fsm.Node) || !fsm.Node[u].Active {
					verdict = append(verdict, fmt.Sprintf("Illegal on defection action from node %v", &n))
				}
			}
		}
	}
	// Check duplicate transitions
	//for _, n := range fsm.Node {
	//	for t1 := 0; t1 < len(n.OnCooperation); t1++ {
	//		for t2 := t1 + 1; t2 < len(n.OnCooperation); t2++ {
	//			if n.OnCooperation[t1] == n.OnCooperation[t2] {
	//				verdict = append(verdict, fmt.Sprintf("Duplicate on cooperation action in node %v", &n))
	//			}
	//		}
	//	}
	//	for t1 := 0; t1 < len(n.OnDefection); t1++ {
	//		for t2 := t1 + 1; t2 < len(n.OnDefection); t2++ {
	//			if n.OnDefection[t1] == n.OnDefection[t2] {
	//				verdict = append(verdict, fmt.Sprintf("Duplicate on defection action in node %v", &n))
	//			}
	//		}
	//	}
	//}

	// Check sink
	sink := fsm.GetSink()
	if len(fsm.Node[sink].OnCooperation)+len(fsm.Node[sink].OnCooperation) != 2 {
		verdict = append(verdict, fmt.Sprintf("Illegal sink: wrong number of transitions %v", fsm.Node[sink]))
	} else {
		if fsm.Node[sink].OnCooperation[0] != NodeIndex(sink) {
			verdict = append(verdict, fmt.Sprintf("Illegal sink: OnCooperation %v", fsm.Node[sink]))
		}
		if fsm.Node[sink].OnDefection[0] != NodeIndex(sink) {
			verdict = append(verdict, fmt.Sprintf("Illegal sink: OnDefection %v", fsm.Node[sink]))
		}
	}
	if !fsm.Node[sink].Active {
		verdict = append(verdict, fmt.Sprintf("Illegal sink: not Active %v", fsm.Node[sink]))
	}

	if len(verdict) > 0 {
		log.Println("Huston, we have a problem")
		fsm.Debug()
		for _, s := range verdict {
			log.Println(s)
		}
		log.Panicln("Check failed")
	}
	return true
}

//////////////////////////////////////////////////////////////////////////////
// I/O

func LoadNdFsms(dbPath string) []*NdFsm {
	files, _ := ioutil.ReadDir(dbPath)
	opp := make([]*NdFsm, len(files))

	num := 0
	for t := 0; t < len(files); t++ {
		if filepath.Ext(os.FileInfo(files[t]).Name()) == ".json" {
			o := LoadNdFsm(dbPath + "/" + os.FileInfo(files[t]).Name())

			base := os.FileInfo(files[t]).Name()[0 : len(os.FileInfo(files[t]).Name())-5]
			if _, err := os.Stat(dbPath + "/" + base + ".gv"); err != nil {
				ioutil.WriteFile(dbPath+"/"+base+".json", o.JsonEncode(), 0644)
				ioutil.WriteFile(dbPath+"/"+base+".gv", o.GvEncode(), 0644)
				Gv2Png(dbPath + "/" + base)
			}
			if _, err := os.Stat(dbPath + "/" + base + ".oy.txt"); err != nil {
				if oy := o.OyunEncode(); oy != nil {
					ioutil.WriteFile(dbPath+"/"+base+".oy.txt", oy, 0644)
				}
			}

			opp[num] = o
			num++
		}
	}
	opp = opp[:num]
	return opp
}

func LoadNdFsm(file string) *NdFsm {
	blob, _ := ioutil.ReadFile(file)
	o := new(NdFsm)
	o.JsonDecode(blob)
	o.UpdateInternals()
	return o
}

func (fsm *NdFsm) JsonEncode() []byte {
	blob, err := json.MarshalIndent(fsm, "", "    ")
	if err != nil {
		log.Panicf("JSON Error: %s\n", err)
	}
	return append(blob, '\n')
}
func (fsm *NdFsm) JsonDecode(blob []byte) {
	err := json.Unmarshal(blob, fsm)
	if err != nil {
		log.Panicf("JSON Error: %v", err)
	}
}

func (fsm *NdFsm) Signature() string {
	var sig []string
	for _, i := range fsm.InitialState {
		sig = append(sig, fmt.Sprintf("I%d", fsm.Node[i].canonic))
	}
	for _, t := range fsm.ActiveNodes() {
		node := fmt.Sprintf("N%d", fsm.Node[t].canonic)
		if fsm.Node[t].Action == COOPERATE {
			node += "C"
		} else if fsm.Node[t].Action == DEFECT {
			node += "D"
		} else {
			node += "X"
		}
		sig = append(sig, node)
	}
	for _, n := range fsm.ActiveNodes() {
		for _, u := range fsm.Node[n].OnDefection {
			sig = append(sig, fmt.Sprintf("T%dc%d", fsm.Node[n].canonic, fsm.Node[u].canonic))
		}
		for _, u := range fsm.Node[n].OnCooperation {
			sig = append(sig, fmt.Sprintf("T%dd%d", fsm.Node[n].canonic, fsm.Node[u].canonic))
		}
	}
	sort.Strings(sig)
	return fmt.Sprint(sig)
}

//////////////////////////////////////////////////////////////////////////////
// cool Gv representation
// Cooperation -> double circle
// Defection -> single circle
// Non active -> grayed

func (fsm *NdFsm) GvEncode() []byte {
	gv := fmt.Sprintf("digraph finite_state_machine {\nlabel=\"%s\";\n", fsm.Name)
	gv += " i [ shape = none; label = \"\"]\n"
	for t := range fsm.Node {
		dump := true
		extra := ""
		shape := ""

		if !fsm.Node[t].Active {
			dump = false
		}
		if fsm.Node[t].Action == INVALID && !fsm.Node[t].reachable {
			dump = false
		}

		if fsm.Node[t].Action == COOPERATE {
			shape = "doublecircle"
		} else if fsm.Node[t].Action == DEFECT {
			shape = "circle"
		} else {
			shape = "square, distortion = 0"
		}
		if *OptionColors && fsm.Node[t].color != "" {
			shape += ", style = filled, color=" + fsm.Node[t].color
		}
		if dump {
			//gv += fmt.Sprintf(" N%d [ label = \"%d/%d\"; shape = %v; %v ]\n", t, t, fsm.Node[t].canonic, shape, extra)
			if fsm.Node[t].canonic >= 0 {
				gv += fmt.Sprintf(" N%d [ label = \"%d\"; shape = %v; %v ]\n", t, fsm.Node[t].canonic, shape, extra)
			} else {
				gv += fmt.Sprintf(" N%d [ label = \"\"; shape = %v; %v ]\n", t, shape, extra)
			}
		}
	}
	for _, i := range fsm.InitialState {
		gv += fmt.Sprintf("i -> N%d [style = bold]\n", i)
	}
	for n := range fsm.Node {
		if fsm.Node[n].Active && fsm.Node[n].Action != INVALID {
			style := ""
			if !fsm.Node[n].reachable {
				style = "style=dotted"
			}
			for _, u := range fsm.Node[n].OnDefection {
				gv += fmt.Sprintf("N%d -> N%d [label = \"D\" %s ]\n", n, u, style)
			}
			for _, u := range fsm.Node[n].OnCooperation {
				gv += fmt.Sprintf("N%d -> N%d [label = \"C\" %s ]\n", n, u, style)
			}
		}
	}
	gv += fmt.Sprint("}\n")
	return []byte(gv)
}

func (fsm *NdFsm) __GvEncode() []byte {
	gv := fmt.Sprintf("digraph finite_state_machine {\nlabel=\"%s\";\n", fsm.Name)
	gv += " i [ shape = none; label = \"\"]\n"
	for t := range fsm.Node {
		dump := true
		extra := ""
		shape := ""

		if !fsm.Node[t].Active {
			dump = false
		}
		if fsm.Node[t].Action == INVALID && !fsm.Node[t].reachable {
			dump = false
		}

		if fsm.Node[t].Action == COOPERATE {
			shape = "doublecircle"
		} else if fsm.Node[t].Action == DEFECT {
			shape = "circle"
		} else {
			shape = "square, distortion = 0, style = filled"
		}

		if !fsm.Node[t].reachable {
			extra = "color = gray; "
		}

		if dump {
			//gv += fmt.Sprintf(" N%d [ label = \"%d/%d\"; shape = %v; %v ]\n", t, t, fsm.Node[t].canonic, shape, extra)
			if fsm.Node[t].canonic >= 0 {
				gv += fmt.Sprintf(" N%d [ label = \"%d\"; shape = %v; %v ]\n", t, fsm.Node[t].canonic, shape, extra)
			} else {
				gv += fmt.Sprintf(" N%d [ label = \"\"; shape = %v; %v ]\n", t, shape, extra)
			}
		}
	}
	for _, i := range fsm.InitialState {
		gv += fmt.Sprintf("i -> N%d [style = bold]\n", i)
	}
	for n := range fsm.Node {
		if fsm.Node[n].Active && fsm.Node[n].Action != INVALID {
			style := ""
			if !fsm.Node[n].reachable {
				style = "style=dotted"
			}
			for _, u := range fsm.Node[n].OnDefection {
				gv += fmt.Sprintf("N%d -> N%d [label = \"D\" %s ]\n", n, u, style)
			}
			for _, u := range fsm.Node[n].OnCooperation {
				gv += fmt.Sprintf("N%d -> N%d [label = \"C\" %s ]\n", n, u, style)
			}
		}
	}
	gv += fmt.Sprint("}\n")
	return []byte(gv)
}

func (fsm *NdFsm) OyunEncode() []byte {
	oy := fmt.Sprintf("Tusna OyunEncode()\n%s\n%d\n", fsm.Name, len(fsm.Node))
	for n := range fsm.Node {
		if len(fsm.Node[n].OnCooperation)+len(fsm.Node[n].OnCooperation) > 2 {
			return nil
		} else if !fsm.Node[n].Active {
			oy += fmt.Sprintf("*, %d, %d\n", n, n)
		} else {
			oy += fmt.Sprintf("%s, %d, %d\n", fsm.Node[n].Action.Mark(), fsm.Node[n].OnCooperation[0], fsm.Node[n].OnDefection[0])
		}
	}
	return []byte(oy)
}

//////////////////////////////////////////////////////////////////////////////
// Pointers to DFSM satisfy the "Player" interface
// Notez bien: *DFSM, not DFSM

func (fsm *NdFsm) GetName() string {
	return fsm.Name
}

func (fsm *NdFsm) FirstPly() PPly {
	fsm.currentState = fsm.InitialState[rand.Intn(len(fsm.InitialState))]
	return fsm.Node[int(fsm.currentState)].Action
}

func (fsm *NdFsm) RePly(oppPly PPly) PPly {
	var alt []NodeIndex
	if oppPly == DEFECT {
		alt = fsm.Node[fsm.currentState].OnDefection
	} else {
		alt = fsm.Node[fsm.currentState].OnCooperation
	}
	if len(alt) < 1 {
		fsm.Debug()
		log.Panicln("No move on ", oppPly, ": ", alt)
	}
	fsm.currentState = alt[rand.Intn(len(alt))]
	return fsm.Node[int(fsm.currentState)].Action
}

// set a state (and check its validity)
func (fsm *NdFsm) SetState(s NodeIndex) PPly {
	fsm.currentState = s
	fsm.Check()
	return fsm.Node[int(fsm.currentState)].Action
}

//////////////////////////////////////////////////////////////////////////////
// expected payoff playing against the opponent
// if starting from CurrentState & following path
func (fsm *NdFsm) ExpectedPayoff(start map[NodeIndex]float64, path GameTrace) float64 {
	ExpectedGlobal := 0.0
	CurrentState := start

	for _, step := range path {
		// payoff
		ExpectedStep := 0.0
		TotStates := 0.0
		for s, m := range CurrentState {
			TotStates += m
			ExpectedStep += float64(Payoff[step][fsm.Node[s].Action]) * m
		}
		ExpectedGlobal += ExpectedStep / TotStates

		// next state
		NextState := make(map[NodeIndex]float64)
		for s, m := range CurrentState {
			next := fsm.Node[s].NextState(step)
			for _, n := range next {
				NextState[n] += m
			}
		}
		CurrentState = NextState
	}

	return ExpectedGlobal
}

//////////////////////////////////////////////////////////////////////////////
// likelihood to find the footprints on the given path

func (fsm *NdFsm) removedLikelihood(path, footprints GameTrace) float64 {
	// log.Printf("Evaluating likelihood(%v, %v)\n", path, footprint)
	// Sanity check
	if len(path) != len(footprints) {
		log.Fatalf("Inconsistent footprints: %v (len %d) on path %v (len %d)", footprints, len(footprints), path, len(path))
	}

	c := 0

	for _, s := range fsm.AllFootprints(path) {
		if s.Compatible(footprints) {
			c++
		}
	}
	return float64(c) / float64(len(_fps))
}

//////////////////////////////////////////////////////////////////////////////
// *very* simple recursive function

var _fps []GameTrace

func (fsm *NdFsm) AllFootprints(path GameTrace) []GameTrace {
	_fps = make([]GameTrace, 0, 32)
	for _, i := range fsm.InitialState {
		fsm.explorePath(path, nil, 0, i)
	}
	return _fps
}

func (fsm *NdFsm) explorePath(path, currentFootprints GameTrace, step int, currentState NodeIndex) {
	if step == len(path) {
		n := make(GameTrace, len(currentFootprints))
		copy(n, currentFootprints)
		_fps = append(_fps, n)
		return
	}

	currentFootprints = append(currentFootprints, fsm.Node[currentState].Action)

	var alt []NodeIndex
	if path[step] == DEFECT {
		alt = fsm.Node[currentState].OnDefection
	} else {
		alt = fsm.Node[currentState].OnCooperation
	}
	for t := range alt {
		fsm.explorePath(path, currentFootprints, step+1, alt[t])
	}
}
