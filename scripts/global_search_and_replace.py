
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
global_search_and_replace.py

A utility script to perform global search and replace operations across all submodules
in the Omni-Workspace ecosystem. Designed to facilitate mass refactoring.
"""

import os
import argparse
import subprocess

def get_all_submodules(root_dir):
    """Parses .gitmodules to find all submodule paths."""
    submodules = []
    gitmodules_path = os.path.join(root_dir, '.gitmodules')
    if not os.path.exists(gitmodules_path):
        return submodules

    with open(gitmodules_path, 'r') as f:
        for line in f:
            if line.strip().startswith('path = '):
                path = line.split('=', 1)[1].strip()
                submodules.append(os.path.join(root_dir, path))
    return submodules

def process_file(filepath, search_str, replace_str):
    """Replaces occurrences of search_str with replace_str in a file."""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()

        if search_str in content:
            new_content = content.replace(search_str, replace_str)
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(new_content)
            return True

    except UnicodeDecodeError:
        pass # Skip binary files
    except Exception as e:
        print(f"Error processing {filepath}: {e}")

    return False

def search_and_replace_in_dir(directory, search_str, replace_str, dry_run=False):
    """Recursively walks a directory and performs search/replace."""
    modifications = 0
    for root, _, files in os.walk(directory):
        if '.git' in root:
            continue # Skip git directories

        for file in files:
            filepath = os.path.join(root, file)
            if dry_run:
                try:
                    with open(filepath, 'r', encoding='utf-8') as f:
                        if search_str in f.read():
                            print(f"[DRY RUN] Would modify: {filepath}")
                            modifications += 1
                except UnicodeDecodeError:
                    pass
            else:
                if process_file(filepath, search_str, replace_str):
                    print(f"Modified: {filepath}")
                    modifications += 1

    return modifications

def main():
    parser = argparse.ArgumentParser(description="Global search and replace across Omni-Workspace.")
    parser.add_argument("search_str", help="The string to search for.")
    parser.add_argument("replace_str", help="The string to replace it with.")
    parser.add_argument("--dry-run", action="store_true", help="Print files that would be modified without changing them.")

    args = parser.parse_args()

    root_dir = os.path.abspath(os.path.join(os.path.dirname(__file__), '..'))
    submodules = get_all_submodules(root_dir)

    if not submodules:
        print("No submodules found. Are you running this from the Omni-Workspace root?")
        return

    total_modifications = 0
    print(f"Searching for '{args.search_str}' in {len(submodules)} submodules...")

    for submodule in submodules:
        if os.path.exists(submodule):
            total_modifications += search_and_replace_in_dir(submodule, args.search_str, args.replace_str, args.dry_run)

    if args.dry_run:
        print(f"Dry run complete. {total_modifications} files would be modified.")
    else:
        print(f"Refactoring complete. {total_modifications} files modified.")

if __name__ == "__main__":
    main()
