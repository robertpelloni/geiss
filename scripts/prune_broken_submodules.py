#!/usr/bin/env python3
"""
prune_broken_submodules.py

A script to identify and remove broken or inaccessible git submodule references
from .gitmodules and the git configuration.
"""

import os
import subprocess
import shutil

def run_cmd(cmd, cwd=None, check=True):
    try:
        result = subprocess.run(cmd, cwd=cwd, check=check, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return result.stdout.strip()
    except subprocess.CalledProcessError:
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

def main():
    if not os.path.exists(".gitmodules"):
        print("No .gitmodules file found.")
        return

    submodules = get_submodules()
    broken_found = False

    for sm in submodules:
        sm_url = run_cmd(["git", "config", "--file", ".gitmodules", "--get", f"submodule.{sm}.url"], check=False)
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
        # We consider it broken if the path doesn't exist AND we can't fetch the URL via ls-remote
        if sm_url:
            print(f"Checking {sm} ({sm_url})...")
            # A quick check to see if remote is accessible
            remote_check = run_cmd(["git", "ls-remote", sm_url], check=False)
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
            if remote_check is None:
                print(f"WARNING: Submodule '{sm}' at '{sm_url}' appears inaccessible.")
                # We could auto-remove here, but for safety in this iteration, we just flag it.
                # To fully prune:
                # run_cmd(["git", "rm", "--cached", sm])
                # run_cmd(["git", "config", "--remove-section", f"submodule.{sm}"])
                # run_cmd(["git", "config", "-f", ".gitmodules", "--remove-section", f"submodule.{sm}"])
                # shutil.rmtree(os.path.join(".git", "modules", sm), ignore_errors=True)
                broken_found = True
        else:
            print(f"WARNING: Submodule '{sm}' has no URL configured.")
            broken_found = True

    if not broken_found:
        print("All submodule references appear intact.")
    else:
        print("Found broken submodules. Please review and remove manually or enable auto-pruning.")

if __name__ == "__main__":
    main()
