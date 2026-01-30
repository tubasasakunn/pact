// Pattern 18: Component with extends relationship
component BaseEntity {
	type EntityData {
		id: string
		createdAt: string
		updatedAt: string
	}
}

component UserEntity {
	type UserData {
		name: string
		email: string
	}
	
	extends BaseEntity
}
