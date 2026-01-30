// Pattern 29: Type with 15+ fields (comprehensive entity)
component CustomerProfile {
	type DetailedCustomer {
		id: string
		firstName: string
		lastName: string
		middleName: string?
		email: string
		phone: string
		mobilePhone: string?
		dateOfBirth: string
		gender: string?
		addressLine1: string
		addressLine2: string?
		city: string
		stateCode: string
		postalCode: string
		country: string
		preferredLanguage: string
		timezone: string
		createdAt: string
		updatedAt: string
		lastLoginAt: string?
		isVerified: bool
		isActive: bool
		loyaltyPoints: int
		membershipTier: string
	}

	provides CustomerAPI {
		GetCustomer(id: string) -> DetailedCustomer
		CreateCustomer(customer: DetailedCustomer) -> DetailedCustomer
		UpdateCustomer(id: string, customer: DetailedCustomer) -> DetailedCustomer
		DeleteCustomer(id: string) -> bool
	}
}
