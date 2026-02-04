package state

import "pact/internal/domain/diagram/common"

// Diagram は状態図を表す
type Diagram struct {
	States      []State
	Transitions []Transition
	Notes       []common.Note
}

func (d *Diagram) Type() common.DiagramType {
	return common.DiagramTypeState
}

// State は状態
type State struct {
	ID          string
	Name        string
	Type        StateType
	Entry       []string
	Exit        []string
	Children    []State  // 階層状態の場合
	Regions     []Region // 並行状態の場合
	Annotations []common.Annotation
}

// StateType は状態の種類
type StateType string

const (
	StateTypeInitial  StateType = "initial"
	StateTypeFinal    StateType = "final"
	StateTypeAtomic   StateType = "atomic"
	StateTypeCompound StateType = "compound"
	StateTypeParallel StateType = "parallel"
)

// Region は並行状態のリージョン
type Region struct {
	Name        string
	States      []State
	Transitions []Transition
}

// Transition は状態遷移
type Transition struct {
	From    string
	To      string
	Trigger Trigger
	Guard   string
	Actions []string
}

// Trigger はトリガー
type Trigger interface {
	triggerNode()
}

// EventTrigger はイベントトリガー
type EventTrigger struct {
	Event string
}

func (t *EventTrigger) triggerNode() {}

// AfterTrigger は時間トリガー
type AfterTrigger struct {
	Duration Duration
}

func (t *AfterTrigger) triggerNode() {}

// Duration は期間
type Duration struct {
	Value int
	Unit  string
}

// WhenTrigger は条件トリガー
type WhenTrigger struct {
	Condition string
}

func (t *WhenTrigger) triggerNode() {}
