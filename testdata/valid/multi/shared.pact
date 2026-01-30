// Shared definitions
type ID = string
type Timestamp = int64

type AuditInfo {
	createdAt: Timestamp
	createdBy: ID
	updatedAt: Timestamp
	updatedBy: ID
}

interface Entity {
	method GetID(): ID
}

interface Auditable {
	method GetAuditInfo(): AuditInfo
}

component BaseEntity {
	implements Entity

	id: ID
	audit: AuditInfo

	method GetID(): ID
}
