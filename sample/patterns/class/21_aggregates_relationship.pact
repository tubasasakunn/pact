// Pattern 21: Component with aggregates relationship
component Library {
	type LibraryData {
		name: string
		location: string
	}
	
	aggregates Book
	aggregates Member
}

component Book {
	type BookData {
		isbn: string
		title: string
		author: string
	}
}

component Member {
	type MemberData {
		memberId: string
		name: string
	}
}
