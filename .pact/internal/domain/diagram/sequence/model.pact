// Sequence diagram domain model
component SequenceDiagram {
    // ParticipantType represents the type of participant
    enum ParticipantType {
        DefaultParticipant
        ActorParticipant
        DatabaseParticipant
        QueueParticipant
        ExternalParticipant
    }

    // MessageType represents the type of message
    enum MessageType {
        SyncMessage
        AsyncMessage
        ReturnMessage
    }

    // FragmentType represents the type of fragment
    enum FragmentType {
        alt
        loop
        opt
    }

    // Participant represents a participant in the sequence diagram
    type Participant {
        id: string
        name: string
        participantType: ParticipantType
    }

    // MessageEvent represents a message event between participants
    type MessageEvent {
        fromParticipant: string
        toParticipant: string
        label: string
        messageType: MessageType
    }

    // FragmentEvent represents a fragment event (alt, loop, opt)
    type FragmentEvent {
        fragmentType: FragmentType
        label: string
        altLabel: string
    }

    // ActivationEvent represents an activation/deactivation event
    type ActivationEvent {
        participant: string
        active: bool
    }

    // Diagram represents a sequence diagram
    type Diagram {
        participants: Participant[]
    }
}
