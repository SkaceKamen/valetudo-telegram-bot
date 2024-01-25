package bot

func localizeAttachmentType(attachmentType string) string {
	switch attachmentType {
	case "mop":
		return "Mop"
	case "dustbin":
		return "Dustbin"
	case "watertank":
		return "Watertank"
	default:
		return attachmentType
	}
}

func localizeRobotStatus(robotStatus string) string {
	switch robotStatus {
	case "docked":
		return "Docked"
	case "idle":
		return "Idle"
	case "cleaning":
		return "Cleaning"
	case "paused":
		return "Paused"
	case "returning":
		return "Returning home"
	case "error":
		return "Error"
	case "manual_control":
		return "Manual control"
	}

	return robotStatus
}

func robotStatusEmoji(robotStatus string) string {
	switch robotStatus {
	case "docked":
		return "🏠"
	case "idle":
		return "💤"
	case "cleaning":
		return "🧹"
	case "paused":
		return "⏸"
	case "returning":
		return "🔙"
	case "error":
		return "❗"
	case "manual_control":
		return "🕹"
	}

	return "🤖"
}
