// Pattern 23: Real-world scenario - Order processing
component OrderController {
    depends on CartService
    depends on InventoryService
    depends on PaymentGateway
    depends on ShippingService
    depends on EmailService

    flow PlaceOrder {
        cart = CartService.getCart(userId)
        valid = self.validateCart(cart)
        if valid {
            available = InventoryService.checkAvailability(cart.items)
            if available {
                InventoryService.reserveItems(cart.items)
                payment = PaymentGateway.processPayment(cart.total)
                if payment.success {
                    order = self.createOrder(cart, payment)
                    shipping = ShippingService.createShipment(order)
                    EmailService.sendConfirmation(order)
                    CartService.clearCart(userId)
                    return order
                } else {
                    InventoryService.releaseItems(cart.items)
                    throw PaymentFailedError
                }
            } else {
                throw OutOfStockError
            }
        } else {
            throw InvalidCartError
        }
    }
}

component CartService { }
component InventoryService { }
component PaymentGateway { }
component ShippingService { }
component EmailService { }
