// State Pattern 02: Multiple states
component Order {
    states OrderStatus {
        initial Pending
        final Delivered
        final Cancelled

        state Pending { }
        state Confirmed { }
        state Processing { }
        state Shipped { }
        state Delivered { }
        state Cancelled { }

        Pending -> Confirmed on confirm
        Pending -> Cancelled on cancel
        Confirmed -> Processing on process
        Processing -> Shipped on ship
        Shipped -> Delivered on deliver
    }
}
