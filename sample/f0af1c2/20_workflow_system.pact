// 20: Complete Workflow System - flows + states + all features
@version("2.5.0")
@author("workflow-team")
@description("Approval workflow system")
component ApprovalWorkflow {
    @entity
    type Request {
        +id: string
        +requester: string
        +title: string
        +amount: float
        +status: string
    }

    enum ApprovalLevel {
        TeamLead
        Manager
        Director
    }

    depends on NotificationHub
    depends on AuditTrail

    provides WorkflowAPI {
        SubmitRequest(req: Request) -> Request throws ValidationError
        ApproveRequest(reqId: string, approverId: string) -> Request
        RejectRequest(reqId: string, reason: string) -> Request
    }

    flow SubmitForApproval {
        validated = self.validateRequest(request)
        if !validated {
            throw ValidationError
        }

        approvers = self.determineApproverChain(request.amount)

        for approver in approvers {
            NotificationHub.sendApprovalRequest(approver, request)
        }

        AuditTrail.log(request.requester, "submit", request.id)
        return request
    }

    flow ProcessApproval {
        request = self.getRequest(requestId)

        if authorized {
            self.recordApproval(request, approverId)
            NotificationHub.sendApproved(request.requester, request)
            AuditTrail.log(approverId, "approve", request.id)
            return request
        } else {
            throw UnauthorizedError
        }
    }

    states RequestLifecycle {
        initial Draft
        final Approved
        final Rejected

        state Draft {
            entry [initializeRequest]
        }

        state Submitted {
            entry [notifyFirstApprover]
        }

        state UnderReview {
            entry [setReviewDeadline]
        }

        state Approved {
            entry [executeApproval]
        }

        state Rejected {
            entry [notifyRequester]
        }

        Draft -> Submitted on submit
        Submitted -> UnderReview on startReview
        UnderReview -> Approved on finalApproval when allApproved
        UnderReview -> Rejected on reject
        UnderReview -> Submitted on requestChanges
    }
}

component NotificationHub {
    provides NotifyAPI {
        SendApprovalRequest(approverId: string, request: string) -> bool
        SendApproved(requesterId: string, request: string) -> bool
    }
}

component AuditTrail {
    type AuditLog {
        actor: string
        action: string
        target: string
    }

    provides AuditAPI {
        Log(actor: string, action: string, target: string)
        GetHistory(targetId: string) -> AuditLog[]
    }
}
