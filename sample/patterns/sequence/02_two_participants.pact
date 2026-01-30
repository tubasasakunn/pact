// Pattern 2: Two participants
component ClientService {
    depends on BackendService

    flow FetchData {
        result = BackendService.getData(id)
        return result
    }
}

component BackendService {
    provides DataAPI {
        GetData(id: string)
    }
}
