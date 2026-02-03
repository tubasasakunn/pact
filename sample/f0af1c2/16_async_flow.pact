// 16: Async Flow with await, error handling
component NotificationDispatcher {
    type Notification {
        id: string
        userId: string
        channel: string
        message: string
        priority: string
    }

    type DeliveryResult {
        success: bool
        channel: string
        timestamp: string
    }

    depends on EmailProvider
    depends on SMSProvider
    depends on PushProvider
    depends on UserPreferences
    depends on DeliveryLog

    flow DispatchNotification {
        prefs = UserPreferences.getPreferences(userId)
        results = self.initResults()

        if prefs.emailEnabled {
            await EmailProvider.send(userId, notification.message)
            self.recordResult(results, "email", true)
        }

        if prefs.smsEnabled {
            await SMSProvider.send(userId, notification.message)
            self.recordResult(results, "sms", true)
        }

        if prefs.pushEnabled {
            await PushProvider.send(userId, notification.message)
            self.recordResult(results, "push", true)
        }

        DeliveryLog.record(notification.id, results)
        return results
    }

    flow BroadcastMessage {
        users = UserPreferences.getAllActiveUsers()
        sent = 0

        for user in users {
            await self.DispatchNotification(user.id, message)
            sent = sent + 1
        }

        return sent
    }
}

component EmailProvider {
    provides EmailAPI {
        async Send(userId: string, message: string) -> bool
    }
}

component SMSProvider {
    provides SMSAPI {
        async Send(userId: string, message: string) -> bool
    }
}

component PushProvider {
    provides PushAPI {
        async Send(userId: string, message: string) -> bool
    }
}

component UserPreferences {
    provides PrefsAPI {
        GetPreferences(userId: string) -> string
        GetAllActiveUsers() -> string[]
    }
}

component DeliveryLog {
    provides LogAPI {
        Record(notificationId: string, results: string)
    }
}
