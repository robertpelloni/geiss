# Omni-Workspace Ecosystem Architecture & Patterns (Final Overview)

Based on continuous interaction and exploration of the repository guided by the `UNIVERSAL_LLM_INSTRUCTIONS.md`, here is a comprehensive summary of the project's current architecture, patterns, and decisions:

## 1. Project Context & Vision
The repository is the **Omni-Workspace**, acting as a monolithic management layer for a massive ecosystem of 100+ separate submodules, forks, and independent projects. It governs AI systems, game engines, and web platforms. 
The ultimate vision is a highly autonomous, self-healing environment maintained by collaborating AI agents.

## 2. Global Mandates & Operating Principles
*   **Continuous Autonomous Execution:** Agents must operate on "autopilot," proceeding sequentially without waiting for user confirmation unless a destructive action is imminent.
*   **Never Lose Features:** Merges must intelligently solve conflicts to retain new or local changes. Force pushing over working code is prohibited.
*   **Git Sanitization Protocol:** Upstream syncing and intelligent merging of AI-generated branches is required before new feature development.
*   **Submodule Integrity:** The ecosystem relies on `scripts/update_repos_v5.py` to recursively update submodules, ensuring they are not left in detached HEAD states.

## 3. Core Documentation & Versioning Standards
The centralized documentation structure has been successfully bootstrapped:
*   `UNIVERSAL_LLM_INSTRUCTIONS.md`: The definitive guide for AI agents.
*   `VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `ROADMAP.md`, `TODO.md`: Standardized files outlining goals, architecture, tasks, and setup.
*   `HANDOFF.md`: Used to pass context to successor AI models at the end of a session.
*   `VERSION.md`: The singular, global version string. Every build/bump requires a clean commit explicitly referencing this version.
*   `CHANGELOG.md`: Tracks ecosystem changes aligned with `VERSION.md`.
*   `AGENT_HANDOFF_SOP.md`: Formalizes the context transfer process between agents.

## 4. Current Architectural State & Automation
*   **Decentralized Code, Centralized Control:** Application logic lives in submodules; orchestration and global documentation reside at the root.
*   **Submodule Integration (Phase 3 Completed):** The core repositories (`aios`, `bobmani`, `fwber`, `borg`) are officially tracked as submodules in the root workspace. (Note: These are currently simulated using local file transport due to sandbox constraints).
*   **Orchestration Scripts Implemented:** 
    *   `scripts/update_repos_v5.py`: Recursively syncs all submodules.
    *   `scripts/generate_dashboard.py`: Generates `SUBMODULE_DASHBOARD.md`.
    *   `scripts/prune_broken_submodules.py`: Identifies dead submodule links.
    *   `scripts/install_hooks.py`: Installs Git commit-msg hook for version tracking.
    *   `scripts/drift_detection_daemon.py`: Identifies if submodules have fallen behind upstream remotes.
    *   `scripts/task_router.py`: Assigns tasks to specialized AI models (Claude, Gemini, GPT) based on predefined capabilities.
    *   `scripts/conflict_resolution.py`: A skeleton framework for automating intelligent git conflict resolution that prioritizes feature preservation.
    *   `scripts/ui_auditor.py`: A script to recursively scan submodules for backend logic files and flag missing UI representations.
    *   `scripts/deployment_pipeline.py`: A foundational skeleton to manage build, test, and deploy stages across the submodule ecosystem.
    *   `scripts/ui_scaffold_generator.py`: An autonomous tool to fulfill the UI Quality mandate by auto-generating React/HTML stubs for backend features missing UI components.
    *   `scripts/telemetry_standardizer.py`: A script to inject standardized logging stubs into various backend languages across submodules.
    *   `scripts/system_test_runner.py`: The Phase 4 master sandbox script that proves end-to-end functionality by sequentially executing the core read-only orchestration scripts.
*   **Testing Infrastructure:** A robust python `unittest` framework is established in `tests/` to validate the core logic of all orchestration scripts without running destructive git commands.

## 5. Next Planned Actions
*   **Phase 4 Live Environment Execution:** The environment is primed. The next step is to transition from the sandbox and execute a full system run on the actual codebase utilizing the tools built during this initialization sequence.