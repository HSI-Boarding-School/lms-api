package utils

func GetPermissionsByRole(role string) []string {
	switch role {
	case "ADMIN":
		return []string{
			"management user (crud, activate user, deactivate user, set role)",
			"management course (crud)",
			"log book management (crud)",
			"feedback management (crud)",
			"quiz management (crud)",
		}

	case "TEACHER":
		return []string{
			"management course (crud)",
			"log book management (crud)",
			"feedback management (crud)",
			"quiz management (crud)",
		}

	case "STUDENT":
		return []string{
			"course (read)",
			"answer (create)",
			"feedback (create, read)",
			"log book (read)",
		}

	default:
		return []string{}
	}
}
