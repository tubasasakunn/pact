// EXPECTED BEHAVIOR: Self-import (circular reference)
// NOTE: Parser currently ALLOWS self-imports

import "10_circular_import.pact"

component CircularComponent {
    provides CircularAPI {
        Process() -> Void
    }
}
