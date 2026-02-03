// 07: Simple State Machine - initial, final, entry/exit, event triggers
component DocumentWorkflow {
    type Document {
        id: string
        title: string
        content: string
        author: string
    }

    states DocumentLifecycle {
        initial Draft
        final Published
        final Archived

        state Draft {
            entry [createDraft]
        }

        state Review {
            entry [notifyReviewers]
            exit [clearReviewQueue]
        }

        state Approved {
            entry [prepareForPublish]
        }

        state Published {
            entry [notifySubscribers, updateIndex]
        }

        state Archived {
            entry [moveToArchive]
        }

        Draft -> Review on submitForReview
        Review -> Draft on requestChanges
        Review -> Approved on approve
        Approved -> Published on publish
        Published -> Archived on archive
        Published -> Draft on unpublish
    }
}
