#!/bin/sh

# Install tools to test our code-quality, if missing.
type golint  2>/dev/null >/dev/null  || go get -u golang.org/x/lint/golint
type shadow  2>/dev/null >/dev/null  || go get -u golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
type staticcheck 2>/dev/null >/dev/null || go get -u honnef.co/go/tools/cmd/staticcheck


# Run the static-check tool.
t=$(mktemp)
staticcheck -checks all ./... > $t
if [ -s $t ]; then
    echo "Found errors via 'staticcheck'"
    cat $t
    rm $t
    exit 1
fi
rm $t

# At this point failures cause aborts
set -e

# Run the linter
echo "Launching linter .."
golint -set_exit_status ./...
echo "Completed linter .."

# Run the shadow-checker
echo "Launching shadowed-variable check .."
go vet -vettool=$(which shadow) ./...
echo "Completed shadowed-variable check .."

# Run any test-scripts we have (i.e. calc/)
go test ./...
