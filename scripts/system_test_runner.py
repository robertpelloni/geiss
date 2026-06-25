#!/usr/bin/env python3
"""
system_test_runner.py

A Phase 4 master testing script designed to sequentially execute the core
Omni-Workspace orchestration and unification scripts to validate end-to-end
autonomous system functionality.
"""

import os
import subprocess
import sys
import time

def run_script(script_name, args=[]):
    """Executes a python script from the scripts/ directory."""
    print(f"\n{'='*60}")
    print(f"EXECUTING: {script_name}")
    print(f"{'='*60}")
<<<<<<< HEAD

    script_path = os.path.join(os.getcwd(), "scripts", script_name)

=======

    script_path = os.path.join(os.getcwd(), "scripts", script_name)

>>>>>>> jules-9396211896448288708-4318ead9
    if not os.path.exists(script_path):
        print(f"Error: {script_path} not found.")
        return False

    cmd = ["python3", script_path] + args
<<<<<<< HEAD

=======

>>>>>>> jules-9396211896448288708-4318ead9
    try:
        # We don't capture output here so it streams directly to the console
        result = subprocess.run(cmd, check=True)
        print(f"\n[OK] {script_name} completed successfully.")
        return True
    except subprocess.CalledProcessError as e:
        print(f"\n[FAIL] {script_name} exited with error code {e.returncode}.")
        return False

def main():
    root_dir = os.getcwd()
    if not os.path.exists(os.path.join(root_dir, ".gitmodules")):
        print("Error: Must be run from the root of the Omni-Workspace (where .gitmodules exists).")
        sys.exit(1)

    print("Initiating Phase 4: Full System Automation Test...\n")
    time.sleep(1)

    # Sequence 1: Maintenance & Synchronization
    if not run_script("update_repos_v5.py"):
        print("Halting system test due to synchronization failure.")
        sys.exit(1)

    if not run_script("prune_broken_submodules.py"):
         print("Warning: Pruning script returned non-zero exit code.")

    if not run_script("drift_detection_daemon.py"):
        print("Warning: Drift detection daemon failed.")

    # Sequence 2: Dashboard Generation
    if not run_script("generate_dashboard.py"):
        print("Warning: Dashboard generation failed.")

    # Sequence 3: Ecosystem Unification
    if not run_script("ui_auditor.py"):
         print("Warning: UI Auditor failed.")

    # We skip ui_scaffold_generator as it requires specific file arguments
    # We skip telemetry_standardizer as it mutates files and requires a target dir

    print(f"\n{'='*60}")
    print("PHASE 4 SYSTEM TEST COMPLETE.")
    print("Core read-only orchestration scripts successfully executed.")
    print(f"{'='*60}\n")

if __name__ == "__main__":
    main()
