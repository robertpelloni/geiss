# Deployment & Environment Setup

## Overview
This document outlines the procedures for setting up the Jules Autopilot Omni-Workspace environment. Given the pivot to a Go-first architecture with a React SPA frontend, strict adherence to these steps is required.

## Prerequisites
*   Git (latest version)
*   Go v1.26.0
*   Node.js v20.20.2
*   pnpm

## Local Development Setup
1.  **Clone the Repository:**
    ```bash
    git clone <repository_url> jules-autopilot
    cd jules-autopilot
    ```
2.  **Frontend Setup:**
    Ensure `ignore-scripts=false` is set in your local `.npmrc` file.
    ```bash
    pnpm install
    pnpm run dev
    ```
3.  **Backend Setup:**
    ```bash
    cd backend-go
    go run main.go
    ```

## Production Build & Run
1.  **Build Frontend:**
    ```bash
    pnpm run build
    ```
    *This generates the static assets in the `dist/` directory.*
2.  **Build and Run Backend:**
    ```bash
    cd backend-go
    go build -o jules-backend
    ./jules-backend
    ```
    *The Go backend serves the API on port 8080 and statically hosts the `../dist` frontend build.*

## Render Deployment Specifications
When deploying to Render, the following environment variables are required:
*   `NODE_VERSION=20.20.2`
*   `GO_VERSION=1.26.0`
*   `CGO_ENABLED=1`
*   `NODE_OPTIONS=--max-old-space-size=4096` (Prevents Vite OOM errors during build)
*   `JULES_API_KEY`

## Legacy Python Orchestration Scripts
Some orchestration scripts in `/scripts` still rely on Python. Ensure Python 3.x is installed to run scripts like `update_repos_v5.py` for submodule syncing.
