// Support Ticket - State Diagram Example
component SupportTicket {
    type Ticket {
        id: string
        title: string
        description: string
        priority: string
        assignee: string?
    }

    states TicketLifecycle {
        initial Open
        final Closed
        final Resolved

        state Open {
            entry [notify_support_team]
        }

        state Assigned {
            entry [notify_assignee]
        }

        state InProgress {
            entry [start_timer]
            exit [stop_timer]
        }

        state OnHold {
            entry [pause_sla]
            exit [resume_sla]
        }

        state Resolved {
            entry [send_resolution_email]
        }

        state Closed {
            entry [archive_ticket]
        }

        Open -> Assigned on assign
        Assigned -> InProgress on start_work
        InProgress -> OnHold on hold
        OnHold -> InProgress on resume
        InProgress -> Resolved on resolve
        Resolved -> Closed on close
        Resolved -> InProgress on reopen
    }

    flow AssignTicket {
        ticket = self.getTicket(ticketId)
        if ticket.status == "Open" {
            agent = self.findAvailableAgent(ticket.priority)
            if agent != null {
                self.assignToAgent(ticket, agent)
                NotificationService.notifyAgent(agent, ticket)
                return ticket
            } else {
                self.escalate(ticket)
                throw NoAgentAvailableError
            }
        } else {
            throw InvalidStateError
        }
    }
}
