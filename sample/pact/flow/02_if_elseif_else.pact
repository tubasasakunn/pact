// Pattern: If-ElseIf-Else Flow
// 複数条件分岐を含むフローパターン

component PricingService {
    type PriceResult {
        finalPrice: float
        discount: float
        tier: string
    }

    flow CalculatePrice {
        basePrice = self.getBasePrice(productId)
        customer = self.getCustomer(customerId)
        purchaseHistory = self.getPurchaseHistory(customerId)

        if purchaseHistory.isPlatinum {
            discount = self.applyPlatinumDiscount(basePrice)
        } else {
            if purchaseHistory.isGold {
                discount = self.applyGoldDiscount(basePrice)
            } else {
                if purchaseHistory.isSilver {
                    discount = self.applySilverDiscount(basePrice)
                } else {
                    discount = self.applyStandardDiscount(basePrice)
                }
            }
        }

        finalPrice = self.calculateFinalPrice(basePrice, discount)
        return finalPrice
    }
}
