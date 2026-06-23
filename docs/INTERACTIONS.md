# AIOS & Borg Interaction Boundaries

**Date:** 2026-06-23

## Overview
This document maps the interaction boundaries between the AI Operating System (`aios`) and the Borg meta-orchestrator (`borg`).

## AIOS Submodule Contents
- `main_aios.py`
- `MainAios.jsx`
- `README.md`

## Borg Submodule Contents
- `core_borg.py`
- `CoreBorg.jsx`
- `README.md`

## Interaction Boundaries
Based on current structural analysis, the `aios` repository handles local agent-level tasks (`main_aios.py`), while `borg` acts as the meta-orchestrator (`core_borg.py`). The communication currently relies on shared standard telemetry `[OMNI]` log tracking injected across both.
