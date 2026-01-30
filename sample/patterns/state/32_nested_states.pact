// Pattern 32: Nested/hierarchical states
component DocumentEditor {
    states EditorStates {
        initial Idle
        final Closed

        state Idle { }

        state Editing {
            entry [lockDocument, startAutoSave]
            exit [unlockDocument, stopAutoSave]

            initial Viewing

            state Viewing { }
            state Modifying {
                entry [enableUndo]
                exit [saveCheckpoint]
            }
            state Selecting {
                entry [showSelectionHandles]
                exit [hideSelectionHandles]
            }

            Viewing -> Modifying on edit
            Modifying -> Viewing on save
            Viewing -> Selecting on select
            Selecting -> Modifying on editSelection
            Selecting -> Viewing on deselect
        }

        state Saving {
            entry [showSaveDialog]
            exit [hideSaveDialog]

            initial Validating

            state Validating { }
            state Writing { }
            state Confirming { }

            Validating -> Writing on valid
            Validating -> Confirming on invalid
            Writing -> Confirming on done
        }

        state Closed { }

        Idle -> Editing on open
        Editing -> Saving on saveRequest
        Saving -> Editing on complete
        Saving -> Editing on cancel
        Editing -> Idle on close
        Idle -> Closed on exitApp
    }
}
