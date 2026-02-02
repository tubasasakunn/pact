// Real-world Pattern 3: Microservices Architecture
// Demonstrates: Service dependencies, API Gateway pattern, Inter-service communication

// ============ API GATEWAY ============

component Gateway {
	type APIRequest {
		method: string
		path: string
		headers: string[]
		body: string?
		clientId: string
		timestamp: string
	}

	type APIResponse {
		statusCode: int
		headers: string[]
		body: string?
		latencyMs: int
	}

	type RateLimitInfo {
		clientId: string
		requestCount: int
		windowStart: string
		limit: int
		remaining: int
	}

	depends on UserService
	depends on OrderService
	depends on ProductService
	depends on NotificationService
	depends on AuthService

	provides GatewayAPI {
		HandleRequest(request: APIRequest) -> APIResponse
		HealthCheck() -> bool
		GetMetrics() -> string
	}

	requires LoadBalancer {
		GetNextInstance(serviceName: string) -> string
		ReportHealth(instanceId: string, healthy: bool)
	}

	requires CircuitBreaker {
		IsOpen(serviceName: string) -> bool
		RecordSuccess(serviceName: string)
		RecordFailure(serviceName: string)
	}

	flow HandleIncomingRequest {
		authenticated = AuthService.ValidateToken(request.token)
		if authenticated {
			rateLimited = self.checkRateLimit(request.clientId)
			if notRateLimited {
				circuitOpen = CircuitBreaker.IsOpen(targetService)
				if circuitOpen {
					throw ServiceUnavailableError
				}
				instance = LoadBalancer.GetNextInstance(targetService)
				response = self.forwardRequest(instance, request)
				if responseSuccessful {
					CircuitBreaker.RecordSuccess(targetService)
					self.updateMetrics(request, response)
					return response
				} else {
					CircuitBreaker.RecordFailure(targetService)
					throw ServiceError
				}
			} else {
				throw RateLimitExceededError
			}
		} else {
			throw UnauthorizedError
		}
	}
}

// ============ USER SERVICE ============

component UserService {
	type UserDTO {
		id: string
		email: string
		name: string
		createdAt: string
		preferences: string?
	}

	type CreateUserRequest {
		email: string
		name: string
		password: string
	}

	provides UserAPI {
		CreateUser(request: CreateUserRequest) -> UserDTO
		GetUser(userId: string) -> UserDTO
		UpdateUser(userId: string, updates: UserDTO) -> UserDTO
		DeleteUser(userId: string) -> bool
		SearchUsers(query: string) -> UserDTO[]
	}

	depends on NotificationService

	requires UserRepository {
		Save(user: UserDTO) -> UserDTO
		FindById(id: string) -> UserDTO?
		FindByEmail(email: string) -> UserDTO?
		Update(user: UserDTO) -> UserDTO
		Delete(id: string) -> bool
	}

	flow CreateNewUser {
		existing = UserRepository.FindByEmail(request.email)
		if userExists {
			throw UserAlreadyExistsError
		}
		user = self.buildUserEntity(request)
		saved = UserRepository.Save(user)
		NotificationService.SendWelcomeEmail(saved.email, saved.name)
		self.publishUserCreatedEvent(saved)
		return saved
	}
}

// ============ ORDER SERVICE ============

component OrderService {
	type OrderDTO {
		id: string
		userId: string
		items: OrderItemDTO[]
		status: string
		total: float
		createdAt: string
		updatedAt: string
	}

	type OrderItemDTO {
		productId: string
		quantity: int
		price: float
	}

	type CreateOrderRequest {
		userId: string
		items: OrderItemDTO[]
		shippingAddress: string
	}

	depends on UserService
	depends on ProductService
	depends on NotificationService
	depends on PaymentService

	provides OrderAPI {
		CreateOrder(request: CreateOrderRequest) -> OrderDTO
		GetOrder(orderId: string) -> OrderDTO
		GetUserOrders(userId: string) -> OrderDTO[]
		UpdateOrderStatus(orderId: string, status: string) -> OrderDTO
		CancelOrder(orderId: string) -> bool
	}

	requires OrderRepository {
		Save(order: OrderDTO) -> OrderDTO
		FindById(id: string) -> OrderDTO?
		FindByUserId(userId: string) -> OrderDTO[]
		Update(order: OrderDTO) -> OrderDTO
	}

	requires EventBus {
		Publish(topic: string, event: string)
		Subscribe(topic: string, handler: string)
	}

	// Order processing states
	states OrderProcessing {
		initial Received

		state Received {
			entry [validateOrder]
		}
		state Validating {
			entry [checkInventory]
		}
		state PaymentPending {
			entry [requestPayment]
		}
		state Confirmed {
			entry [reserveInventory]
		}
		state Fulfilling {
			entry [startFulfillment]
		}
		state Shipped {
			entry [updateTracking]
		}
		state Completed {
			entry [requestFeedback]
		}
		state Failed {
			entry [notifyFailure]
		}
		state Cancelled {
			entry [releaseReservations]
		}

		Received -> Validating on validationStarted
		Validating -> PaymentPending on validationPassed
		Validating -> Failed on validationFailed
		PaymentPending -> Confirmed on paymentSucceeded
		PaymentPending -> Failed on paymentFailed
		Confirmed -> Fulfilling on fulfillmentStarted
		Confirmed -> Cancelled on cancellationRequested
		Fulfilling -> Shipped on itemsShipped
		Shipped -> Completed on deliveryConfirmed
		Shipped -> Cancelled on returnInitiated
	}

	flow ProcessNewOrder {
		user = UserService.GetUser(request.userId)
		if userNotFound {
			throw UserNotFoundError
		}
		for item in request.items {
			product = ProductService.GetProduct(item.productId)
			available = ProductService.CheckAvailability(item.productId, item.quantity)
			if notAvailable {
				throw InsufficientStockError
			}
		}
		order = self.buildOrder(request, user)
		saved = OrderRepository.Save(order)
		payment = PaymentService.ProcessPayment(saved.id, saved.total)
		if paymentSuccessful {
			for item in request.items {
				ProductService.ReserveStock(item.productId, item.quantity)
			}
			updated = OrderRepository.Update(saved)
			EventBus.Publish(orderCreatedTopic, orderEvent)
			NotificationService.SendOrderConfirmation(user.email, saved)
			return updated
		} else {
			OrderRepository.Update(failedOrder)
			NotificationService.SendPaymentFailure(user.email, saved)
			throw PaymentFailedError
		}
	}
}

// ============ PRODUCT SERVICE ============

component ProductService {
	type ProductDTO {
		id: string
		name: string
		description: string
		price: float
		stock: int
		category: string
	}

	provides ProductAPI {
		GetProduct(productId: string) -> ProductDTO
		SearchProducts(query: string, category: string?) -> ProductDTO[]
		CheckAvailability(productId: string, quantity: int) -> bool
		ReserveStock(productId: string, quantity: int) -> bool
		ReleaseStock(productId: string, quantity: int) -> bool
		UpdateStock(productId: string, quantity: int) -> ProductDTO
	}

	requires ProductRepository {
		FindById(id: string) -> ProductDTO?
		Search(query: string, category: string?) -> ProductDTO[]
		Update(product: ProductDTO) -> ProductDTO
	}

	requires InventoryCache {
		Get(productId: string) -> int?
		Set(productId: string, quantity: int)
		Decrement(productId: string, amount: int) -> int
		Increment(productId: string, amount: int) -> int
	}
}

// ============ NOTIFICATION SERVICE ============

component NotificationService {
	type EmailMessage {
		to: string
		subject: string
		body: string
		templateId: string?
	}

	type PushNotification {
		userId: string
		title: string
		body: string
		data: string?
	}

	provides NotificationAPI {
		SendWelcomeEmail(email: string, name: string)
		SendOrderConfirmation(email: string, order: string)
		SendPaymentFailure(email: string, order: string)
		SendShippingUpdate(email: string, trackingInfo: string)
		SendPushNotification(notification: PushNotification)
	}

	requires EmailProvider {
		Send(message: EmailMessage) -> bool
		SendBatch(messages: EmailMessage[]) -> bool
	}

	requires PushProvider {
		Send(notification: PushNotification) -> bool
	}

	requires MessageQueue {
		Enqueue(queue: string, message: string)
		Dequeue(queue: string) -> string?
	}
}

// ============ SUPPORTING SERVICES ============

component AuthService {
	provides AuthAPI {
		ValidateToken(token: string) -> bool
		GetUserIdFromToken(token: string) -> string?
		RefreshToken(refreshToken: string) -> string
	}
}

component PaymentService {
	type PaymentResult {
		success: bool
		transactionId: string?
		errorCode: string?
		errorMessage: string?
	}

	provides PaymentAPI {
		ProcessPayment(orderId: string, amount: float) -> PaymentResult
		RefundPayment(transactionId: string, amount: float) -> PaymentResult
		GetPaymentStatus(transactionId: string) -> string
	}

	requires PaymentGateway {
		Charge(amount: float, paymentMethod: string) -> PaymentResult
		Refund(transactionId: string, amount: float) -> PaymentResult
	}
}
