// Pattern 25: Real-world scenario - Event handling
component EventBus {
    depends on EventStore
    depends on SubscriberRegistry
    depends on DeadLetterQueue
    depends on MetricsCollector

    flow PublishEvent {
        EventStore.persist(event)
        MetricsCollector.recordEvent(event)
        subscribers = SubscriberRegistry.getSubscribers(eventType)
        for subscriber in subscribers {
            delivered = self.deliverToSubscriber(subscriber, event)
            if delivered == false {
                DeadLetterQueue.enqueue(event, subscriber)
            }
        }
        MetricsCollector.recordDelivery(event)
        return true
    }

    flow ConsumeEvent {
        event = EventStore.getNext(consumerId)
        if event != null {
            processed = self.processEvent(event)
            if processed {
                EventStore.acknowledge(event)
                MetricsCollector.recordProcessed(event)
            } else {
                EventStore.nack(event)
                MetricsCollector.recordFailed(event)
            }
        }
        return event
    }
}

component EventStore { }
component SubscriberRegistry { }
component DeadLetterQueue { }
component MetricsCollector { }
