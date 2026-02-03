// 12: Sequence Diagram - multi-service interactions
component CheckoutService {
    type Cart {
        userId: string
        items: string[]
        total: float
    }

    type OrderConfirmation {
        orderId: string
        estimatedDelivery: string
    }

    depends on CartService
    depends on InventoryService
    depends on PaymentService
    depends on ShippingService
    depends on NotificationService
    depends on AnalyticsService

    flow Checkout {
        cart = CartService.getCart(userId)

        for item in cart.items {
            stock = InventoryService.checkAvailability(item)
            if !stock {
                throw OutOfStockError
            }
        }

        total = CartService.calculateTotal(cart)
        payment = PaymentService.processPayment(userId, total)

        if payment.success {
            InventoryService.reserveItems(cart.items)
            shipment = ShippingService.createShipment(cart, userId)
            order = self.createOrder(cart, payment, shipment)
            NotificationService.sendOrderConfirmation(userId, order)
            AnalyticsService.trackPurchase(userId, order)
            CartService.clearCart(userId)
            return order
        } else {
            NotificationService.sendPaymentFailed(userId)
            throw PaymentFailedError
        }
    }
}

component CartService {
    provides CartAPI {
        GetCart(userId: string) -> string
        CalculateTotal(cart: string) -> float
        ClearCart(userId: string) -> bool
    }
}

component InventoryService {
    provides InventoryAPI {
        CheckAvailability(item: string) -> bool
        ReserveItems(items: string[]) -> bool
    }
}

component PaymentService {
    provides PaymentAPI {
        ProcessPayment(userId: string, amount: float) -> string
    }
}

component ShippingService {
    provides ShippingAPI {
        CreateShipment(cart: string, userId: string) -> string
    }
}

component NotificationService {
    provides NotifyAPI {
        SendOrderConfirmation(userId: string, order: string)
        SendPaymentFailed(userId: string)
    }
}

component AnalyticsService {
    provides AnalyticsAPI {
        TrackPurchase(userId: string, order: string)
    }
}
