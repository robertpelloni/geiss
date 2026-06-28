# Agent Handoff Standard Operating Procedure (SOP)

This document formalizes the procedure for transferring context between autonomous AI agents operating within the Omni-Workspace.

## The Goal
To ensure continuous autonomous execution ("autopilot") without human intervention, maintaining state and architectural context across context windows or model shifts (e.g., Gemini to Claude to GPT).

## Procedure

1. **Session Wrap-Up & Summarization:**
   Before ending a session or execution loop, the active agent must synthesize its findings.
   *   What structural shifts occurred?
   *   What roadblocks were encountered?
   *   What is the *immediate* next step?

2. **Update `HANDOFF.md`:**
   Overwrite `HANDOFF.md` in the repository root with the synthesized summary.
   *   The format must include: `## Context`, `## Actions Taken`, `## Findings & System Memories`, and `## Next Steps for Successor Agent`.
   *   *(Future Enhancement: Archive past handoffs in `logs/handoffs/YYYY-MM-DD-model.md`)*

3. **Update Global Status:**
   *   Check off completed tasks in `TODO.md` and `ROADMAP.md`.
   *   Append new discoveries to `MEMORY.md`.

4. **Version & Commit:**
   *   If any files were modified, bump `VERSION.md`.
   *   Commit changes explicitly referencing the new version (enforced by Git hooks).
   *   Run `git push` (including all submodules if modified).

5. **Successor Initialization:**
   The successor agent must begin its session by reading `UNIVERSAL_LLM_INSTRUCTIONS.md` and `HANDOFF.md` to instantly restore context before taking any action.
