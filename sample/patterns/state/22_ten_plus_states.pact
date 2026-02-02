// Pattern 22: State machine with 10+ states
component TenPlusStates {
    states OrderFulfillment {
        initial New
        final Completed
        final Cancelled

        state New { }
        state Validated { }
        state PaymentPending { }
        state PaymentConfirmed { }
        state Picking { }
        state Packing { }
        state ReadyToShip { }
        state Shipped { }
        state InTransit { }
        state OutForDelivery { }
        state Delivered { }
        state Completed { }
        state Cancelled { }

        New -> Validated on validate
        Validated -> PaymentPending on requestPayment
        PaymentPending -> PaymentConfirmed on paymentReceived
        PaymentConfirmed -> Picking on startPicking
        Picking -> Packing on pickingComplete
        Packing -> ReadyToShip on packingComplete
        ReadyToShip -> Shipped on handToCarrier
        Shipped -> InTransit on carrierPickup
        InTransit -> OutForDelivery on outForDelivery
        OutForDelivery -> Delivered on delivered
        Delivered -> Completed on confirm
        New -> Cancelled on cancel
        Validated -> Cancelled on cancel
        PaymentPending -> Cancelled on cancel
    }
}
