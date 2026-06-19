#!/usr/bin/env python3
"""
task_router.py

A foundational script for analyzing repository needs and assigning tasks to specific
LLM agents based on predefined capabilities (as outlined in UNIVERSAL_LLM_INSTRUCTIONS).
"""

import sys

# Define model capabilities based on the global mandate
MODEL_CAPABILITIES = {
    "Claude": ["architecture", "planning", "documentation", "refactoring", "system_understanding"],
    "Gemini": ["speed", "performance", "large_context", "scanning", "scripting"],
    "GPT": ["code_generation", "unit_testing", "algorithm_implementation"]
}

def analyze_task_requirements(task_description):
    """
    Analyzes a task string and returns a list of matched keywords.
    """
    task_lower = task_description.lower()
    requirements = []

    # Simple keyword matching logic for the foundation
    if any(kw in task_lower for kw in ["doc", "plan", "architect", "refactor"]):
        requirements.extend(["architecture", "documentation", "planning"])

    if any(kw in task_lower for kw in ["scan", "performance", "fast", "scripting", "all repos"]):
        requirements.extend(["speed", "scripting", "large_context"])

    if any(kw in task_lower for kw in ["code", "test", "implement", "algorithm", "function", "script"]):
        requirements.extend(["code_generation", "unit_testing"])

    return requirements

def route_task(task_description):
    """
    Routes a task to the most appropriate model based on matched requirements.
    """
    if not task_description:
        return "Unknown", "No task description provided."

    requirements = analyze_task_requirements(task_description)

    if not requirements:
        return "Claude", "Default routing applied due to ambiguous task."

    # Scoring models
    scores = {model: 0 for model in MODEL_CAPABILITIES}

    for req in requirements:
        for model, capabilities in MODEL_CAPABILITIES.items():
            if req in capabilities:
                scores[model] += 1

    # Find model with highest score
    best_model = max(scores, key=scores.get)

    # If tie, default to Claude (Architect)
    if scores[best_model] == 0:
        return "Claude", "No strong match, defaulting to Architect."

    return best_model, f"Matched requirements: {', '.join(requirements)}"

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 task_router.py 'Task description here'")
        sys.exit(1)

    task_desc = sys.argv[1]
    model, reason = route_task(task_desc)

    print(f"Task: '{task_desc}'")
    print(f"Assigned Model: {model}")
    print(f"Reasoning: {reason}")

if __name__ == "__main__":
    main()
