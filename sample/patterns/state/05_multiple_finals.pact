// Pattern 5: Multiple final states
component MultipleFinals {
    states OrderOutcome {
        initial Pending
        final Completed
        final Cancelled
        final Refunded

        state Pending { }
        state Processing { }
        state Completed { }
        state Cancelled { }
        state Refunded { }

        Pending -> Processing on process
        Pending -> Cancelled on cancel
        Processing -> Completed on complete
        Processing -> Refunded on refund
    }
}
