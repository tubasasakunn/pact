// Pattern 31: Type using all visibility modifiers (+, -, #, ~)
component VisibilityDemo {
	type PublicEntity {
		+id: string
		+name: string
		+description: string
	}

	type PrivateEntity {
		-internalId: string
		-secretKey: string
		-hashedPassword: string
	}

	type ProtectedEntity {
		#baseId: string
		#inheritedField: string
		#overridable: bool
	}

	type PackageEntity {
		~packageId: string
		~internalStatus: string
		~moduleConfig: string
	}

	type MixedVisibility {
		+publicId: string
		+publicName: string
		-privateSecret: string
		-privateKey: string
		#protectedBase: string
		#protectedConfig: string
		~packageInternal: string
		~packageModule: string
	}

	provides VisibilityAPI {
		GetPublic() -> PublicEntity
		GetMixed() -> MixedVisibility
	}
}
