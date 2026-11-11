#!/bin/sh
# Prints the CHANGELOG.md section of the given tag (v1.2.3 -> "## [1.2.3]")
# and fails if the section is missing or empty.
set -eu

tag=${1:?usage: release-notes.sh <tag>}
version=${tag#v}

notes=$(awk -v ver="$version" '
	/^## \[/ { if (found) exit; if (index($0, "## [" ver "]") == 1) { found = 1; next } }
	found { print }
' CHANGELOG.md)

if [ -z "$(printf '%s' "$notes" | tr -d '[:space:]')" ]; then
	echo "no CHANGELOG.md section for $version" >&2
	exit 1
fi
printf '%s\n' "$notes"
