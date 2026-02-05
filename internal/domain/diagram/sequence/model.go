package sequence

import "pact/internal/domain/diagram/common"

// Diagram はシーケンス図を表す
type Diagram struct {
	Participants []Participant
	Events       []Event
	Notes        []common.Note
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
	Type      FragmentType
	Label     string
	Events    []Event // メイン(then)部分のイベント
	AltLabel  string  // alt の else 部分用ラベル
	AltEvents []Event // alt の else 部分のイベント
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

// NoteEvent はシーケンス図内の注釈イベント
type NoteEvent struct {
	Participant string   // 注釈を付ける参加者（空の場合は全体）
	Text        string   // 注釈テキスト
	NoteType    NoteType // 注釈の種類
}

func (e *NoteEvent) eventNode() {}

// NoteType は注釈の種類
type NoteType string

const (
	NoteTypeNote   NoteType = "note"   // 通常の注釈
	NoteTypeReturn NoteType = "return" // return文
	NoteTypeThrow  NoteType = "throw"  // throw文
)
