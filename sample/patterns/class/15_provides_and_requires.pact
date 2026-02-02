// Pattern 15: Component with both provides and requires
component ShoppingCart {
	type CartItem {
		productId: string
		quantity: int
		price: float
	}
	
	provides CartAPI {
		AddItem(item: CartItem)
		RemoveItem(productId: string)
		GetTotal() -> float
	}
	
	requires InventoryService {
		CheckStock(productId: string) -> int
		ReserveStock(productId: string, quantity: int) -> bool
	}
	
	requires PricingService {
		GetPrice(productId: string) -> float
		ApplyDiscount(total: float, code: string) -> float
	}
}
