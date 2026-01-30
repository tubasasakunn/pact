// Pattern 31: Transition do actions
component OrderProcessor {
    states OrderWorkflow {
        initial Pending
        final Completed
        final Cancelled

        state Pending { }
        state Processing { }
        state Shipping { }
        state Delivered { }
        state Completed { }
        state Cancelled { }

        // Simple do action
        Pending -> Processing on confirm do [notifyWarehouse]

        // Multiple do actions
        Processing -> Shipping on ship do [updateInventory, createShipment, sendTrackingEmail]

        // Do action with guard
        Shipping -> Delivered on deliver when addressVerified do [confirmDelivery, updateStatus]

        // Do action to final state
        Delivered -> Completed on complete do [archiveOrder, sendSurvey]

        // Do action on cancellation
        Pending -> Cancelled on cancel do [releaseReservation, notifyCustomer, logCancellation]
        Processing -> Cancelled on cancel do [revertInventory, cancelShipment, refundPayment]
    }
}
