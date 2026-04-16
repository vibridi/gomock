#!/bin/bash

export GOTOOLCHAIN="go1.26.2+auto"

COV=$(go tool cover -func=coverage.out | fgrep total | awk '{print $3}' | cut -d. -f1)

if [ "$COV" -le 49 ]; then
  COLOR="red"
elif [ "$COV" -le 79 ]; then
  COLOR="yellow"
else
  COLOR="green"
fi

sed -i '' -E "s/coverage-[0-9]+%25-[a-z]+/coverage-${COV}%25-${COLOR}/g" README.md
