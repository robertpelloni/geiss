#!/usr/bin/env python3
"""
generate_dashboard.py

A script to generate SUBMODULE_DASHBOARD.md, outlining the status, current commit,
and branch of all submodules in the Omni-Workspace.
"""

import os
import subprocess
import datetime

def run_cmd(cmd, cwd=None):
    try:
        result = subprocess.run(cmd, cwd=cwd, check=True, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return result.stdout.strip()
    except subprocess.CalledProcessError:
        return "Unknown"

def get_submodules():
    output = run_cmd(["git", "config", "--file", ".gitmodules", "--get-regexp", "path"])
    if output == "Unknown" or not output:
        return []

    paths = []
    for line in output.splitlines():
        parts = line.split(" ", 1)
        if len(parts) == 2:
            paths.append(parts[1])
    return paths

def get_repo_status(path):
    if not os.path.exists(os.path.join(path, ".git")):
        return "Not Initialized", "N/A", "N/A"

    branch = run_cmd(["git", "rev-parse", "--abbrev-ref", "HEAD"], cwd=path)
    commit = run_cmd(["git", "rev-parse", "--short", "HEAD"], cwd=path)

    # Check if working tree is clean
    status_output = run_cmd(["git", "status", "--porcelain"], cwd=path)
    status = "Clean" if not status_output else "Dirty"

    return branch, commit, status

def main():
    dashboard_file = "SUBMODULE_DASHBOARD.md"
    submodules = get_submodules()

    with open(dashboard_file, "w") as f:
        f.write("# Omni-Workspace Submodule Dashboard\n\n")
        f.write(f"*Last updated: {datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')}*\n\n")

        if not submodules:
            f.write("No submodules detected in `.gitmodules`.\n")
            return

        f.write("| Submodule | Branch | Commit | Status |\n")
        f.write("| --- | --- | --- | --- |\n")

        for sm in sorted(submodules):
            branch, commit, status = get_repo_status(sm)
            f.write(f"| `{sm}` | `{branch}` | `{commit}` | {status} |\n")

    print(f"Generated {dashboard_file}")

if __name__ == "__main__":
    main()
