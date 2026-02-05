// Pattern: Callback
// コールバックを含むパターン

component EventEmitter {
    type Event {
        eventType: string
        data: string
        timestamp: string
    }

    depends on EventHandler

    flow EmitEvent {
        event = self.createEvent(eventType, data)
        self.logEvent(event)
        EventHandler.onEvent(event)
        confirmation = EventHandler.getConfirmation()
        self.recordDelivery(event, confirmation)
        return confirmation
    }
}

component EventHandler {
    type HandlerResult {
        handled: bool
        message: string
    }

    provides EventHandlerAPI {
        OnEvent(event: string) -> bool
        GetConfirmation() -> string
    }

    flow OnEvent {
        self.parseEvent(event)
        result = self.processEvent(event)
        self.notifyObservers(result)
        return result
    }
}
