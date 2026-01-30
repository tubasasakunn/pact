// Pattern 30: Full-featured state machine
// States with entry, exit, all trigger types (on, when, after), guards, and actions
component FullFeatured {
    states OrderFulfillment {
        initial Received
        final Delivered
        final Cancelled

        state Received {
            entry [logOrderReceived, notifyWarehouse]
            exit [clearReceiveBuffer]
        }
        state Validating {
            entry [checkInventory, verifyAddress]
        }
        state PaymentPending {
            entry [initiatePayment, setPaymentTimer]
            exit [clearPaymentTimer]
        }
        state PaymentProcessing {
            entry [processPayment]
        }
        state Confirmed {
            entry [sendConfirmationEmail, reserveInventory]
            exit [finalizeReservation]
        }
        state Picking {
            entry [assignPicker, printPickList]
            exit [confirmItemsPicked]
        }
        state Packing {
            entry [assignPacker, generateShippingLabel]
            exit [sealPackage]
        }
        state ReadyToShip {
            entry [schedulePickup, notifyCarrier]
        }
        state Shipped {
            entry [updateTrackingInfo, notifyCustomer]
        }
        state OutForDelivery {
            entry [sendDeliveryAlert]
        }
        state Delivered {
            entry [confirmDelivery, requestFeedback]
        }
        state Cancelled {
            entry [releaseInventory, processRefund, notifyCancellation]
        }

        // Event triggers (on)
        Received -> Validating on processOrder
        Validating -> PaymentPending on validationPassed
        PaymentPending -> PaymentProcessing on paymentReceived
        PaymentProcessing -> Confirmed on paymentApproved
        Confirmed -> Picking on startFulfillment
        Picking -> Packing on itemsPicked
        Packing -> ReadyToShip on packingComplete
        ReadyToShip -> Shipped on carrierPickup
        Shipped -> OutForDelivery on outForDelivery
        OutForDelivery -> Delivered on deliveryConfirmed

        // Event triggers with guards (on ... when)
        Validating -> Cancelled on validationFailed when invalidAddress
        Validating -> Cancelled on validationFailed when outOfStock
        PaymentProcessing -> Cancelled on paymentFailed when cardDeclined
        PaymentProcessing -> PaymentPending on paymentFailed when retryAvailable

        // Condition triggers (when)
        Received -> Validating when autoProcessEnabled
        Confirmed -> Picking when warehouseReady
        ReadyToShip -> Shipped when carrierAvailable

        // Time triggers (after)
        PaymentPending -> Cancelled after 24h
        ReadyToShip -> Cancelled after 7d
        Shipped -> Delivered after 14d

        // Self transitions
        Picking -> Picking on itemNotFound
        Packing -> Packing on qualityCheckFailed
        OutForDelivery -> OutForDelivery on deliveryAttemptFailed

        // Cross-cutting cancellation
        Received -> Cancelled on customerCancelled
        Validating -> Cancelled on customerCancelled
        PaymentPending -> Cancelled on customerCancelled
        Confirmed -> Cancelled on customerCancelled
        Picking -> Cancelled on customerCancelled
    }
}
