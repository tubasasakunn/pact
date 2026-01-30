// Order Processing - Sequence/Flow Diagram Example
component OrderProcessor {
    type Order {
        id: string
        customerId: string
        items: string[]
        total: float
        status: string
    }

    type PaymentResult {
        success: bool
        transactionId: string
    }

    depends on PaymentGateway
    depends on InventoryService
    depends on NotificationService

    flow ProcessOrder {
        // Validate order
        valid = self.validateOrder(order)
        if valid {
            // Check inventory
            available = InventoryService.checkStock(order.items)
            if available {
                // Process payment
                payment = PaymentGateway.charge(order.total)
                if payment.success {
                    // Reserve inventory
                    InventoryService.reserve(order.items)
                    // Send confirmation
                    NotificationService.sendConfirmation(order)
                    return order
                } else {
                    throw PaymentFailedError
                }
            } else {
                throw OutOfStockError
            }
        } else {
            throw ValidationError
        }
    }

    flow CancelOrder {
        order = self.getOrder(orderId)
        if order.status == "pending" {
            InventoryService.release(order.items)
            PaymentGateway.refund(order.id)
            NotificationService.sendCancellation(order)
            return true
        } else {
            throw CannotCancelError
        }
    }
}

component PaymentGateway {
    type ChargeRequest {
        amount: float
        currency: string
    }

    provides PaymentAPI {
        Charge(amount: float) -> PaymentResult
        Refund(orderId: string) -> bool
    }
}

component InventoryService {
    type StockItem {
        productId: string
        quantity: int
    }

    provides InventoryAPI {
        CheckStock(items: string[]) -> bool
        Reserve(items: string[]) -> bool
        Release(items: string[]) -> bool
    }
}

component NotificationService {
    provides NotificationAPI {
        SendConfirmation(order: Order)
        SendCancellation(order: Order)
    }
}
