// Real-world Pattern 5: Message Queue System
// Demonstrates: Async processing, Message states, Producer-Consumer pattern

// ============ TYPE DEFINITIONS ============

component Message {
	type MessageEnvelope {
		id: string
		correlationId: string?
		messageType: string
		payload: string
		headers: string[]
		priority: int
		ttl: int?
		createdAt: string
		scheduledFor: string?
	}

	type MessageMetadata {
		id: string
		attempts: int
		maxAttempts: int
		lastAttemptAt: string?
		nextRetryAt: string?
		errorMessage: string?
		processingTime: int?
	}

	type DeadLetterEntry {
		messageId: string
		queueName: string
		reason: string
		originalPayload: string
		failedAt: string
		retryable: bool
	}

	// Message lifecycle states
	states MessageLifecycle {
		initial Pending

		state Pending {
			entry [recordEnqueue]
		}
		state Scheduled {
			entry [setTimer]
		}
		state Processing {
			entry [lockMessage]
		}
		state Completed {
			entry [recordSuccess]
		}
		state Failed {
			entry [recordFailure]
		}
		state Retrying {
			entry [calculateBackoff]
		}
		state DeadLetter {
			entry [moveToDeadLetter]
		}

		Pending -> Processing on consumerAcquired
		Pending -> Scheduled on delayedMessage
		Scheduled -> Pending on scheduleReached
		Scheduled -> DeadLetter on ttlExpired
		Processing -> Completed on processedSuccessfully
		Processing -> Failed on processingError
		Processing -> Completed on acknowledgeReceived
		Failed -> Retrying on retryScheduled
		Failed -> DeadLetter on maxRetriesExceeded
		Retrying -> Pending on retryReady
		Retrying -> DeadLetter on retryExpired
		DeadLetter -> Pending on manualRequeue
	}
}

component Queue {
	type QueueConfig {
		name: string
		maxSize: int
		maxRetries: int
		retryDelayMs: int
		visibilityTimeout: int
		deadLetterQueue: string?
		fifo: bool
	}

	type QueueStats {
		name: string
		size: int
		inFlight: int
		completed: int
		failed: int
		deadLettered: int
		avgProcessingTime: float
	}

	type QueueMessage {
		message: MessageEnvelope
		metadata: MessageMetadata
		receiptHandle: string?
	}

	provides QueueService {
		CreateQueue(config: QueueConfig) -> bool
		DeleteQueue(queueName: string) -> bool
		GetQueueStats(queueName: string) -> QueueStats
		PurgeQueue(queueName: string) -> int
		ListQueues() -> QueueConfig[]
		Enqueue(queueName: string, message: MessageEnvelope) -> string
		Dequeue(queueName: string, count: int) -> QueueMessage[]
		Acknowledge(queueName: string, receiptHandle: string) -> bool
		Reject(queueName: string, receiptHandle: string, requeue: bool) -> bool
		ExtendVisibility(receiptHandle: string, seconds: int) -> bool
	}

	requires StorageBackend {
		Put(key: string, value: string) -> bool
		Get(key: string) -> string?
		Delete(key: string) -> bool
		List(prefix: string) -> string[]
		Lock(key: string, ttl: int) -> bool
		Unlock(key: string) -> bool
	}
}

component Producer {
	type PublishOptions {
		priority: int?
		delay: int?
		ttl: int?
		correlationId: string?
		headers: string[]
	}

	type PublishResult {
		messageId: string
		queueName: string
		timestamp: string
		scheduled: bool
	}

	type BatchPublishResult {
		successful: PublishResult[]
		failed: string[]
	}

	depends on Queue
	depends on MessageRouter

	provides ProducerAPI {
		Publish(queueName: string, message: MessageEnvelope) -> PublishResult
		PublishWithOptions(queueName: string, message: MessageEnvelope, options: PublishOptions) -> PublishResult
		PublishBatch(queueName: string, messages: MessageEnvelope[]) -> BatchPublishResult
		PublishToTopic(topic: string, message: MessageEnvelope) -> PublishResult[]
		ScheduleMessage(queueName: string, message: MessageEnvelope, scheduleAt: string) -> PublishResult
	}

	flow PublishMessage {
		validated = self.validateMessage(message)
		if validated {
			enriched = self.enrichMessage(message)
			if hasTopic {
				queues = MessageRouter.GetQueuesForTopic(topic)
				results = self.createEmptyResults()
				for queue in queues {
					result = Queue.Enqueue(queue.name, enriched)
					results = self.appendResult(results, result)
				}
				return results
			} else {
				queueExists = Queue.GetQueueStats(queueName)
				if queueExists {
					if hasDelay {
						enriched = self.setScheduledTime(enriched, delay)
					}
					messageId = Queue.Enqueue(queueName, enriched)
					self.recordPublishMetrics(queueName, enriched)
					return PublishResult
				} else {
					throw QueueNotFoundError
				}
			}
		} else {
			throw InvalidMessageError
		}
	}
}

component Consumer {
	type ConsumerConfig {
		queueName: string
		batchSize: int
		pollInterval: int
		concurrency: int
		autoAcknowledge: bool
		visibilityTimeout: int
	}

	type ProcessingContext {
		message: QueueMessage
		attempt: int
		startedAt: string
		consumerGroup: string?
	}

	type ConsumerStats {
		consumerId: string
		queueName: string
		processed: int
		failed: int
		avgLatency: float
		lastActivity: string
	}

	depends on Queue
	depends on DeadLetterHandler
	depends on MetricsCollector

	provides ConsumerAPI {
		StartConsumer(config: ConsumerConfig) -> string
		StopConsumer(consumerId: string) -> bool
		PauseConsumer(consumerId: string) -> bool
		ResumeConsumer(consumerId: string) -> bool
		GetConsumerStats(consumerId: string) -> ConsumerStats
	}

	requires MessageHandler {
		Process(context: ProcessingContext) -> bool
		OnError(context: ProcessingContext, error: string)
		OnComplete(context: ProcessingContext)
	}

	flow ConsumeMessages {
		config = self.getConsumerConfig(consumerId)
		while consumerActive {
			messages = Queue.Dequeue(config.queueName, config.batchSize)
			if hasMessages {
				for message in messages {
					context = self.buildContext(message)
					MetricsCollector.RecordDequeue(config.queueName)
					processed = MessageHandler.Process(context)
					if processed {
						Queue.Acknowledge(config.queueName, message.receiptHandle)
						MessageHandler.OnComplete(context)
						MetricsCollector.RecordSuccess(config.queueName, processingTime)
					} else {
						if canRetry {
							Queue.Reject(config.queueName, message.receiptHandle, true)
							MetricsCollector.RecordRetry(config.queueName)
						} else {
							Queue.Reject(config.queueName, message.receiptHandle, false)
							DeadLetterHandler.HandleFailure(message, lastError)
							MessageHandler.OnError(context, lastError)
							MetricsCollector.RecordFailure(config.queueName)
						}
					}
				}
			} else {
				self.sleep(config.pollInterval)
			}
		}
	}
}

component DeadLetterHandler {
	depends on Queue
	depends on AlertService

	provides DeadLetterAPI {
		HandleFailure(message: QueueMessage, error: string) -> DeadLetterEntry
		GetDeadLetters(queueName: string, limit: int) -> DeadLetterEntry[]
		RequeueMessage(messageId: string) -> bool
		DeleteDeadLetter(messageId: string) -> bool
		RequeueAll(queueName: string) -> int
	}

	flow HandleDeadLetter {
		deadLetterRecord = self.createDeadLetterEntry(message, error)
		dlqName = self.getDeadLetterQueue(message.queueName)
		Queue.Enqueue(dlqName, message.message)
		self.recordDeadLetter(deadLetterRecord)
		threshold = self.getAlertThreshold(message.queueName)
		count = self.getDeadLetterCount(dlqName)
		if countExceedsThreshold {
			AlertService.SendAlert(highDeadLetterCountAlert)
		}
		return deadLetterRecord
	}
}

component MessageRouter {
	type TopicSubscription {
		topic: string
		queueName: string
		filter: string?
		createdAt: string
	}

	type RoutingRule {
		id: string
		sourcePattern: string
		destinationQueue: string
		condition: string?
		transform: string?
		priority: int
	}

	provides RouterAPI {
		Subscribe(topic: string, queueName: string, filter: string?) -> TopicSubscription
		Unsubscribe(topic: string, queueName: string) -> bool
		GetQueuesForTopic(topic: string) -> QueueConfig[]
		AddRoutingRule(rule: RoutingRule) -> bool
		RemoveRoutingRule(ruleId: string) -> bool
		RouteMessage(message: MessageEnvelope) -> string[]
	}
}

component MetricsCollector {
	type QueueMetrics {
		queueName: string
		enqueued: int
		dequeued: int
		processed: int
		failed: int
		retried: int
		deadLettered: int
		avgLatency: float
		p99Latency: float
		timestamp: string
	}

	provides MetricsAPI {
		RecordEnqueue(queueName: string)
		RecordDequeue(queueName: string)
		RecordSuccess(queueName: string, latencyMs: int)
		RecordFailure(queueName: string)
		RecordRetry(queueName: string)
		GetMetrics(queueName: string, from: string, to: string) -> QueueMetrics[]
		GetDashboard() -> QueueMetrics[]
	}

	requires MetricsStore {
		Write(metric: string, value: float, tags: string[])
		Query(metric: string, from: string, to: string) -> float[]
	}
}

component AlertService {
	type Alert {
		id: string
		severity: string
		source: string
		message: string
		timestamp: string
		acknowledged: bool
	}

	provides AlertAPI {
		SendAlert(alert: Alert) -> bool
		GetActiveAlerts() -> Alert[]
		AcknowledgeAlert(alertId: string) -> bool
	}

	requires NotificationChannel {
		SendEmail(to: string, subject: string, body: string) -> bool
		SendSlack(channel: string, message: string) -> bool
		SendPagerDuty(severity: string, message: string) -> bool
	}
}
