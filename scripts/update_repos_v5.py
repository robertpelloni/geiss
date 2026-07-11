
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
update_repos_v5.py

A script to recursively update git submodules and ensure they are on their default branches,
preventing detached HEAD states, as required by Omni-Workspace directives.
"""

import os
import subprocess
import sys

def run_cmd(cmd, cwd=None):
    """Run a command and return its output."""
    try:
        result = subprocess.run(cmd, cwd=cwd, check=True, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        print(f"Error running {' '.join(cmd)} in {cwd or os.getcwd()}:")
        print(e.stderr)
        return None

def get_submodules(cwd=None):
    """List all submodule paths relative to the current directory."""
    output = run_cmd(["git", "config", "--file", ".gitmodules", "--get-regexp", "path"], cwd)
    if not output:
        return []

    paths = []
    for line in output.splitlines():
        # Line format: submodule.<name>.path <path>
        parts = line.split(" ", 1)
        if len(parts) == 2:
            paths.append(parts[1])
    return paths

def get_default_branch(repo_path):
    """Determine the default branch (main or master) for a repository."""
    # Try querying remote origin HEAD
    output = run_cmd(["git", "remote", "show", "origin"], cwd=repo_path)
    if output:
        for line in output.splitlines():
            if "HEAD branch" in line:
                return line.split(":")[1].strip()

    # Fallback checking local branches
    branches = run_cmd(["git", "branch"], cwd=repo_path)
    if branches:
        if "main" in branches:
            return "main"
        if "master" in branches:
            return "master"
    return "main" # Final fallback

def update_submodule(submodule_path, root_dir):
    """Update a specific submodule, checkout its default branch, and pull."""
    full_path = os.path.join(root_dir, submodule_path)
    print(f"\n--- Processing Submodule: {submodule_path} ---")

    # Initialize and update if not done
    run_cmd(["git", "submodule", "update", "--init", "--recursive", submodule_path], cwd=root_dir)

    if not os.path.exists(os.path.join(full_path, ".git")):
         print(f"Warning: {submodule_path} does not appear to be a valid git repository. Skipping.")
         return

    # Find the default branch
    default_branch = get_default_branch(full_path)
    print(f"Detected default branch: {default_branch}")

    # Checkout default branch
    print(f"Checking out {default_branch}...")
    run_cmd(["git", "checkout", default_branch], cwd=full_path)

    # Pull latest changes
    print("Pulling latest changes...")
    run_cmd(["git", "pull", "origin", default_branch], cwd=full_path)

    # Recursively handle submodules within this submodule
    nested_submodules = get_submodules(full_path)
    for nested in nested_submodules:
        update_submodule(nested, full_path)

def main():
    root_dir = os.getcwd()

    # Ensure we are in a git repository
    if not os.path.exists(".git"):
        print("Error: Must be run from the root of a git repository.")
        sys.exit(1)

    print("Initializing global submodule update...")

    # Basic init for root
    run_cmd(["git", "submodule", "update", "--init", "--recursive"])

    submodules = get_submodules(root_dir)

    if not submodules:
        print("No submodules found in .gitmodules.")
        return

    for sm_path in submodules:
        update_submodule(sm_path, root_dir)

    print("\nAll submodules updated successfully.")

if __name__ == "__main__":
    main()
