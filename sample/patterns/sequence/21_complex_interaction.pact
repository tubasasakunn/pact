// Pattern 21: Complex interaction (5+ messages)
component WorkflowEngine {
    depends on TaskManager
    depends on NotificationService
    depends on AuditLog
    depends on StateManager

    flow ExecuteWorkflow {
        task = TaskManager.createTask(workflow)
        StateManager.initialize(task)
        AuditLog.logStart(task)
        NotificationService.notifyStakeholders(task)
        result = TaskManager.execute(task)
        StateManager.update(task, result)
        AuditLog.logComplete(task)
        NotificationService.notifyCompletion(task)
        TaskManager.archive(task)
        return result
    }
}

component TaskManager { }
component NotificationService { }
component AuditLog { }
component StateManager { }
