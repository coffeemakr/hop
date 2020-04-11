#!/bin/bash
ROOT=$(pwd)
ruckd() {
	(cd "$ROOT/server" && go run app/ruckd.go $*)
}

ruck() {
	(cd "$ROOT/cli" && go run app/ruck.go $*)
}
