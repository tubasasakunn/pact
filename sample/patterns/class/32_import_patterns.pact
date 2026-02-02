// Pattern 32: Import patterns
import "common/types" as types
import "shared/utils"
import "external/api" as externalApi

component ImportDemo {
    depends on CommonService : Service as commonSvc
    depends on ExternalClient : Client as extClient

    type LocalType {
        id: string
        externalRef: string
        data: string
    }

    provides LocalAPI {
        processWithImported(input: InputData) -> OutputData
        useExternal(id: string) -> Result
    }

    flow UseImportedTypes {
        input = self.createInput(rawData)
        processed = commonSvc.process(input)
        result = extClient.send(processed)
        return result
    }
}
