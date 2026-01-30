// Pattern 30: Real-world checkout flow
component CheckoutService {
    depends on CartService
    depends on InventoryService
    depends on PricingService
    depends on PaymentService
    depends on OrderService
    depends on NotificationService
    depends on FraudService

    flow RealCheckout {
        // Step 1: Validate cart
        cart = CartService.getCart(userId)
        if cartEmpty {
            throw EmptyCartError
        }

        cartValidation = CartService.validate(cart)
        if cartInvalid {
            errors = CartService.getValidationErrors(cartValidation)
            throw CartValidationError
        }

        // Step 2: Check inventory for all items
        for item in cart.items {
            availability = InventoryService.checkStock(item.productId, item.quantity)
            if outOfStock {
                NotificationService.sendOutOfStockAlert(userId, item)
                throw InsufficientStockError
            }
        }

        // Step 3: Calculate pricing
        subtotal = PricingService.calculateSubtotal(cart)
        discounts = PricingService.applyDiscounts(cart, userId)
        taxes = PricingService.calculateTaxes(subtotal, shippingAddress)
        shippingCost = PricingService.calculateShipping(cart, shippingAddress)
        total = PricingService.calculateTotal(subtotal, discounts, taxes, shippingCost)

        // Step 4: Fraud check
        fraudScore = FraudService.analyze(userId, total, paymentMethod)
        if fraudDetected {
            FraudService.flagTransaction(userId, fraudScore)
            NotificationService.sendFraudAlert(userId)
            throw FraudDetectedError
        }

        // Step 5: Reserve inventory
        for item in cart.items {
            reservation = InventoryService.reserve(item.productId, item.quantity)
            if reservationFailed {
                InventoryService.releaseAllReservations(cart)
                throw InventoryReservationError
            }
        }

        // Step 6: Process payment
        paymentAuth = PaymentService.authorize(paymentMethod, total)
        if authorizationFailed {
            InventoryService.releaseAllReservations(cart)
            throw PaymentAuthorizationError
        }

        paymentCapture = PaymentService.capture(paymentAuth.authorizationId)
        if captureFailed {
            PaymentService.voidAuthorization(paymentAuth.authorizationId)
            InventoryService.releaseAllReservations(cart)
            throw PaymentCaptureError
        }

        // Step 7: Create order
        order = OrderService.create(cart, paymentCapture, shippingAddress)
        if orderCreationFailed {
            PaymentService.refund(paymentCapture.transactionId)
            InventoryService.releaseAllReservations(cart)
            throw OrderCreationError
        }

        // Step 8: Commit inventory
        for item in cart.items {
            InventoryService.commitReservation(item.reservationId)
        }

        // Step 9: Clear cart and send notifications
        CartService.clear(userId)
        NotificationService.sendOrderConfirmation(userId, order)
        NotificationService.sendReceipt(userId, order, paymentCapture)

        return order
    }
}
