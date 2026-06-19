# Changelog

All notable changes to the Omni-Workspace ecosystem will be documented in this file.

The format is loosely based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres strictly to the versioning defined in the global `VERSION.md` file.

## [1.0.4] - Phase 2 Orchestration Scripts
### Added
- Created `scripts/drift_detection_daemon.py` to identify if submodules are out of sync with upstream.
- Created `scripts/task_router.py` as a foundational framework for assigning AI models to tasks based on capabilities.
- Added corresponding unit tests in `tests/`.

## [1.0.3] - Conclude Initialization Phase
### Added
- Concluded the "Initialization & Context Restoration" phase.
- Synced with upstream and prepared `HANDOFF.md` for Phase 2 (Submodule Injection).

## [1.0.2] - Git Hooks & SOP Implementation
### Added
- Created `scripts/install_hooks.py` to enforce explicit versioning in git commits via a `commit-msg` hook.
- Formalized `docs/AGENT_HANDOFF_SOP.md` to ensure autonomous agents maintain context during handoffs.
- Updated `HANDOFF.md`, `ROADMAP.md`, and `TODO.md` to map out the next phase.

## [1.0.1] - Add Testing Framework
### Added
- Created `tests/` directory and implemented initial Python unit tests for orchestration scripts.
- Updated `TODO.md` to reflect the completion of the testing task.

## [1.0.0] - Initial Setup
### Added
- Created foundational Omni-Workspace documentation files (`VISION.md`, `MEMORY.md`, `DEPLOY.md`, `IDEAS.md`, `ROADMAP.md`, `TODO.md`).
- Established `VERSION.md` as the single source of truth for build tracking.
- Initialized `CHANGELOG.md` to track global ecosystem changes.
- Defined autonomous agent protocols and handoff procedures based on `UNIVERSAL_LLM_INSTRUCTIONS.md`.
