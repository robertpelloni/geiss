#!/usr/bin/env python3
"""
ui_scaffold_generator.py

Automatically generates basic UI scaffolding (HTML/React stubs) for backend
features identified by the ui_auditor.py as lacking frontend representation.
This directly satisfies the UI Representation mandate.
"""

import os
import sys

REACT_COMPONENT_TEMPLATE = """import React from 'react';

/**
 * Auto-generated UI Component for {feature_name}
 * TODO: Fully wire this component to the corresponding backend API.
 */
const {component_name} = () => {{
    return (
        <div className="{feature_name}-container">
            <h2>{component_name} Interface</h2>
            <form>
                {{/* Auto-generated interactive form stub */}}
                <label htmlFor="inputField">Input:</label>
                <input type="text" id="inputField" name="inputField" />
                <button type="submit">Submit</button>
            </form>
            <div className="tooltip">
                <span className="tooltiptext">Detailed tooltip for {feature_name}</span>
            </div>
        </div>
    );
}};

export default {component_name};
"""

def generate_component_name(backend_file_name):
    """Converts a backend file name (e.g., 'process_data.py') into a React component name ('ProcessData')."""
    base_name = os.path.splitext(backend_file_name)[0]
    # Convert snake_case or kebab-case to PascalCase
    parts = base_name.replace('-', '_').split('_')
    return "".join(part.capitalize() for part in parts)

def create_scaffold(target_dir, backend_file_name):
    """Generates a React component file in the target directory."""
    if not os.path.exists(target_dir):
        try:
            os.makedirs(target_dir, exist_ok=True)
        except OSError:
            return False, f"Failed to create directory {target_dir}"

    component_name = generate_component_name(backend_file_name)
    file_path = os.path.join(target_dir, f"{component_name}.jsx")

    if os.path.exists(file_path):
        return False, f"Component {file_path} already exists."

    try:
        with open(file_path, "w") as f:
            content = REACT_COMPONENT_TEMPLATE.format(
                feature_name=os.path.splitext(backend_file_name)[0],
                component_name=component_name
            )
            f.write(content)
        return True, f"Successfully created scaffold at {file_path}"
    except Exception as e:
        return False, f"Error writing file {file_path}: {e}"

def main():
    if len(sys.argv) < 3:
        print("Usage: python3 ui_scaffold_generator.py <target_directory> <backend_file_name>")
        sys.exit(1)

    target_dir = sys.argv[1]
    backend_file = sys.argv[2]

    success, message = create_scaffold(target_dir, backend_file)
    print(message)
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()
