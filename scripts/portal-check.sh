#!/bin/bash

cd services/portal

echo "Running npm format..."
yarn format
FORMAT_STATUS=$?

echo "Running npm lint..."
yarn lint
LINT_STATUS=$?

echo "Running npm check..."
yarn check
CHECK_STATUS=$?

# check if any of the npm command exited with a non-zero status
if [ $FORMAT_STATUS -ne 0 -o $LINT_STATUS -ne 0 -o $CHECK_STATUS -ne 0 ]; then
    echo "One or more checks failed. Please fix the errors before committing."
    exit 1
fi

echo "All checks passed. Ready to commit."
exit 0
