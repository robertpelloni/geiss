import os

def audit_directory(path):
    print(f"--- Auditing: {path} ---")
    files = os.listdir(path)
    if not files:
        print("Directory is empty.")
        return

    for file in files:
        full_path = os.path.join(path, file)
        if os.path.isdir(full_path) and file != '.git':
            audit_directory(full_path)
        elif os.path.isfile(full_path):
            with open(full_path, 'r', encoding='utf-8', errors='ignore') as f:
                content = f.read()
                print(f"File: {file} | Size: {len(content)} bytes")

audit_directory("bobmani")
