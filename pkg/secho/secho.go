package secho

import "log"

func CheckError(err error) {
  if (err != nil) {
    log.Fatalf("Fatal error: %s", err)
  }
}
