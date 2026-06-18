# Global TODO

This document tracks granular short-term features, immediate bug fixes, and explicit tasks across the Omni-Workspace.

## Immediate Tasks (Current Sprint)
- [x] Investigate the state of the `/scripts` directory. Check if `update_repos_v5.py`, `generate_dashboard.py`, and `prune_broken_submodules.py` actually exist or need to be implemented.
- [x] Implement the missing scripts (`update_repos_v5.py`, `generate_dashboard.py`, `prune_broken_submodules.py`).
- [ ] Perform a test run of the recursive submodule update process across the actual ecosystem (once submodules are added).
- [ ] Review any existing code files outside of `/docs` to identify immediate refactoring needs or missing UI components.
- [ ] Establish initial unit testing frameworks for the core Python orchestration scripts.

## Backlog
- [ ] Audit `bobmani` rhythm game engine suite for undocumented features.
- [ ] Audit `aios` and `borg` projects to map out their current interaction boundaries.
- [ ] Implement a system to auto-generate `SUBMODULE_DASHBOARD.md`.
- [ ] Verify that Git hooks are correctly configured to prevent commits without version increments.
