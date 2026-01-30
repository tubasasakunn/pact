// Pattern 25a: Shared definitions
component SharedTypes {
	type EntityBase {
		id: string
		createdAt: string
		updatedAt: string
	}
	
	type AuditInfo {
		createdBy: string
		updatedBy: string
	}
	
	enum Status {
		active
		inactive
		deleted
	}
}
