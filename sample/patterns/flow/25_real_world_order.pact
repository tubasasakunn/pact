// Pattern 25: Complex real-world flow (order processing)
component OrderProcessingService {
    depends on InventoryService
    depends on PaymentGateway
    depends on ShippingService
    depends on NotificationService
    depends on FraudDetectionService
    depends on LoyaltyService

    flow ProcessOrder {
        order = self.validateOrder(orderRequest)
        if orderValid {
            fraudCheck = FraudDetectionService.analyze(order)
            if fraudCheckPassed {
                inventory = InventoryService.checkAvailability(order.items)
                if itemsAvailable {
                    reserved = InventoryService.reserve(order.items)
                    payment = PaymentGateway.authorize(order.payment)
                    if paymentAuthorized {
                        captured = PaymentGateway.capture(payment.authId)
                        if captureSuccess {
                            LoyaltyService.awardPoints(order.customerId, order.total)
                            shipment = ShippingService.createShipment(order)
                            trackingNumber = ShippingService.getTracking(shipment)
                            order = self.updateOrderStatus(order, trackingNumber)
                            NotificationService.sendConfirmation(order)
                            return order
                        } else {
                            InventoryService.release(reserved)
                            PaymentGateway.void(payment.authId)
                            throw PaymentCaptureError
                        }
                    } else {
                        InventoryService.release(reserved)
                        NotificationService.sendPaymentFailure(order)
                        throw PaymentAuthorizationError
                    }
                } else {
                    NotificationService.sendOutOfStock(order)
                    throw OutOfStockError
                }
            } else {
                NotificationService.sendFraudAlert(order)
                throw FraudDetectedError
            }
        } else {
            throw OrderValidationError
        }
    }
}
