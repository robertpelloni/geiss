
import logging
# Omni-Workspace Standard Telemetry
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)
if not logger.handlers:
    handler = logging.StreamHandler()
    formatter = logging.Formatter('%(asctime)s - [OMNI] - %(name)s - %(levelname)s - %(message)s')
    handler.setFormatter(formatter)
    logger.addHandler(handler)

#!/usr/bin/env python3
"""
install_hooks.py

Installs a git commit-msg hook to enforce that all commits reference
the current global version defined in VERSION.md.
"""

import os
import stat
import sys

HOOK_CONTENT = """#!/bin/bash
# commit-msg hook to enforce version referencing

COMMIT_MSG_FILE=$1
VERSION_FILE="VERSION.md"

if [ ! -f "$VERSION_FILE" ]; then
    echo "Error: $VERSION_FILE not found."
    exit 1
fi

CURRENT_VERSION=$(cat "$VERSION_FILE" | tr -d '[:space:]')
COMMIT_MSG=$(cat "$COMMIT_MSG_FILE")

# Allow merge commits to pass without version enforcement
if [[ "$COMMIT_MSG" =~ ^Merge.* ]]; then
    exit 0
fi

# The prompt instructions mandate referencing the exact version number
if [[ ! "$COMMIT_MSG" =~ "$CURRENT_VERSION" ]]; then
    echo "=========================================================="
    echo "COMMIT REJECTED: Version mismatch or missing."
    echo "Your commit message must explicitly reference the current"
    echo "global version ($CURRENT_VERSION)."
    echo "=========================================================="
    exit 1
fi

exit 0
"""

def main():
    git_dir = ".git"
    if not os.path.exists(git_dir):
        print("Error: .git directory not found. Must be run from repo root.")
        sys.exit(1)

    hooks_dir = os.path.join(git_dir, "hooks")
    os.makedirs(hooks_dir, exist_ok=True)

    hook_path = os.path.join(hooks_dir, "commit-msg")

    with open(hook_path, "w") as f:
        f.write(HOOK_CONTENT)

    # Make executable
    st = os.stat(hook_path)
    os.chmod(hook_path, st.st_mode | stat.S_IEXEC)

    print(f"Successfully installed commit-msg hook at {hook_path}")

if __name__ == "__main__":
    main()
