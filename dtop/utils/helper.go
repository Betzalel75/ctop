package utils

import (
	"fmt"
	"os/exec"
	"os/user"
)

// Capitalize the first letter of a string
func Capitalize(s string) string {
	n := []rune(s)
	first := true

	for i := range n {
		if isValid(n[i]) && first {
			if n[i] >= 'a' && n[i] <= 'z' {
				n[i] = n[i] - 32
			}
			first = false
		} else if n[i] >= 'A' && n[i] <= 'Z' {
			n[i] = n[i] + 32
		} else if !isValid(n[i]) {
			first = true
		}
	}

	return string(n)
}

func isValid(f rune) bool {
	if (f >= 'a' && f <= 'z') || (f >= 'A' && f <= 'Z') || (f >= '0' && f <= '9') {
		return true
	}
	return false
}

func CheckDockerPermissions() (bool, string) {
	// Step 1: Check if the 'docker' command exists in the user's PATH
	_, err := exec.LookPath("docker")
	if err != nil {
		return false, "The 'docker' command was not found. Please ensure Docker is installed."
	}

	// Step 2: Check if the user is in the 'docker' group
	currentUser, err := user.Current()
	if err != nil {
		return false, fmt.Sprintf("Error getting the current user: %v", err)
	}

	groups, err := currentUser.GroupIds()
	if err != nil {
		return false, fmt.Sprintf("Error getting user's groups: %v", err)
	}

	for _, groupID := range groups {
		group, err := user.LookupGroupId(groupID)
		if err != nil {
			// We can safely ignore lookup errors, as we may not have permissions
			// to look up all groups.
			continue
		}
		if group.Name == "docker" {
			return true, "The user has permissions to run Docker."
		}
	}

	return false, "The user is not in the 'docker' group. To add them, run 'sudo usermod -aG docker <your_username>' and then log out and log back in."
}
