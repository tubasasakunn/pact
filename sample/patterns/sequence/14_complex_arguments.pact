// Pattern 14: Call with complex arguments
component SearchService {
    depends on SearchEngine

    flow PerformSearch {
        results = SearchEngine.search(query, filters, pagination, sorting)
        facets = SearchEngine.getFacets(query, filters)
        suggestions = SearchEngine.getSuggestions(query, userContext)
        return results
    }
}

component SearchEngine { }
