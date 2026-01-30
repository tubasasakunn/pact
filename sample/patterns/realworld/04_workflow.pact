// Real-world Pattern 4: Approval Workflow System
// Demonstrates: Multi-step approval, State machines, Complex business logic

// ============ TYPE DEFINITIONS ============

component Task {
	type TaskDetails {
		id: string
		title: string
		description: string
		creatorId: string
		assigneeId: string?
		priority: string
		category: string
		dueDate: string?
		attachments: string[]
		createdAt: string
		updatedAt: string
		status: string
	}

	type TaskComment {
		id: string
		taskId: string
		authorId: string
		content: string
		createdAt: string
	}

	type TaskHistory {
		id: string
		taskId: string
		action: string
		actorId: string
		previousValue: string?
		newValue: string?
		timestamp: string
	}

	depends on User
	depends on Approval
	depends on NotificationService

	provides TaskService {
		CreateTask(details: TaskDetails) -> TaskDetails
		GetTask(taskId: string) -> TaskDetails
		UpdateTask(taskId: string, updates: TaskDetails) -> TaskDetails
		DeleteTask(taskId: string) -> bool
		SubmitForApproval(taskId: string) -> bool
		GetTaskHistory(taskId: string) -> TaskHistory[]
		AddComment(taskId: string, comment: TaskComment) -> TaskComment
	}

	// Task lifecycle states
	states TaskLifecycle {
		initial Draft

		state Draft {
			entry [notifyCreator]
		}
		state Submitted {
			entry [assignReviewers]
		}
		state UnderReview {
			entry [notifyReviewers]
		}
		state PendingChanges {
			entry [notifyCreatorOfChanges]
		}
		state Approved {
			entry [notifyAllStakeholders]
		}
		state Rejected {
			entry [notifyCreatorOfRejection]
		}
		state Archived {
			entry [cleanupResources]
		}

		Draft -> Submitted on submitTask
		Draft -> Archived on deleteTask
		Submitted -> UnderReview on reviewStarted
		Submitted -> Draft on withdrawTask
		UnderReview -> PendingChanges on changesRequested
		UnderReview -> Approved on allApproved
		UnderReview -> Rejected on rejected
		PendingChanges -> Submitted on changesSubmitted
		PendingChanges -> Draft on withdrawTask
		Approved -> Archived on archiveTask
		Rejected -> Draft on reviseTask
		Rejected -> Archived on archiveTask
	}

	flow SubmitTaskForApproval {
		task = self.GetTask(taskId)
		if taskInDraft {
			validated = self.validateTask(task)
			if validated {
				approvers = Approval.DetermineApprovers(task)
				if noApproversFound {
					throw NoApproversConfiguredError
				}
				for approver in approvers {
					Approval.CreateApprovalRequest(task.id, approver.id)
				}
				task = self.updateStatus(task.id, submitted)
				self.recordHistory(task.id, submitted, creatorId)
				NotificationService.NotifySubmission(task, approvers)
				return task
			} else {
				throw TaskValidationError
			}
		} else {
			throw InvalidTaskStateError
		}
	}
}

component Approval {
	type ApprovalRequest {
		id: string
		taskId: string
		approverId: string
		status: string
		level: int
		decision: string?
		comments: string?
		decidedAt: string?
		createdAt: string
	}

	type ApprovalRule {
		id: string
		category: string
		minApprovers: int
		requiredLevels: int[]
		escalationTimeout: int
		autoApproveAmount: float?
	}

	type ApprovalChain {
		taskId: string
		requests: ApprovalRequest[]
		currentLevel: int
		totalLevels: int
		status: string
	}

	depends on User
	depends on NotificationService

	provides ApprovalService {
		CreateApprovalRequest(taskId: string, approverId: string) -> ApprovalRequest
		GetApprovalRequest(requestId: string) -> ApprovalRequest
		GetTaskApprovals(taskId: string) -> ApprovalRequest[]
		DetermineApprovers(task: TaskDetails) -> User[]
		ApproveTask(requestId: string, comments: string?) -> ApprovalRequest
		RejectTask(requestId: string, comments: string) -> ApprovalRequest
		RequestChanges(requestId: string, comments: string) -> ApprovalRequest
		GetApprovalRules(category: string) -> ApprovalRule
		EscalateApproval(requestId: string) -> ApprovalRequest
	}

	// Approval request states
	states ApprovalRequestLifecycle {
		initial Pending

		state Pending {
			entry [startTimer]
		}
		state InProgress {
			entry [lockForReview]
		}
		state Approved {
			entry [recordApproval]
		}
		state Rejected {
			entry [recordRejection]
		}
		state ChangesRequested {
			entry [recordChangeRequest]
		}
		state Escalated {
			entry [notifyEscalation]
		}
		state Expired {
			entry [handleExpiration]
		}

		Pending -> InProgress on reviewStarted
		Pending -> Escalated on timeoutReached
		Pending -> Expired on abandonedTimeout
		InProgress -> Approved on approvalGiven
		InProgress -> Rejected on rejectionGiven
		InProgress -> ChangesRequested on changesNeeded
		InProgress -> Escalated on manualEscalation
		Escalated -> InProgress on escalationAccepted
		Escalated -> Expired on escalationTimeout
	}

	flow ProcessApprovalDecision {
		request = self.GetApprovalRequest(requestId)
		if requestPending {
			approver = User.GetUser(approverId)
			authorized = self.checkApproverAuthorization(approver, request)
			if authorized {
				if decisionApprove {
					request = self.ApproveTask(requestId, comments)
					allApprovals = self.GetTaskApprovals(request.taskId)
					rules = self.GetApprovalRules(taskCategory)
					if allRequiredApprovalsReceived {
						Task.UpdateStatus(request.taskId, approved)
						NotificationService.NotifyApproval(request.taskId)
						return request
					} else {
						nextLevel = self.determineNextLevel(allApprovals, rules)
						nextApprovers = self.getApproversForLevel(nextLevel)
						for approver in nextApprovers {
							self.CreateApprovalRequest(request.taskId, approver.id)
						}
						NotificationService.NotifyNextLevelApprovers(nextApprovers)
						return request
					}
				} else {
					if decisionReject {
						request = self.RejectTask(requestId, comments)
						Task.UpdateStatus(request.taskId, rejected)
						NotificationService.NotifyRejection(request.taskId, comments)
						return request
					} else {
						request = self.RequestChanges(requestId, comments)
						Task.UpdateStatus(request.taskId, pendingChanges)
						NotificationService.NotifyChangesRequested(request.taskId, comments)
						return request
					}
				}
			} else {
				throw UnauthorizedApproverError
			}
		} else {
			throw InvalidApprovalStateError
		}
	}
}

component User {
	type UserProfile {
		id: string
		name: string
		email: string
		department: string
		role: string
		managerId: string?
		approvalLevel: int
		delegateTo: string?
		delegateUntil: string?
	}

	provides UserService {
		GetUser(userId: string) -> UserProfile
		GetUsersByDepartment(department: string) -> UserProfile[]
		GetUsersByRole(role: string) -> UserProfile[]
		GetApproversForLevel(level: int) -> UserProfile[]
		GetManager(userId: string) -> UserProfile?
		SetDelegate(userId: string, delegateId: string, until: string)
		GetEffectiveApprover(userId: string) -> UserProfile
	}
}

component NotificationService {
	type NotificationMessage {
		recipientId: string
		notificationType: string
		subject: string
		body: string
		taskId: string?
		actionUrl: string?
	}

	provides NotificationAPI {
		NotifySubmission(task: TaskDetails, approvers: UserProfile[])
		NotifyApproval(taskId: string)
		NotifyRejection(taskId: string, reason: string)
		NotifyChangesRequested(taskId: string, comments: string)
		NotifyNextLevelApprovers(approvers: UserProfile[])
		NotifyEscalation(taskId: string, escalatedTo: UserProfile)
		SendReminder(approverId: string, taskId: string)
	}

	requires EmailService {
		Send(to: string, subject: string, body: string) -> bool
	}

	requires SlackService {
		PostMessage(channel: string, message: string) -> bool
		SendDirectMessage(userId: string, message: string) -> bool
	}
}

// ============ WORKFLOW ORCHESTRATOR ============

component WorkflowOrchestrator {
	depends on Task
	depends on Approval
	depends on User
	depends on NotificationService

	provides WorkflowAPI {
		StartWorkflow(taskId: string) -> bool
		GetWorkflowStatus(taskId: string) -> string
		CancelWorkflow(taskId: string) -> bool
		RetryWorkflow(taskId: string) -> bool
	}

	flow ExecuteApprovalWorkflow {
		task = Task.GetTask(taskId)
		rules = Approval.GetApprovalRules(task.category)
		currentLevel = 1
		while levelNotComplete {
			approvers = User.GetApproversForLevel(currentLevel)
			for approver in approvers {
				effectiveApprover = User.GetEffectiveApprover(approver.id)
				Approval.CreateApprovalRequest(taskId, effectiveApprover.id)
				NotificationService.SendApprovalNotification(effectiveApprover, task)
			}
			self.waitForApprovals(taskId, currentLevel, rules.timeout)
			approvals = Approval.GetTaskApprovals(taskId)
			levelApprovals = self.filterByLevel(approvals, currentLevel)
			if levelApprovalsMet {
				currentLevel = currentLevel + 1
				if currentLevel > rules.requiredLevels {
					Task.UpdateStatus(taskId, approved)
					NotificationService.NotifyApproval(taskId)
					return success
				}
			} else {
				if anyRejection {
					Task.UpdateStatus(taskId, rejected)
					NotificationService.NotifyRejection(taskId, rejectionReason)
					return failure
				}
			}
		}
	}
}
