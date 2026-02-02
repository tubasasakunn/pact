// Pattern 21: Complex guard conditions (simulated via event names)
component ComplexGuards {
    states TransactionProcessor {
        initial Validating

        state Validating { }
        state Processing { }
        state Approved { }
        state Rejected { }
        state Pending { }

        Validating -> Processing on validationPassed
        Validating -> Rejected on validationFailed
        Processing -> Approved on amountWithinLimit
        Processing -> Pending on amountExceedsLimit
        Pending -> Approved on manualApproval
        Pending -> Rejected on manualRejection
    }
}
