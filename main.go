package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const signaturePrefix = "sha1="
const signatureLength = 45

var events = []string{"wildcard", "check_run", "check_suite", "commit_comment", "create", "delete", "deployment", "deployment_status", "fork", "github_app_authorization", "gollum", "installation", "installation_repositories", "issue_comment", "issues", "label", "marketplace_purchase", "member", "membership", "milestone", "organization", "org_block", "page_build", "ping", "project_card", "project_column", "project", "public", "pull_request_review_comment", "pull_request_review", "pull_request", "push", "repository", "repository_import", "repository_vulnerability_alert", "release", "security_advisory", "status", "team", "team_add", "watch"}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var secret = os.Getenv("WEBHOOKS_SECRET")
		if secret == "" {
			respond500(w, "secret not found in environment")
			return
		}

		//var delivery = r.Header.Get("X-GitHub-Delivery")
		var signature = r.Header.Get("X-Hub-Signature")
		var event = r.Header.Get("X-GitHub-Event")

		if !sliceHasString(events, event) {
			respond500(w, "invalid event type")
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respond500(w, "can't read body")
			return
		}

		var validSig = checkMAC(body, []byte(signature), []byte(secret))
		if !validSig {
			respond500(w, "invalid secret")
			return
		}

		scripts, err := readYaml()
		if err != nil {
			respond500(w, "can't read yaml: ")
			return
		}

		for _, v := range scripts {
			out, err := runCommand(v.Run)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(string(out))
			}
		}

		respond200(w)
	})

	http.ListenAndServe(":8099", nil)
}

func respond500(w http.ResponseWriter, message string) {
	fmt.Println(strconv.Itoa(500) + message)
	w.WriteHeader(500)
	w.Write([]byte(message))
}

func respond200(w http.ResponseWriter) {
	fmt.Println(200)
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func runCommand(cmd string) (out []byte, err error) {

	fmt.Println("Executing: " + cmd)
	//return
	return exec.Command("sh", "-c", cmd).Output()
}

func checkMAC(message, messageMAC, key []byte) bool {

	if len(messageMAC) != signatureLength || !strings.HasPrefix(string(messageMAC), signaturePrefix) {
		return false
	}

	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signaturePrefix+expectedMAC), messageMAC)
}

func readYaml() (ret []YamlScript, err error) {

	yamlFile, err := ioutil.ReadFile("scripts/scripts.yaml")
	if err != nil {
		return ret, err
	}

	err = yaml.Unmarshal(yamlFile, &ret)
	if err != nil {
		return ret, err
	}

	return ret, err

}

func sliceHasString(slice []string, i string) bool {
	for _, v := range slice {
		if v == i {
			return true
		}
	}
	return false
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
