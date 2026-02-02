// Pattern 27: Transitions with complex guard conditions
// Demonstrates guard conditions on transitions
component ComplexGuardsDemo {
    states PaymentProcessor {
        initial Validating
        final Success
        final Rejected

        state Validating {
            entry [validateInput]
        }
        state CheckingBalance { }
        state CheckingFraud { }
        state Processing { }
        state PendingApproval { }
        state Success { }
        state Rejected { }

        // Transitions with guards (using when for guard conditions)
        Validating -> CheckingBalance on validated when isValidCard
        Validating -> Rejected on validated when isInvalidCard

        CheckingBalance -> CheckingFraud on balanceChecked when hasSufficientFunds
        CheckingBalance -> Rejected on balanceChecked when insufficientFunds

        CheckingFraud -> Processing on fraudChecked when lowRisk
        CheckingFraud -> PendingApproval on fraudChecked when mediumRisk
        CheckingFraud -> Rejected on fraudChecked when highRisk

        Processing -> Success on processed when amountBelowLimit
        Processing -> PendingApproval on processed when amountAboveLimit

        PendingApproval -> Success on approved when managerApproved
        PendingApproval -> Rejected on denied when managerDenied
    }
}
