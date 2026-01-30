// Real-world Pattern 1: E-commerce System
// Demonstrates: Types, Components, States, Dependencies, Flows

// ============ TYPE DEFINITIONS ============

component User {
	type Address {
		street: string
		city: string
		province: string
		zipCode: string
		country: string
	}

	type UserProfile {
		id: string
		email: string
		name: string
		shippingAddress: Address?
		billingAddress: Address?
	}

	provides UserService {
		GetUser(userId: string) -> UserProfile
		UpdateProfile(userId: string, profile: UserProfile)
		GetAddresses(userId: string) -> Address[]
	}
}

component Product {
	type ProductDetails {
		id: string
		name: string
		description: string
		price: float
		stock: int
		category: string
	}

	provides ProductCatalog {
		GetProduct(productId: string) -> ProductDetails
		SearchProducts(query: string) -> ProductDetails[]
		CheckStock(productId: string) -> int
		UpdateStock(productId: string, quantity: int)
	}
}

component Cart {
	type CartItem {
		productId: string
		productName: string
		quantity: int
		unitPrice: float
		subtotal: float
	}

	type ShoppingCart {
		userId: string
		items: CartItem[]
		total: float
		itemCount: int
	}

	depends on Product

	provides CartService {
		GetCart(userId: string) -> ShoppingCart
		AddItem(userId: string, productId: string, quantity: int)
		RemoveItem(userId: string, productId: string)
		UpdateQuantity(userId: string, productId: string, quantity: int)
		ClearCart(userId: string)
	}

	flow AddToCart {
		stock = Product.CheckStock(productId)
		if stockAvailable {
			cart = self.GetCart(userId)
			self.AddItem(userId, productId, quantity)
			return cart
		} else {
			throw OutOfStockError
		}
	}
}

component Order {
	type OrderItem {
		productId: string
		productName: string
		quantity: int
		unitPrice: float
		subtotal: float
	}

	type OrderRecord {
		id: string
		userId: string
		items: OrderItem[]
		subtotal: float
		tax: float
		shipping: float
		total: float
		status: string
		createdAt: string
		updatedAt: string
	}

	depends on Cart
	depends on Product
	depends on Payment

	provides OrderService {
		CreateOrder(userId: string, paymentMethod: string) -> OrderRecord
		GetOrder(orderId: string) -> OrderRecord
		GetUserOrders(userId: string) -> OrderRecord[]
		UpdateStatus(orderId: string, status: string)
		CancelOrder(orderId: string)
	}

	// Order lifecycle states
	states OrderLifecycle {
		initial Pending

		state Pending {
			entry [notifyOrderCreated]
		}
		state Confirmed {
			entry [reserveInventory]
		}
		state Processing {
			entry [beginFulfillment]
		}
		state Shipped {
			entry [sendShipmentNotification]
		}
		state Delivered {
			entry [requestReview]
		}
		state Cancelled {
			entry [releaseInventory]
			exit [processRefund]
		}

		Pending -> Confirmed on paymentConfirmed
		Pending -> Cancelled on paymentFailed
		Pending -> Cancelled on userCancelled
		Confirmed -> Processing on processingStarted
		Confirmed -> Cancelled on userCancelled
		Processing -> Shipped on itemsShipped
		Shipped -> Delivered on deliveryConfirmed
		Shipped -> Cancelled on returnRequested
	}
}

component Payment {
	type PaymentDetails {
		orderId: string
		amount: float
		method: string
		cardLast4: string?
		transactionId: string?
		status: string
	}

	depends on User

	provides PaymentGateway {
		ProcessPayment(orderId: string, amount: float, method: string) -> PaymentDetails
		RefundPayment(transactionId: string, amount: float) -> bool
		GetPaymentStatus(transactionId: string) -> string
	}

	requires ExternalPaymentProvider {
		Authorize(amount: float, card: string) -> string
		Capture(authorizationId: string) -> bool
		Void(authorizationId: string) -> bool
	}
}

// ============ CHECKOUT FLOW ============

component Checkout {
	depends on User
	depends on Cart
	depends on Product
	depends on Order
	depends on Payment

	flow ProcessCheckout {
		user = User.GetUser(userId)
		cart = Cart.GetCart(userId)
		if cartNotEmpty {
			validated = self.validateCartItems(cart)
			if validated {
				for item in cart.items {
					stock = Product.CheckStock(item.productId)
					if insufficientStock {
						throw InsufficientStockError
					}
				}
				order = Order.CreateOrder(userId, paymentMethod)
				payment = Payment.ProcessPayment(order.id, cart.total, paymentMethod)
				if paymentSuccessful {
					Order.UpdateStatus(order.id, confirmed)
					for item in cart.items {
						Product.UpdateStock(item.productId, newQuantity)
					}
					Cart.ClearCart(userId)
					self.sendConfirmationEmail(user, order)
					return order
				} else {
					Order.UpdateStatus(order.id, cancelled)
					throw PaymentFailedError
				}
			} else {
				throw InvalidCartError
			}
		} else {
			throw EmptyCartError
		}
	}
}
