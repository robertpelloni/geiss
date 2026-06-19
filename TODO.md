# Global TODO

This document tracks granular short-term features, immediate bug fixes, and explicit tasks across the Omni-Workspace.

## Immediate Tasks (Current Sprint)
- [x] Investigate the state of the `/scripts` directory. Check if `update_repos_v5.py`, `generate_dashboard.py`, and `prune_broken_submodules.py` actually exist or need to be implemented.
- [x] Implement the missing scripts (`update_repos_v5.py`, `generate_dashboard.py`, `prune_broken_submodules.py`).
- [x] Perform a test run of the recursive submodule update process across the actual ecosystem (once submodules are added).
- [x] Review any existing code files outside of `/docs` to identify immediate refactoring needs or missing UI components. (N/A - Repo is empty of application code).
- [ ] Audit newly injected submodules for missing UI components.
- [x] Implement Git commit hook to enforce version referencing.
- [x] Establish initial unit testing frameworks for the core Python orchestration scripts.

## Backlog
- [ ] Audit `bobmani` rhythm game engine suite for undocumented features.
- [ ] Audit `aios` and `borg` projects to map out their current interaction boundaries.
- [x] Implement a system to auto-generate `SUBMODULE_DASHBOARD.md`.
- [x] Verify that Git hooks are correctly configured to prevent commits without version increments.
- [ ] Develop logic for the automated conflict-resolution handlers.
