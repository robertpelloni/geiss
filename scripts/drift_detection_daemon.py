#!/usr/bin/env python3
"""
drift_detection_daemon.py

A script to detect if submodules have drifted (fallen behind) from their upstream remote branches.
"""

import os
import subprocess

def run_cmd(cmd, cwd=None, check=True):
    try:
        result = subprocess.run(cmd, cwd=cwd, check=check, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        print(f"Error running {' '.join(cmd)} in {cwd or os.getcwd()}:")
        print(e.stderr)
        return None

def get_submodules():
    output = run_cmd(["git", "config", "--file", ".gitmodules", "--get-regexp", "path"], check=False)
    if not output:
        return []
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
    paths = []
    for line in output.splitlines():
        parts = line.split(" ", 1)
        if len(parts) == 2:
            paths.append(parts[1])
    return paths

def check_drift(sm_path):
    if not os.path.exists(os.path.join(sm_path, ".git")):
        return f"{sm_path}: Not initialized"

    # Fetch latest from remote for the submodule
    run_cmd(["git", "fetch", "origin"], cwd=sm_path, check=False)

    # Get local commit hash
    local_hash = run_cmd(["git", "rev-parse", "HEAD"], cwd=sm_path)
    if not local_hash:
        return f"{sm_path}: Could not determine local HEAD"

    # Get remote commit hash for the current tracking branch (usually origin/main or origin/master)
    # This simplified version assumes origin/HEAD exists or fallback to branch tracking
    branch_output = run_cmd(["git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"], cwd=sm_path, check=False)
    remote_branch = branch_output if branch_output else "origin/HEAD"

    remote_hash = run_cmd(["git", "rev-parse", remote_branch], cwd=sm_path, check=False)

    if not remote_hash:
        return f"{sm_path}: Could not determine remote {remote_branch}"

    if local_hash != remote_hash:
        return f"{sm_path}: DRIFT DETECTED (Local: {local_hash[:7]} != Remote: {remote_hash[:7]})"
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
    return f"{sm_path}: In sync"

def main():
    print("Starting drift detection across submodules...")
    submodules = get_submodules()
    if not submodules:
        print("No submodules found.")
        return

    for sm in submodules:
        print(check_drift(sm))

if __name__ == "__main__":
    main()
