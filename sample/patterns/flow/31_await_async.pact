// Pattern 31: Await and async operations
component AsyncService {
    depends on ExternalAPI : APIClient as api
    depends on Database : DB as db

    provides AsyncOperations {
        async fetchData(id: string) -> Data
        async saveData(data: Data) -> bool
        async processAsync(items: Data[]) -> Result[]
    }

    flow FetchAndSave {
        await api.fetchData(resourceId)
        data = api.getResult()
        validated = self.validate(data)
        if validated {
            await db.save(data)
            result = db.getResult()
            return result
        } else {
            throw ValidationError
        }
    }

    flow ParallelFetch {
        await api.fetchUser(userId)
        result1 = api.getResult()
        await api.fetchOrders(userId)
        result2 = api.getResult()
        combined = self.combine(result1, result2)
        await db.saveUserData(combined)
        return combined
    }

    flow ChainedAsync {
        await self.firstStep()
        step1 = self.getResult()
        await self.secondStep(step1)
        step2 = self.getResult()
        await self.thirdStep(step2)
        step3 = self.getResult()
        return step3
    }
}
