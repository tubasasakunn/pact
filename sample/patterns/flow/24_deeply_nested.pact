// Pattern 24: Flow with 5+ nested conditions
component DeeplyNestedService {
    flow DeeplyNested {
        level1 = self.checkLevel1()
        if level1Valid {
            level2 = self.checkLevel2()
            if level2Valid {
                level3 = self.checkLevel3()
                if level3Valid {
                    level4 = self.checkLevel4()
                    if level4Valid {
                        level5 = self.checkLevel5()
                        if level5Valid {
                            result = self.executeDeepOperation()
                            return result
                        } else {
                            throw Level5Error
                        }
                    } else {
                        throw Level4Error
                    }
                } else {
                    throw Level3Error
                }
            } else {
                throw Level2Error
            }
        } else {
            throw Level1Error
        }
    }
}
