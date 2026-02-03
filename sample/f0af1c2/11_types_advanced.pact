// 11: Advanced Types - arrays, nullable, nested types
component DataModel {
    type Organization {
        +id: string
        +name: string
        +departments: Department[]
        +owner: User?
        +metadata: string[]
    }

    type Department {
        +id: string
        +name: string
        +manager: User?
        +members: User[]
        +budget: float
    }

    type User {
        +id: string
        +name: string
        +email: string
        -passwordHash: string
        #permissions: Permission[]
        ~lastLogin: string?
    }

    type Permission {
        resource: string
        action: string
        granted: bool
    }

    type AuditEntry {
        timestamp: string
        actor: User
        action: string
        target: string?
        details: string?
    }

    enum AccessLevel {
        None
        Read
        Write
        Admin
        SuperAdmin
    }

    enum Department_Status {
        Active
        Inactive
        Restructuring
        Dissolved
    }
}
