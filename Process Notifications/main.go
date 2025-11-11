package main

type notification interface {
	importance() int
}

type directMessage struct {
	senderUsername string
	messageContent string
	priorityLevel  int
	isUrgent       bool
}

type groupMessage struct {
	groupName      string
	messageContent string
	priorityLevel  int
}

type systemAlert struct {
	alertCode      string
	messageContent string
}

// ?
func (s systemAlert) importance() int {
	return 100
}
func (g groupMessage) importance() int {
	return g.priorityLevel
}
func (d directMessage) importance() int {
	if d.isUrgent {
		return 50
	}
	return d.priorityLevel
}

func processNotification(n notification) (string, int) {
	// ?
	switch c := n.(type) {
	case systemAlert:
		return c.alertCode, c.importance()
	case groupMessage:
		return c.groupName, c.importance()
	case directMessage:
		return c.senderUsername, c.importance()
	default:
		return "", 0
	}
}
