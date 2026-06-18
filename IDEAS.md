# Future Ideas & Pivots

This document is a sandbox for aggressive ideas, creative pivots, structural overhauls, and massive feature expansions for the Omni-Workspace ecosystem.

## Workspace Architecture Overhaul Ideas
*   **Decentralized Agent Coordination:** Move away from purely script-based orchestration (`scripts/`) to an event-driven message bus where AI agents subscribe to repository events (commits, issue creation, submodule drift) and react autonomously in real-time.
*   **Virtual File System (VFS) for Submodules:** Instead of standard git submodules, implement a custom VFS layer that abstracts the physical location of the repositories, allowing for faster cross-repo refactoring and unified codebase searches.
*   **Language Porting Targets:**
    *   Rewrite performance-critical orchestration Python scripts into Rust for speed and concurrency during massive (100+ repo) sync operations.
    *   Consider compiling core `bobmani` rhythm game logic to WebAssembly to allow unified web-based testing interfaces directly within the Omni-Workspace documentation.

## Feature Expansions
*   **"Time Machine" Build System:** The ability to instantly revert the entire Omni-Workspace (all submodules) to an exact state corresponding to a specific `VERSION.md` timestamp, ensuring reproducible builds across the massive ecosystem.
*   **Automated UI Generation:** For any backend feature detected in a submodule, deploy an AI agent specifically tasked with generating a corresponding web-based UI component, ensuring zero "hidden" features.
*   **Self-Writing Tests:** A daemon that constantly monitors the workspace, identifies untested code paths, writes tests, runs them, and creates pull requests completely autonomously.

## Radical Pivots
*   **From Workspace to OS:** Evolve the Omni-Workspace from a git-management hub into a full-fledged "Development Operating System" where the IDE, testing environment, deployment pipeline, and documentation generation are entirely contained within and managed by the AI agents.
