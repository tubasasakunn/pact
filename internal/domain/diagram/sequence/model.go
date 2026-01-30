package sequence

import "pact/internal/domain/diagram/common"

// Diagram はシーケンス図を表す
type Diagram struct {
	Participants []Participant
	Events       []Event
}

func (d *Diagram) Type() common.DiagramType {
	return common.DiagramTypeSequence
}

// Participant は参加者
type Participant struct {
	ID   string
	Name string
	Type ParticipantType
}

// ParticipantType は参加者の種類
type ParticipantType string

const (
	ParticipantTypeDefault  ParticipantType = "default"
	ParticipantTypeActor    ParticipantType = "actor"
	ParticipantTypeDatabase ParticipantType = "database"
	ParticipantTypeQueue    ParticipantType = "queue"
	ParticipantTypeExternal ParticipantType = "external"
)

// Event はイベント
type Event interface {
	eventNode()
}

// MessageEvent はメッセージイベント
type MessageEvent struct {
	From        string
	To          string
	Label       string
	MessageType MessageType
}

func (e *MessageEvent) eventNode() {}

// MessageType はメッセージの種類
type MessageType string

const (
	MessageTypeSync   MessageType = "sync"
	MessageTypeAsync  MessageType = "async"
	MessageTypeReturn MessageType = "return"
)

// FragmentEvent はフラグメントイベント
type FragmentEvent struct {
	Type     FragmentType
	Label    string
	Events   []Event
	AltLabel string // alt の else 部分用
}

func (e *FragmentEvent) eventNode() {}

// FragmentType はフラグメントの種類
type FragmentType string

const (
	FragmentTypeAlt  FragmentType = "alt"
	FragmentTypeLoop FragmentType = "loop"
	FragmentTypeOpt  FragmentType = "opt"
)

// ActivationEvent はアクティベーションイベント
type ActivationEvent struct {
	Participant string
	Active      bool
}

func (e *ActivationEvent) eventNode() {}
