// Shared definitions
component BaseEntity {
	type AuditInfo {
		createdAt: int
		createdBy: string
		updatedAt: int
		updatedBy: string
	}

	type EntityData {
		id: string
		audit: AuditInfo
	}

	provides EntityAPI {
		GetID() -> string
	}
}
