# System Memory

## Architectural Observations
*   **Omni-Workspace Paradigm:** The repository acts as a monolithic management layer for a vast array of submodules. The architecture relies heavily on scripts (e.g., `scripts/update_repos_v5.py`, `scripts/generate_dashboard.py`) to orchestrate tasks across boundaries.
*   **Decentralized Codebases, Centralized Control:** While the actual logic of the applications (games, web apps, AI OS) lives in separate repositories (submodules), the orchestration, synchronization rules, and global documentation live here at the root.

## Codebase Traits
*   **Heavy Automation Reliance:** Python scripts in the `scripts/` directory are crucial for daily operations. If these scripts fail, the ecosystem risks desynchronization.
*   **Universal Instructions:** The `docs/UNIVERSAL_LLM_INSTRUCTIONS.md` acts as the definitive guide for all AI agents. It dictates behavior, merging strategies, and handoff procedures.
*   **Model Specialization:** The system expects different LLMs to perform specific roles (Claude for architecture/docs, Gemini for speed/large scans, GPT for code/tests). This implies a pipeline where tasks are routed or where models hand off context via `HANDOFF.md`.

## Design Preferences & Mandates
*   **Aggressive Preservation:** "Never Lose Features." Merges must prioritize retaining local or new feature code. Conflicts must be resolved intelligently, not forcefully overridden.
*   **Continuous Autonomous Execution:** Agents are expected to operate on "autopilot," proceeding to the next task automatically without waiting for user confirmation.
*   **Explicit Versioning:** A strict versioning protocol is enforced. `VERSION.md` is the single source of truth, and every build/bump requires a commit with a specific message format referencing the version.
*   **Complete UI Wiring:** Backend features must be fully represented in the UI. No hidden or inaccessible logic.
*   **In-Depth Documentation:** Code comments must focus on *why*, structural side effects, and optimization choices. Self-evident comments are discouraged.
