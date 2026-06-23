import os

def generate_interaction_map():
    print("Generating AIOS and Borg interaction map...")

    aios_files = os.listdir("aios")
    borg_files = os.listdir("borg")

    report = "# AIOS & Borg Interaction Boundaries\n\n"
    report += "**Date:** 2026-06-23\n\n"
    report += "## Overview\n"
    report += "This document maps the interaction boundaries between the AI Operating System (`aios`) and the Borg meta-orchestrator (`borg`).\n\n"

    report += "## AIOS Submodule Contents\n"
    for f in aios_files:
        if f != ".git":
            report += f"- `{f}`\n"

    report += "\n## Borg Submodule Contents\n"
    for f in borg_files:
        if f != ".git":
            report += f"- `{f}`\n"

    report += "\n## Interaction Boundaries\n"
    report += "Based on current structural analysis, the `aios` repository handles local agent-level tasks (`main_aios.py`), while `borg` acts as the meta-orchestrator (`core_borg.py`). The communication currently relies on shared standard telemetry `[OMNI]` log tracking injected across both.\n"

    with open("docs/INTERACTIONS.md", "w") as f:
        f.write(report)

    print("Interaction map generated at docs/INTERACTIONS.md")

if __name__ == "__main__":
    generate_interaction_map()
