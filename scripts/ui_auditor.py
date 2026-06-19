#!/usr/bin/env python3
"""
ui_auditor.py

A script to scan submodules for backend logic files and attempt to verify
if corresponding UI representations exist, as mandated by the implementation guidelines.
"""

import os
import sys

# Simplified heuristic extensions for backend vs frontend
BACKEND_EXTS = {'.py', '.go', '.rs', '.java', '.c', '.cpp', '.rb'}
FRONTEND_EXTS = {'.html', '.jsx', '.tsx', '.vue', '.css', '.scss', '.js', '.ts'}

def get_subdirectories(root_dir):
    return [d for d in os.listdir(root_dir) if os.path.isdir(os.path.join(root_dir, d)) and not d.startswith('.')]

def scan_directory_for_extensions(directory, extensions):
    """Recursively counts files matching specific extensions."""
    count = 0
    for root, _, files in os.walk(directory):
        if '.git' in root or 'node_modules' in root or 'venv' in root:
            continue
        for file in files:
            if any(file.endswith(ext) for ext in extensions):
                count += 1
    return count

def audit_submodule(sm_path):
    """
    Audits a submodule to determine if it has a severe imbalance of backend to frontend code.
    This is a basic heuristic to flag projects needing UI work.
    """
    backend_count = scan_directory_for_extensions(sm_path, BACKEND_EXTS)
    frontend_count = scan_directory_for_extensions(sm_path, FRONTEND_EXTS)

    status = "OK"
    if backend_count > 0 and frontend_count == 0:
         status = "WARNING: Backend logic detected with ZERO frontend files."
    elif backend_count > 10 and frontend_count < 2:
         status = "WARNING: High backend-to-frontend ratio. UI may be incomplete."

    return {
        "submodule": sm_path,
        "backend_files": backend_count,
        "frontend_files": frontend_count,
        "status": status
    }

def main():
    root_dir = os.getcwd()
    submodules = get_subdirectories(root_dir)

    # Filter out known non-submodule directories
    submodules = [sm for sm in submodules if sm not in ('scripts', 'docs', 'tests')]

    if not submodules:
        print("No submodules found to audit.")
        return

    print("Starting UI Coverage Audit...\n")

    for sm in submodules:
        result = audit_submodule(sm)
        print(f"[{result['submodule']}]")
        print(f"  Backend Files: {result['backend_files']}")
        print(f"  Frontend Files: {result['frontend_files']}")
        print(f"  Status: {result['status']}\n")

if __name__ == "__main__":
    main()
