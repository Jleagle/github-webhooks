package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const signaturePrefix = "sha1="
const signatureLength = 45

var events = []string{"wildcard", "check_run", "check_suite", "commit_comment", "create", "delete", "deployment", "deployment_status", "fork", "github_app_authorization", "gollum", "installation", "installation_repositories", "issue_comment", "issues", "label", "marketplace_purchase", "member", "membership", "milestone", "organization", "org_block", "page_build", "project_card", "project_column", "project", "public", "pull_request_review_comment", "pull_request_review", "pull_request", "push", "repository", "repository_import", "repository_vulnerability_alert", "release", "security_advisory", "status", "team", "team_add", "watch"}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var secret = os.Getenv("WEBHOOKS_SECRET")
		if secret == "" {
			panic("secret not set")
		}

		var signature = r.Header.Get("X-Hub-Signature")
		var delivery = r.Header.Get("X-GitHub-Delivery")
		var event = r.Header.Get("X-GitHub-Event")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// todo
		}

		var validSig = CheckMAC(body, []byte(signature), []byte(secret))
		if err != nil {
			// todo
		}

		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	http.ListenAndServe(":80", nil)
}

func CheckMAC(message, messageMAC, key []byte) bool {

	if len(messageMAC) != signatureLength || !strings.HasPrefix(messageMAC, signaturePrefix) {
		return false
	}

	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func readYaml() (ret []YamlScript) {

	yamlFile, err := ioutil.ReadFile("scripts/scripts.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, ret)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return ret

}

type YamlScript struct {
	Event  string `json:"event"`
	Branch string `json:"branch"`
	Sender string `json:"sender"`
	Run    string `json:"run"`
}

type CreateEvent struct {
	Ref        string `json:"ref"`
	RefType    string `json:"ref_type"`
	PusherType string `json:"pusher_type"`
}
