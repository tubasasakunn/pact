// Real-world example: Order Service
@version("1.0")
@module("orders")

type Money {
	amount: float
	currency: string
}

type Address {
	line1: string
	line2: string?
	city: string
	state: string
	postalCode: string
	country: string
}

type OrderItem {
	productId: string
	productName: string
	quantity: int
	unitPrice: Money
	discount: Money?
}

type PaymentInfo {
	method: string
	cardLast4: string?
	transactionId: string?
}

type ShippingInfo {
	carrier: string
	trackingNumber: string?
	estimatedDelivery: timestamp?
}

interface OrderRepository {
	method FindById(id: string): Order?
	method FindByCustomer(customerId: string): []Order
	method Save(order: Order): Order
	method UpdateStatus(orderId: string, status: string): void
}

interface PaymentGateway {
	method Authorize(amount: Money, paymentInfo: PaymentInfo): PaymentResult
	method Capture(transactionId: string): PaymentResult
	method Refund(transactionId: string, amount: Money): PaymentResult
}

interface InventoryService {
	method CheckAvailability(items: []OrderItem): AvailabilityResult
	method Reserve(orderId: string, items: []OrderItem): ReservationResult
	method Release(orderId: string): void
}

interface ShippingService {
	method CreateShipment(order: Order): ShipmentResult
	method GetTrackingInfo(trackingNumber: string): TrackingInfo
	method CancelShipment(trackingNumber: string): void
}

@aggregate
component Order {
	id: string
	customerId: string
	items: []OrderItem
	shippingAddress: Address
	billingAddress: Address
	subtotal: Money
	tax: Money
	shipping: Money
	discount: Money
	total: Money
	status: string
	paymentInfo: PaymentInfo?
	shippingInfo: ShippingInfo?
	notes: string?
	createdAt: timestamp
	updatedAt: timestamp

	method AddItem(item: OrderItem): void
	method RemoveItem(productId: string): void
	method UpdateQuantity(productId: string, quantity: int): void
	method ApplyDiscount(discount: Money): void
	method CalculateTotals(): void

	states OrderStatus {
		initial -> Draft

		Draft -> Submitted: submit
		Draft -> Abandoned: abandon

		Submitted -> PaymentPending: awaitPayment
		Submitted -> Cancelled: cancel

		PaymentPending -> PaymentAuthorized: authorizePayment
		PaymentPending -> PaymentFailed: paymentFail
		PaymentPending -> Cancelled: cancel

		PaymentAuthorized -> Processing: startProcessing
		PaymentAuthorized -> Cancelled: cancel

		PaymentFailed -> PaymentPending: retryPayment
		PaymentFailed -> Cancelled: cancel

		Processing -> ReadyToShip: packComplete
		Processing -> Cancelled: cancel

		ReadyToShip -> Shipped: ship
		ReadyToShip -> Cancelled: cancel

		Shipped -> OutForDelivery: outForDelivery
		Shipped -> Delivered: deliver
		Shipped -> DeliveryFailed: deliveryFail

		OutForDelivery -> Delivered: deliver
		OutForDelivery -> DeliveryFailed: deliveryFail

		DeliveryFailed -> Shipped: reship
		DeliveryFailed -> Returned: returnToSender

		Delivered -> Completed: complete
		Delivered -> ReturnRequested: requestReturn

		ReturnRequested -> ReturnApproved: approveReturn
		ReturnRequested -> ReturnDenied: denyReturn

		ReturnApproved -> Returned: receiveReturn

		Returned -> Refunded: processRefund

		Completed -> [*]
		Cancelled -> [*]
		Abandoned -> [*]
		Refunded -> [*]
		ReturnDenied -> [*]
	}
}

@service
component OrderService {
	private orderRepo: OrderRepository
	private paymentGateway: PaymentGateway
	private inventoryService: InventoryService
	private shippingService: ShippingService
	private notificationService: NotificationService
	private eventPublisher: EventPublisher

	relation OrderRepository: uses
	relation PaymentGateway: uses
	relation InventoryService: uses
	relation ShippingService: uses

	@transactional
	method CreateOrder(request: CreateOrderRequest): Order

	@transactional
	method SubmitOrder(orderId: string): Order

	method ProcessPayment(orderId: string): PaymentResult

	@transactional
	method CancelOrder(orderId: string, reason: string): void

	method ShipOrder(orderId: string): ShipmentResult

	method GetOrderStatus(orderId: string): OrderStatusResponse

	method RequestReturn(orderId: string, reason: string): ReturnRequest

	@transactional
	method ProcessRefund(orderId: string): RefundResult

	flow SubmitOrder {
		start: "Receive Order Submission"
		validateOrder: "Validate Order"
		if orderValid {
			checkInventory: "Check Inventory Availability"
			if inventoryAvailable {
				reserveInventory: "Reserve Inventory"
				if reservationSuccess {
					calculateTotals: "Calculate Order Totals"
					authorizePayment: "Authorize Payment"
					if paymentAuthorized {
						createOrder: "Create Order Record"
						publishEvent: "Publish OrderCreated Event"
						sendConfirmation: "Send Order Confirmation"
						success: "Order Submitted Successfully"
					} else {
						releaseInventory: "Release Inventory"
						paymentError: "Return Payment Error"
					}
				} else {
					inventoryError: "Return Inventory Error"
				}
			} else {
				outOfStock: "Return Out of Stock Error"
			}
		} else {
			validationError: "Return Validation Error"
		}
		end: "Complete"
	}

	flow ProcessReturn {
		start: "Receive Return Request"
		validateReturn: "Validate Return Eligibility"
		if eligible {
			createReturnLabel: "Create Return Shipping Label"
			notifyCustomer: "Notify Customer"
			waitForReturn: "Wait for Package Return"
			if packageReceived {
				inspectItems: "Inspect Returned Items"
				if itemsAcceptable {
					processRefund: "Process Refund"
					updateInventory: "Update Inventory"
					notifyRefund: "Notify Refund Complete"
					success: "Return Processed"
				} else {
					rejectItems: "Reject Items"
					notifyRejection: "Notify Customer of Rejection"
				}
			} else {
				timeout: "Return Timeout"
			}
		} else {
			notEligible: "Return Not Eligible"
		}
		end: "Complete"
	}
}

@service
component OrderEventHandler {
	private orderService: OrderService
	private notificationService: NotificationService
	private analyticsService: AnalyticsService

	method HandleOrderCreated(event: OrderCreatedEvent): void
	method HandleOrderShipped(event: OrderShippedEvent): void
	method HandleOrderDelivered(event: OrderDeliveredEvent): void
	method HandlePaymentFailed(event: PaymentFailedEvent): void
}
