package main

import (
	"fmt"
	"net/http"
)

var events = []string{"wildcard", "check_run", "check_suite", "commit_comment", "create", "delete", "deployment", "deployment_status", "fork", "github_app_authorization", "gollum", "installation", "installation_repositories", "issue_comment", "issues", "label", "marketplace_purchase", "member", "membership", "milestone", "organization", "org_block", "page_build", "project_card", "project_column", "project", "public", "pull_request_review_comment", "pull_request_review", "pull_request", "push", "repository", "repository_import", "repository_vulnerability_alert", "release", "security_advisory", "status", "team", "team_add", "watch"}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	http.ListenAndServe(":80", nil)
}

type CreateEvent struct {
	Ref          string      `json:"ref"`
	RefType      string      `json:"ref_type"`
	MasterBranch string      `json:"master_branch"`
	Description  interface{} `json:"description"`
	PusherType   string      `json:"pusher_type"`
}
