// Pattern 17: Component with multiple annotations
@version("2.1.0")
@author("dev-team")
@deprecated("Use NewService instead")
component AnnotatedService {
	type Config {
		setting: string
		enabled: bool
	}
	
	provides ConfigAPI {
		@cache(ttl: "300")
		GetConfig() -> Config
		
		@async
		@transactional
		UpdateConfig(config: Config)
	}
}
