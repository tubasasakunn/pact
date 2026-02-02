// Sequence Pattern 02: Multiple components
component OrderService {
    depends on UserService
    depends on PaymentService
    depends on InventoryService

    flow ProcessOrder {
        user = UserService.getUser(userId)
        available = InventoryService.checkStock(items)
        if available {
            payment = PaymentService.charge(user, amount)
            if payment.success {
                InventoryService.reserve(items)
                return self.createOrder(user, items)
            } else {
                throw PaymentError
            }
        } else {
            throw OutOfStockError
        }
    }
}

component UserService { }
component PaymentService { }
component InventoryService { }
