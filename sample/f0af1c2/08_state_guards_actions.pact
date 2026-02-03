// 08: State Machine with guards, actions, time triggers, self-transitions
component OrderFulfillment {
    type Order {
        id: string
        total: float
        status: string
    }

    states OrderProcess {
        initial Received
        final Delivered
        final Cancelled

        state Received {
            entry [logReceived]
        }

        state Validating {
            entry [startValidation]
            exit [stopValidation]
        }

        state PaymentPending {
            entry [requestPayment]
        }

        state Processing {
            entry [startProcessing]
        }

        state Shipping {
            entry [createShipment, notifyCustomer]
            exit [updateTracking]
        }

        state Delivered {
            entry [confirmDelivery, sendSurvey]
        }

        state Cancelled {
            entry [processRefund, notifyCustomer]
        }

        // Event triggers
        Received -> Validating on process

        // Guard conditions
        Validating -> PaymentPending on validated when isValid
        Validating -> Cancelled on validated when isInvalid

        // Time-based triggers
        PaymentPending -> Cancelled after 24h

        // Event with actions
        PaymentPending -> Processing on paymentReceived do [updateLedger, sendReceipt]

        Processing -> Shipping on readyToShip
        Shipping -> Delivered on delivered

        // Self transition
        Processing -> Processing on updateProgress

        // Any state cancellation
        Received -> Cancelled on cancel
        Validating -> Cancelled on cancel
        Processing -> Cancelled on cancel
    }
}
