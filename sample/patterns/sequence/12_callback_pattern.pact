// Pattern 12: Callback pattern (A->B->A)
component EventPublisher {
    depends on EventSubscriber

    flow PublishWithCallback {
        EventSubscriber.onEventReceived(event)
        self.confirmDelivery(event)
        return true
    }
}

component EventSubscriber {
    depends on EventPublisher

    flow OnEventReceived {
        processed = self.handleEvent(event)
        EventPublisher.acknowledge(event)
        return processed
    }
}
