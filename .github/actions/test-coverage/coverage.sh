# This code is sourced from the go-coverage-report Repository by ncruces.
# Original code: https://github.com/ncruces/go-coverage-report
#
# MIT License
#
# Copyright (c) 2023 Nuno Cruces
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

#!/usr/bin/env bash
set -euo pipefail

# This is a simple script to generate an HTML coverage report,
# and SVG badge for your Go project.
#
# It's meant to be used manually or as a pre-commit hook.
#
# Place it some where in your code tree and execute it.
# If your tests pass, next to the script you'll find
# the coverage.html report and coverage.svg badge.
#
# You can add the badge to your README.md as such:
#  [![Go Coverage](PATH_TO/coverage.svg)](https://raw.githack.com/URL/coverage.html)
#
# Visit https://raw.githack.com/ to find the correct URL.
#
# To have the script run as a pre-commmit hook,
# symlink the script to .git/hooks/pre-commit:
#
#  ln -s PATH_TO/coverage.sh .git/hooks/pre-commit
#
# Or, if you have other pre-commit hooks,
# call it from your main hook.

# Get the script's directory after resolving a possible symlink.
SCRIPT_DIR="$(dirname -- "$(readlink -f "${BASH_SOURCE[0]}")")"

OUT_DIR="${1-$SCRIPT_DIR}"
OUT_FILE="$(mktemp)"

# Get coverage for all packages in the current directory; store next to script.
go test -short ./...  -coverprofile "$OUT_FILE"

if [[ "${INPUT_REPORT-true}" == "true" ]]; then
	# Create an HTML report; store next to script.
	go tool cover -html="$OUT_FILE" -o "$OUT_DIR/coverage.html"
fi

# Extract total coverage: the decimal number from the last line of the function report.
COVERAGE=$(go tool cover -func="$OUT_FILE" | tail -1 | grep -Eo '[0-9]+\.[0-9]')

echo "coverage: $COVERAGE% of statements"

date "+%s,$COVERAGE" >> "$OUT_DIR/coverage.log"
sort -u -o "$OUT_DIR/coverage.log" "$OUT_DIR/coverage.log"

# Pick a color for the badge.
if awk "BEGIN {exit !($COVERAGE >= 90)}"; then
	COLOR=brightgreen
elif awk "BEGIN {exit !($COVERAGE >= 80)}"; then
	COLOR=green
elif awk "BEGIN {exit !($COVERAGE >= 70)}"; then
	COLOR=yellowgreen
elif awk "BEGIN {exit !($COVERAGE >= 60)}"; then
	COLOR=yellow
elif awk "BEGIN {exit !($COVERAGE >= 50)}"; then
	COLOR=orange
else
	COLOR=red
fi

# Download the badge; store next to script.
curl -s "https://img.shields.io/badge/coverage-$COVERAGE%25-$COLOR" > "$OUT_DIR/coverage.svg"

if [[ "${INPUT_CHART-false}" == "true" ]]; then
	# Download the chart; store next to script.
	curl -s -H "Content-Type: text/plain" --data-binary "@$OUT_DIR/coverage.log" \
		https://go-coverage-report.nunocruces.workers.dev/chart/ > \
		"$OUT_DIR/coverage-chart.svg"
fi

# When running as a pre-commit hook, add the report and badge to the commit.
if [[ -n "${GIT_INDEX_FILE-}" ]]; then
	git add "$OUT_DIR/coverage.html" "$OUT_DIR/coverage.svg"
fi
