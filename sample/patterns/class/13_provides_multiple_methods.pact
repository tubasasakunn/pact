// Pattern 13: Component with provides interface (multiple methods)
component ProductService {
	type Product {
		id: string
		name: string
		price: float
	}
	
	provides ProductAPI {
		GetProduct(id: string) -> Product
		CreateProduct(product: Product) -> Product
		UpdateProduct(id: string, product: Product) -> Product
		DeleteProduct(id: string)
		ListProducts() -> Product[]
	}
}
