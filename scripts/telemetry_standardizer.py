#!/usr/bin/env python3
"""
telemetry_standardizer.py

A foundational script to enforce standard logging and telemetry across the Omni-Workspace submodules.
This script scans for known backend languages and suggests/injects standard
Omni-Workspace logging stubs to feed back into the central memory.
"""

import os
import sys

# Standard Omni-Workspace Telemetry Stubs
PYTHON_LOG_STUB = """
import logging
# Omni-Workspace Standard Telemetry
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)
if not logger.handlers:
    handler = logging.StreamHandler()
    formatter = logging.Formatter('%(asctime)s - [OMNI] - %(name)s - %(levelname)s - %(message)s')
    handler.setFormatter(formatter)
    logger.addHandler(handler)
"""

NODEJS_LOG_STUB = """
// Omni-Workspace Standard Telemetry
const omniLogger = {
    info: (msg) => console.log(`[${new Date().toISOString()}] - [OMNI] - INFO - ${msg}`),
    error: (msg) => console.error(`[${new Date().toISOString()}] - [OMNI] - ERROR - ${msg}`)
};
"""

def detect_language(file_path):
    if file_path.endswith('.py'):
        return 'python'
    if file_path.endswith('.js') or file_path.endswith('.ts'):
        return 'nodejs'
    return 'unknown'

def check_file_for_telemetry(file_path):
    """
    Checks if a file already contains the '[OMNI]' standard telemetry signature.
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
            return "[OMNI]" in content
    except Exception:
        return False

def standardize_file(file_path):
    """
    Attempts to inject the standardized telemetry stub into the file.
    Returns True if injected, False if already present or unsupported.
    """
    if check_file_for_telemetry(file_path):
        return False, "Telemetry already present"

    lang = detect_language(file_path)
    stub = ""

    if lang == 'python':
        stub = PYTHON_LOG_STUB
    elif lang == 'nodejs':
        stub = NODEJS_LOG_STUB
    else:
        return False, f"Unsupported language for {file_path}"

    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            original_content = f.read()

        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(stub + "\n" + original_content)

        return True, "Telemetry injected successfully"
    except Exception as e:
        return False, f"Error writing to file: {e}"

def scan_and_standardize(target_dir):
    """Recursively scans a directory to standardize logging."""
    results = []
    for root, _, files in os.walk(target_dir):
        if '.git' in root or 'node_modules' in root or 'venv' in root:
            continue
        for file in files:
            file_path = os.path.join(root, file)
            lang = detect_language(file_path)
            if lang != 'unknown':
                success, msg = standardize_file(file_path)
                results.append((file_path, success, msg))
    return results

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 telemetry_standardizer.py <target_directory>")
        sys.exit(1)

    target_dir = sys.argv[1]

    if not os.path.exists(target_dir):
        print(f"Directory {target_dir} does not exist.")
        sys.exit(1)

    print(f"Starting telemetry standardization scan in {target_dir}...\n")
    results = scan_and_standardize(target_dir)

    for path, success, msg in results:
        status = "SUCCESS" if success else "SKIPPED"
        print(f"[{status}] {path} - {msg}")

if __name__ == "__main__":
    main()
