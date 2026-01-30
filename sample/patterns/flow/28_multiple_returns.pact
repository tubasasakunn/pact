// Pattern 28: Flow with 5+ return statements in different branches
component MultipleReturnsService {
    depends on CacheService
    depends on DatabaseService
    depends on ExternalApi

    flow MultipleReturns {
        // Early return 1: Check cache first
        cached = CacheService.get(key)
        if cacheHit {
            self.logCacheHit(key)
            return cached
        }

        // Early return 2: Check local database
        localData = DatabaseService.findLocal(key)
        if localFound {
            CacheService.set(key, localData)
            self.logLocalHit(key)
            return localData
        }

        // Early return 3: Check replica database
        replicaData = DatabaseService.findReplica(key)
        if replicaFound {
            CacheService.set(key, replicaData)
            self.logReplicaHit(key)
            return replicaData
        }

        // Early return 4: Check external API with fallback
        apiAvailable = ExternalApi.healthCheck()
        if apiAvailable {
            apiData = ExternalApi.fetch(key)
            if apiFetchSuccess {
                DatabaseService.save(key, apiData)
                CacheService.set(key, apiData)
                self.logApiHit(key)
                return apiData
            }
        }

        // Early return 5: Return stale data if available
        staleData = CacheService.getStale(key)
        if staleAvailable {
            self.logStaleReturn(key)
            return staleData
        }

        // Final return 6: Return default value
        defaultValue = self.getDefaultValue(key)
        self.logDefaultReturn(key)
        return defaultValue
    }
}
