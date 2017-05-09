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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

//////////////////////////////////////////////////////////////////////////////
// Basic structures

// DFsmNode is a deterministic FSM node
type DFsmNode struct {
	Active    bool
	Action    PPly
	NextState [3]NodeIndex
}

// DFsm is a deterministic FSM
type DFsm struct {
	Name         string
	InitialState NodeIndex
	currentState NodeIndex
	Node         []DFsmNode
}

//////////////////////////////////////////////////////////////////////////////
// Stringer interfaces

func (n *DFsmNode) String() string {
	var m = [...]byte{'D', 'C'}
	var s string
	if n.Active {
		s = fmt.Sprintf("{%c;%d/%d}", m[n.Action], n.NextState[0], n.NextState[1])
	} else {
		s = "{-;-/-}"
	}
	return s
}

func (f *DFsm) String() string {
	s := fmt.Sprintf("{i:%d [", f.InitialState)
	for t := 0; t < len(f.Node); t++ {
		if t == int(f.currentState) {
			s += " " + fmt.Sprint(t) + "*"
		} else {
			s += " " + fmt.Sprint(t)
		}
		s += ":" + f.Node[t].String()
	}
	s += " ]}"
	return s
}

//////////////////////////////////////////////////////////////////////////////
// Self explanatory

// NewDFSM creates a deterministic FSM with a given name
func NewDFSM(name string) DFsm {
	f := *new(DFsm)
	f.Name = name
	f.currentState = f.InitialState
	f.Node = make([]DFsmNode, 1)
	f.Node[0].Active = false
	f.Check()
	return f
}

func (f *DFsm) AddNode(action PPly) NodeIndex {
	found := -1
	for t := 0; t < len(f.Node); t++ {
		if !f.Node[t].Active {
			found = t
			break
		}
	}
	if found == -1 {
		var n DFsmNode
		f.Node = append(f.Node, n)
		found = len(f.Node) - 1
	}

	f.Node[found].Active = true
	f.Node[found].Action = action
	f.Node[found].NextState[COOPERATE] = NodeIndex(found)
	f.Node[found].NextState[DEFECT] = NodeIndex(found)

	f.Check()
	return NodeIndex(found)
}

func (f *DFsm) RemoveNode(nodeR NodeIndex) []NodeIndex {
	ret := []NodeIndex{}
	for t := 0; t < len(f.Node); t++ {
		if NodeIndex(t) == nodeR {
			f.Node[t].Active = false
		}
		if f.Node[t].NextState[COOPERATE] == nodeR {
			f.Node[t].NextState[COOPERATE] = NodeIndex(t)
			ret = append(ret, NodeIndex(t))
		}
		if f.Node[t].NextState[DEFECT] == nodeR {
			f.Node[t].NextState[DEFECT] = NodeIndex(t)
			ret = append(ret, NodeIndex(t))
		}
	}

	f.Check()
	return ret
}

func (f *DFsm) ModifyNodeAction(nodeCurrent NodeIndex, action PPly) {
	action.Check()
	f.Node[nodeCurrent].Action = action
	f.Check()
}

func (f *DFsm) SetTransition(nodeCurrent NodeIndex, action PPly, nodeTarget NodeIndex) {
	action.Check()
	f.Node[nodeCurrent].NextState[action] = nodeTarget
	f.Check()
}

//////////////////////////////////////////////////////////////////////////////
// Check if valid
// Rationale: quit asap, don't propagate errors

func (f *DFsm) Check() bool {
	if f.InitialState < 0 || int(f.InitialState) >= len(f.Node) {
		log.Panicf("Invalid initial state: %s\n", f)
	}
	if f.currentState != -1 {
		if f.currentState < 0 || int(f.currentState) >= len(f.Node) || !f.Node[f.currentState].Active {
			log.Panicf("Invalid current state: %s\n", f)
		}
	}
	for t := 0; t < len(f.Node); t++ {
		if f.Node[t].Active {
			f.Node[t].Action.Check()
			if f.Node[t].NextState[0] < 0 || int(f.Node[t].NextState[0]) >= len(f.Node) {
				log.Panicf("Invalid index: %v\n", f.Node[t])
			}
			if !f.Node[f.Node[t].NextState[DEFECT]].Active {
				log.Panicf("Illegal action on DEFECT: %v\n", f.Node[t])
			}
			if !f.Node[f.Node[t].NextState[1]].Active {
				log.Panicf("Illegal action on COOPERATE: %v\n", f.Node[t])
			}
		}
	}
	return true
}

//////////////////////////////////////////////////////////////////////////////
// I/O

func LoadDFsm(dbPath string) []*DFsm {
	files, _ := ioutil.ReadDir(dbPath)
	opp := make([]*DFsm, len(files))
	num := 0
	for t := 0; t < len(files); t++ {
		if filepath.Ext(os.FileInfo(files[t]).Name()) == ".json" {
			base := os.FileInfo(files[t]).Name()[0 : len(os.FileInfo(files[t]).Name())-5]
			blob, _ := ioutil.ReadFile(dbPath + "/" + base + ".json")
			o := new(DFsm)
			o.JsonDecode(blob)
			_, err := os.Stat(dbPath + "/" + base + ".gv")
			if err != nil {
				ioutil.WriteFile(dbPath+"/"+base+".gv", o.GvEncode(), 0644)
			}
			log.Printf("Loaded [%d] \"%s\"\n", num, o.Name)

			opp[num] = o
			num++
			// ioutil.WriteFile(dbPath + "/"+base+".json", opp[num].JsonEncode(), 0644)
		}
	}
	opp = opp[:num]
	return opp
}

func (f *DFsm) JsonEncode() []byte {
	blob, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Panicf("JSON Error: %s\n", err)
	}
	return append(blob, '\n')
}
func (f *DFsm) JsonDecode(blob []byte) {
	err := json.Unmarshal(blob, f)
	if err != nil {
		log.Panicf("JSON Error: %v", err)
	}
}
func (f *DFsm) GvEncode() []byte {
	gv := fmt.Sprintf("digraph finite_state_machine {\nlabel=\"%s\";\n", f.Name)
	gv += " i [ shape = none; label = \"\"]\n"
	for t := 0; t < len(f.Node); t++ {
		gv += fmt.Sprintf(" N%d [ label = %d; fixedsize = true; ", t, t)
		if !f.Node[t].Active {
			gv += "color = gray; "
		}
		if f.Node[t].Action == COOPERATE {
			gv += "shape = doublecircle; "
		} else {
			gv += "shape = circle; "
		}
		gv += "]\n"
	}
	gv += fmt.Sprintf("i -> N%d [style = bold]\n", f.InitialState)
	for t := 0; t < len(f.Node); t++ {
		if f.Node[t].NextState[0] == f.Node[t].NextState[1] {
			gv += fmt.Sprintf("N%d -> N%d [label = \"*\"]\n", t, f.Node[t].NextState[0])
		} else {
			gv += fmt.Sprintf("N%d -> N%d [label = \"D\"]\n", t, f.Node[t].NextState[0])
			gv += fmt.Sprintf("N%d -> N%d [label = \"C\"]\n", t, f.Node[t].NextState[1])
		}
	}
	gv += fmt.Sprint("}\n")
	return []byte(gv)
}

//////////////////////////////////////////////////////////////////////////////
// Pointers to DFsm satisfy the "Player" interface
// Notez bien: *DFsm, not DFsm

func (fsm *DFsm) GetName() string {
	return fsm.Name
}

func (fsm *DFsm) FirstPly() PPly {
	fsm.currentState = fsm.InitialState
	return fsm.Node[int(fsm.currentState)].Action
}

func (fsm *DFsm) RePly(oppPly PPly) PPly {
	fsm.currentState = fsm.Node[fsm.currentState].NextState[oppPly]
	return fsm.Node[int(fsm.currentState)].Action
}

//////////////////////////////////////////////////////////////////////////////
// footprints

func (fsm *DFsm) Footprints(path GameTrace) GameTrace {
	fp := make(GameTrace, 0, 100)
	cur := fsm.InitialState
	for _, om := range path {
		fp = append(fp, fsm.Node[cur].Action)
		cur = fsm.Node[fsm.currentState].NextState[om]
	}
	return fp
}

func (fsm *DFsm) removedLikelihood(path, footprint GameTrace) float64 {
	// Sanity check
	if len(path) != len(footprint) {
		log.Fatalf("Inconsistent footprint: %v (len %d) vs. %v (len %d)", footprint, len(footprint), path, len(path))
	}

	if footprint.Compatible(fsm.Footprints(path)) {
		return 1.0
	} else {
		return 0.0
	}
}
