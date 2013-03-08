package main

import "fmt"
import "os"
import "os/exec"
import "regexp"
import "strings"

func main() {
  fetch()

  branch := currentBranch()

  if needsForce() {
    if confirm(fmt.Sprintf("\n'%s' needs to be \x1b[31mFORCE PUSHED\x1b[0m, is this OK?", branch)) {
      push(branch, true)
    }
  } else {
    push(branch, false)
  }

  fmt.Printf("\n")
}

func fetch() {
  gitCommand("fetch", "--quiet")
}

func currentBranch() string {
  return gitCommand("rev-parse", "--abbrev-ref", "HEAD")
}

func push(branch string, force bool) {
  if force {
    fmt.Printf("Force pushing to origin/%s\n", branch)
    gitCommand("push", "origin", branch, "--force")

  } else {
    fmt.Printf("Pushing to origin/%s\n", branch)
    gitCommand("push", "origin", branch)
  }
}

func needsForce() bool {
  branch_status := gitCommand("status", "--porcelain", "-b")
  matched, err := regexp.MatchString(`^##.*\[ahead \d+, behind \d]$`, branch_status)
  return err == nil && matched
}

func gitCommand(command ...string) string {
  result, err := exec.Command("git", command...).Output()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  return strings.Trim(string(result), "\n")
}

func confirm(message string) bool {
  fmt.Printf("%s (y/n):", message)

  var confirmed string
  n, err := fmt.Scanln(&confirmed)
  if err != nil || n != 1 {
    return false
  }

  matched, err := regexp.MatchString("(?i)^y(es)?$", confirmed)
  return err == nil && matched
}

