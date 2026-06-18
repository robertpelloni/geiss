# Deployment & Environment Setup

## Overview
This document outlines the procedures for setting up the Omni-Workspace environment and deploying the various sub-projects it manages. Given the complexity of the ecosystem (100+ nested repositories), strict adherence to these steps is required.

## Prerequisites
*   Git (latest version)
*   Python 3.x (for orchestration scripts)
*   *Note: Individual submodules may have their own language/framework dependencies (Node.js, C++, Rust, etc.). Refer to their respective `DEPLOY.md` or `README.md` files.*

## Initial Setup & Synchronization
1.  **Clone the Root Repository:**
    ```bash
    git clone <repository_url> omni-workspace
    cd omni-workspace
    ```
2.  **Initialize and Update Submodules (Crucial):**
    The standard `git submodule update --init --recursive` might be insufficient for this complex setup.
    You MUST run the designated Python synchronization script:
    ```bash
    python scripts/update_repos_v5.py
    ```
    *This script recursively syncs all submodules, handles nested structures, and ensures no submodule is left in a detached HEAD state.*
3.  **Sanitize Git State:**
    Before any major operations, ensure the local repo is in sync with the server:
    *   Fetch and pull the root repo.
    *   Sync upstream parent forks.
    *   Resolve any conflicts intelligently.
    *   Ensure all personal feature branches are updated with `main`.

## Maintenance & Pruning
*   **Clean Submodules:** To maintain the integrity of `.gitmodules` and remove broken links, run:
    ```bash
    python scripts/prune_broken_submodules.py
    ```

## Dashboard Generation
*   **Update Views:** To regenerate the global status views, run:
    ```bash
    python scripts/generate_dashboard.py
    ```
    This updates `SUBMODULE_DASHBOARD.md` and related structural documents.

## Deployment Processes
*(Specific deployment pipelines for `aios`, `bobmani`, `fwber`, etc., should be documented here as they are developed or imported from their respective repositories. Currently, rely on the automated synchronization scripts as the primary "deployment" mechanism for the workspace structure.)*
