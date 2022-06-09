package utils

import "strings"

// GetProjectFromIssueKey - ...
func GetProjectFromIssueKey(issuekey string) string {
	if data := strings.Split(issuekey, "-"); len(data) > 0 {
		return data[0]
	}
	return ""
}
