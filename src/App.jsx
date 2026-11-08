import { useEffect, useState } from 'react'

function App() {
  // Core System Telemetry States
  const [submodules, setSubmodules] = useState([])
  const [anomaly, setAnomaly] = useState(null)
  const [queue, setQueue] = useState([])
  const [logs, setLogs] = useState("")
  const [apiKey, setApiKey] = useState("")
  const [trayStatus, setTrayStatus] = useState(null)

  // Interactive Form States
  const [taskInput, setTaskInput] = useState("")
  const [taskRoutingResult, setTaskRoutingResult] = useState(null)
  const [routingLoading, setRoutingLoading] = useState(false)

  const [conflicts, setConflicts] = useState([])
  const [driftReport, setDriftReport] = useState([])
  const [auditReport, setAuditReport] = useState([])

  // Global Refactoring
  const [searchStr, setSearchStr] = useState("")
  const [replaceStr, setReplaceStr] = useState("")
  const [dryRun, setDryRun] = useState(true)
  const [refactorResult, setRefactorResult] = useState(null)
  const [refactorLoading, setRefactorLoading] = useState(false)

  // Pipeline Generator & Runner
  const [pipelineSubmodule, setPipelineSubmodule] = useState("")
  const [pipelineEnv, setPipelineEnv] = useState("staging")
  const [pipelineResult, setPipelineResult] = useState(null)
  const [pipelineLoading, setPipelineLoading] = useState(false)

  // Status Notification
  const [statusMessage, setStatusMessage] = useState("")
  const [activeTab, setActiveTab] = useState("all") // 'all', 'orchestration', 'submodules', 'telemetry'

  const notify = (msg) => {
    setStatusMessage(msg)
    setTimeout(() => setStatusMessage(""), 4000)
  }

  const fetchTray = () => {
    fetch('/api/system/tray')
      .then(res => res.json())
      .then(data => setTrayStatus(data))
      .catch(err => console.error("Tray status error:", err))
  }

  const fetchLogs = () => {
    fetch('/api/system/logs')
      .then(res => res.json())
      .then(data => setLogs(data.logs || ""))
      .catch(err => console.error("Failed to fetch logs:", err))
  }

  const fetchSubmodules = () => {
    fetch('/system/status')
      .then(res => res.json())
      .then(data => setSubmodules(data || []))
      .catch(err => console.error("Failed to fetch submodules:", err))
  }

  const fetchAnomaly = () => {
    fetch('/api/shadow/diff')
      .then(res => res.json())
      .then(data => setAnomaly(data))
      .catch(err => console.error("Failed to fetch diff logic:", err))
  }

  useEffect(() => {
    fetchSubmodules()
    fetchAnomaly()
    fetchLogs()
    fetchTray()

    // Polling interval for queue telemetry & tray uptime
    const interval = setInterval(() => {
      fetch('/api/queue/telemetry')
        .then(res => res.json())
        .then(data => setQueue(data || []))
        .catch(err => console.error("Failed to fetch queue telemetry:", err))
      fetchTray()
    }, 2500)

    return () => clearInterval(interval)
  }, [])

  // Auto-Fix trigger
  const handleTriggerAutoFix = () => {
    fetch('/api/shadow/autofix', { method: 'POST', headers: { 'X-API-KEY': apiKey } })
      .then(res => res.json())
      .then(data => {
        notify(data.message || "Auto-fix routine finished.")
        fetchAnomaly()
      })
      .catch(() => notify("Auto-fix trigger failed."))
  }

  // Task Route trigger
  const handleRouteTask = (e) => {
    e.preventDefault()
    if (!taskInput.trim()) return
    setRoutingLoading(true)
    fetch('/api/tasks/route', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ task_description: taskInput })
    })
      .then(res => res.json())
      .then(data => {
        setTaskRoutingResult(data)
        setRoutingLoading(false)
        notify(`Task assigned to ${data.assigned_model}`)
      })
      .catch(() => {
        setRoutingLoading(false)
        notify("Failed to route task.")
      })
  }

  // Conflict Resolution check
  const handleCheckConflicts = () => {
    fetch('/api/conflicts/resolve')
      .then(res => res.json())
      .then(data => {
        setConflicts(data || [])
        notify(`Conflict check complete. Found ${data?.length || 0} items.`)
      })
      .catch(() => notify("Error running conflict check."))
  }

  // Drift Detection check
  const handleCheckDrift = () => {
    fetch('/api/system/drift')
      .then(res => res.json())
      .then(data => {
        setDriftReport(data || [])
        notify(`Drift scan complete. ${data?.length || 0} submodules analyzed.`)
      })
      .catch(() => notify("Error scanning git drift."))
  }

  // Prune broken submodules
  const handlePruneSubmodules = () => {
    fetch('/api/system/prune')
      .then(res => res.json())
      .then(data => {
        notify(`Prune scan complete: ${data?.length || 0} modules verified.`)
        fetchSubmodules()
      })
      .catch(() => notify("Error pruning submodules."))
  }

  // Generate Dashboard
  const handleGenerateDashboard = () => {
    fetch('/api/system/dashboard', {
      method: 'POST',
      headers: { 'X-API-KEY': apiKey }
    })
      .then(res => res.json())
      .then(data => notify(data.message || "Dashboard regenerated."))
      .catch(() => notify("Failed to regenerate dashboard. Verify API Key."))
  }

  // Global Search and Replace
  const handleGlobalRefactor = (e) => {
    e.preventDefault()
    if (!searchStr) return
    setRefactorLoading(true)
    fetch('/api/system/refactor', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ search_str: searchStr, replace_str: replaceStr, dry_run: dryRun })
    })
      .then(res => res.json())
      .then(data => {
        setRefactorResult(data)
        setRefactorLoading(false)
        notify(`Refactor execution: ${data.modified_files} files affected.`)
      })
      .catch(() => {
        setRefactorLoading(false)
        notify("Error during search and replace execution.")
      })
  }

  // Run or Generate Pipeline
  const handlePipelineAction = (generate) => {
    if (!pipelineSubmodule) {
      notify("Please select or enter a submodule name.")
      return
    }
    setPipelineLoading(true)
    fetch('/api/system/pipeline', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'X-API-KEY': apiKey },
      body: JSON.stringify({ submodule: pipelineSubmodule, environment: pipelineEnv, generate })
    })
      .then(res => res.json())
      .then(data => {
        setPipelineResult(data)
        setPipelineLoading(false)
        notify(data.message)
      })
      .catch(() => {
        setPipelineLoading(false)
        notify("Pipeline operation failed.")
      })
  }

  // UI Auditor
  const handleRunAudit = () => {
    fetch('/api/system/audit', {
      headers: { 'X-API-KEY': apiKey }
    })
      .then(res => res.json())
      .then(data => {
        setAuditReport(data || [])
        notify(`Audit finished: ${data?.length || 0} submodules analyzed.`)
      })
      .catch(() => notify("UI Auditor failed. Check JULES_API_KEY."))
  }

  // Standardize Telemetry
  const handleStandardizeTelemetry = () => {
    fetch('/api/system/telemetry', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ target_dir: "." })
    })
      .then(res => res.json())
      .then(data => notify(`Telemetry standardized across ${data?.length || 0} repos.`))
      .catch(() => notify("Telemetry standardizer failed."))
  }

  // Graceful Server Shutdown via System Tray
  const handleQuitServers = () => {
    if (window.confirm("Are you sure you want to shut down the Jules Autopilot Go backend and daemons?")) {
      fetch('/api/system/shutdown', { method: 'POST' })
        .then(res => res.json())
        .then(data => {
          alert(data.message)
          fetchTray()
        })
        .catch(() => alert("Shutdown call failed or server already stopped."))
    }
  }

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-blue-600 selection:text-white">
      {/* Top Banner / Toast Notification */}
      {statusMessage && (
        <div className="fixed top-4 right-4 z-50 bg-blue-600 text-white px-5 py-3 rounded-lg shadow-2xl flex items-center space-x-3 transition-transform duration-300 ease-out border border-blue-400">
          <span className="font-bold">⚡ Status:</span>
          <span>{statusMessage}</span>
        </div>
      )}

      {/* Main Top Header: High-Level Autopilot Telemetry & System Tray Control */}
      <header className="bg-slate-900 border-b border-slate-800 px-6 py-4 flex flex-wrap items-center justify-between gap-4 sticky top-0 z-40 shadow-md">
        <div className="flex items-center space-x-4">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-blue-600 via-indigo-500 to-purple-600 flex items-center justify-center font-black text-xl shadow-lg ring-2 ring-blue-500/30">
            J
          </div>
          <div>
            <h1 className="text-xl font-extrabold tracking-tight text-white flex items-center gap-2" title="Unified Autonomous Management Hub for Jules Autopilot and the Omni-Workspace">
              Jules Autopilot
              <span className="text-xs px-2 py-0.5 font-mono uppercase bg-blue-900/60 text-blue-300 border border-blue-700/50 rounded-full" title="Global architecture: Native Go Backend + Pure React 19 SPA">
                Go-First v1.0.39
              </span>
            </h1>
            <p className="text-xs text-slate-400">Autonomous Orchestration • Extreme Telemetry • Omnipresent Tray</p>
          </div>
        </div>

        {/* Global Filter Tabs */}
        <div className="flex items-center bg-slate-950 p-1 rounded-lg border border-slate-800 text-xs font-medium">
          <button
            onClick={() => setActiveTab("all")}
            className={`px-3 py-1.5 rounded-md transition-all ${activeTab === 'all' ? 'bg-blue-600 text-white shadow' : 'text-slate-400 hover:text-white'}`}
            title="Consolidated view displaying all system sections seamlessly"
          >
            All Features
          </button>
          <button
            onClick={() => setActiveTab("orchestration")}
            className={`px-3 py-1.5 rounded-md transition-all ${activeTab === 'orchestration' ? 'bg-blue-600 text-white shadow' : 'text-slate-400 hover:text-white'}`}
            title="Focus specifically on Autonomous Intelligence, AI Routing, and Mass Refactoring"
          >
            Orchestration & AI
          </button>
          <button
            onClick={() => setActiveTab("submodules")}
            className={`px-3 py-1.5 rounded-md transition-all ${activeTab === 'submodules' ? 'bg-blue-600 text-white shadow' : 'text-slate-400 hover:text-white'}`}
            title="Manage git submodules, conflict resolution, drift status, and pipeline runs"
          >
            Ecosystem Repos
          </button>
          <button
            onClick={() => setActiveTab("telemetry")}
            className={`px-3 py-1.5 rounded-md transition-all ${activeTab === 'telemetry' ? 'bg-blue-600 text-white shadow' : 'text-slate-400 hover:text-white'}`}
            title="Inspect real-time queue tasks and live backend system logs"
          >
            Telemetry & Logs
          </button>
        </div>

        {/* System Tray Daemon & Global Quit Server Controls */}
        <div className="flex items-center space-x-3 bg-slate-950 px-4 py-2 rounded-xl border border-slate-800">
          <div className="flex items-center space-x-2" title="OS System Tray Supervisor Daemon status">
            <span className={`w-3 h-3 rounded-full animate-pulse ${trayStatus?.active ? 'bg-emerald-500' : 'bg-rose-500'}`} />
            <div className="text-left">
              <div className="text-[10px] text-slate-400 uppercase tracking-wider font-semibold">System Tray</div>
              <div className="text-xs font-mono font-medium text-slate-200">
                {trayStatus ? `${trayStatus.platform} (${trayStatus.uptime_sec}s up)` : 'Connecting...'}
              </div>
            </div>
          </div>

          <div className="h-6 w-px bg-slate-800" />

          {/* Destructive / Auth API Key Input */}
          <div className="relative">
            <input
              type="password"
              value={apiKey}
              onChange={e => setApiKey(e.target.value)}
              placeholder="API Key..."
              className="w-28 text-xs bg-slate-900 border border-slate-700 px-2 py-1 rounded text-slate-200 focus:outline-none focus:border-blue-500 focus:w-44 transition-all"
              title="Enter JULES_API_KEY to unlock destructive actions (Auto-Fix, Dashboard Regen, UI Audit)"
            />
          </div>

          {/* Safe Quit Servers Button */}
          <button
            onClick={handleQuitServers}
            className="bg-rose-600/20 hover:bg-rose-600 text-rose-300 hover:text-white border border-rose-500/30 px-3 py-1 rounded text-xs font-bold transition-all flex items-center gap-1.5 shadow-sm"
            title="Gracefully terminates the Go backend server, workers, and system tray supervisor safely."
          >
            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Quit Server
          </button>
        </div>
      </header>

      {/* Main Single-Page Dashboard Layout */}
      <main className="flex-1 p-6 space-y-8 max-w-7xl mx-auto w-full">

        {/* SECTION 1: HIGHEST-VALUE ORCHESTRATION & AUTONOMY (Always visible or in Orchestration tab) */}
        {(activeTab === 'all' || activeTab === 'orchestration') && (
          <section className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-xl font-bold text-white flex items-center gap-2">
                  <span className="w-2 h-2 rounded-full bg-blue-500" />
                  Primary Orchestration & Autonomous Intelligence
                </h2>
                <p className="text-xs text-slate-400">High-priority automated decision making, structural diff monitoring, and AI multi-model routing.</p>
              </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Card 1: Shadow Pilot & Anomaly Auto-Fix */}
              <div className="bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-lg flex flex-col justify-between hover:border-slate-700 transition-colors">
                <div>
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-base font-semibold text-amber-400 flex items-center gap-2" title="Background daemon monitoring git diffs for mass automated refactoring anomalies.">
                      Shadow Pilot Sentinel
                    </h3>
                    <span className="text-[10px] uppercase font-mono px-2 py-0.5 bg-amber-950/80 text-amber-300 border border-amber-800 rounded">
                      Realtime Sentinel
                    </span>
                  </div>
                  <p className="text-xs text-slate-300 mb-4">
                    Continuously calculates worktree byte diffs to prevent regressions and accidental code overwrites.
                  </p>

                  {anomaly ? (
                    <div className="space-y-2 bg-slate-950 p-3 rounded-xl border border-slate-800/80 text-xs font-mono">
                      <div className="flex justify-between">
                        <span className="text-slate-400">Diff Size:</span>
                        <span className="text-slate-200 font-bold">{anomaly.diff_size_bytes} bytes</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-slate-400">Last Scan:</span>
                        <span className="text-slate-300">{new Date(anomaly.timestamp).toLocaleTimeString()}</span>
                      </div>
                      {anomaly.warning ? (
                        <div className="mt-2 p-2.5 bg-rose-950/40 border border-rose-500/50 rounded-lg text-rose-300">
                          <p className="font-bold flex items-center gap-1">⚠️ {anomaly.warning}</p>
                          <p className="text-[11px] text-rose-200 mt-1">Automated mass refactor or git divergence flagged.</p>
                        </div>
                      ) : (
                        <div className="mt-1 text-emerald-400 text-[11px]">✓ Worktree delta within safe autonomous thresholds.</div>
                      )}
                    </div>
                  ) : (
                    <p className="text-xs text-slate-500 animate-pulse">Syncing sentinel telemetry...</p>
                  )}
                </div>

                <div className="mt-4 pt-4 border-t border-slate-800/80 flex items-center justify-between">
                  <span className="text-[11px] text-slate-400" title="Trigger self-healing CI fix routine">CI Anomaly Auto-Fix</span>
                  <button
                    onClick={handleTriggerAutoFix}
                    className="bg-amber-600 hover:bg-amber-500 text-slate-950 font-bold px-3 py-1.5 rounded-lg text-xs transition-colors shadow"
                    title="Directly executes CI pipeline auto-fix in backend-go"
                  >
                    Trigger Auto-Fix
                  </button>
                </div>
              </div>

              {/* Card 2: Centralized Task Router */}
              <div className="bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-lg flex flex-col justify-between hover:border-slate-700 transition-colors">
                <div>
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-base font-semibold text-indigo-400 flex items-center gap-2" title="Assigns complex developer tasks to specialized models (Claude 3.5 Sonnet, Gemini 1.5 Pro, GPT-4o).">
                      AI Model Task Router
                    </h3>
                    <span className="text-[10px] uppercase font-mono px-2 py-0.5 bg-indigo-950/80 text-indigo-300 border border-indigo-800 rounded">
                      Model Orchestrator
                    </span>
                  </div>
                  <p className="text-xs text-slate-300 mb-3">
                    Evaluates task keywords and computational requirements to route requests to the best AI model.
                  </p>

                  <form onSubmit={handleRouteTask} className="space-y-2">
                    <div className="relative">
                      <input
                        type="text"
                        value={taskInput}
                        onChange={e => setTaskInput(e.target.value)}
                        placeholder="e.g., Refactor backend telemetry or write tests..."
                        className="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-2 text-xs text-slate-200 placeholder-slate-500 focus:outline-none focus:border-indigo-500 transition-colors"
                        title="Describe any engineering task to discover optimal AI model allocation."
                      />
                    </div>
                    <button
                      type="submit"
                      disabled={routingLoading}
                      className="w-full bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white font-semibold py-1.5 rounded-lg text-xs transition-colors shadow"
                    >
                      {routingLoading ? "Analyzing Task Complexity..." : "Route Task to Model"}
                    </button>
                  </form>
                </div>

                {taskRoutingResult && (
                  <div className="mt-3 p-3 bg-slate-950 border border-indigo-500/30 rounded-xl text-xs space-y-1">
                    <div className="flex justify-between items-center">
                      <span className="text-slate-400">Assigned Model:</span>
                      <span className="font-bold text-indigo-300 font-mono px-2 py-0.5 bg-indigo-950 border border-indigo-700 rounded">
                        {taskRoutingResult.assigned_model}
                      </span>
                    </div>
                    <p className="text-[11px] text-slate-300 mt-1"><span className="text-slate-400">Rationale:</span> {taskRoutingResult.reason}</p>
                  </div>
                )}
              </div>

              {/* Card 3: Mass Global Search & Replace */}
              <div className="bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-lg flex flex-col justify-between hover:border-slate-700 transition-colors">
                <div>
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-base font-semibold text-emerald-400 flex items-center gap-2" title="Executes high-performance text replacements across all 100+ submodules in the Omni-Workspace.">
                      Global Code Refactor
                    </h3>
                    <span className="text-[10px] uppercase font-mono px-2 py-0.5 bg-emerald-950/80 text-emerald-300 border border-emerald-800 rounded">
                      Omni-Workspace
                    </span>
                  </div>
                  <p className="text-xs text-slate-300 mb-3">
                    Multi-repository mass refactoring tool. Test with dry-run before applying destructive edits.
                  </p>

                  <form onSubmit={handleGlobalRefactor} className="space-y-2">
                    <div className="grid grid-cols-2 gap-2">
                      <input
                        type="text"
                        value={searchStr}
                        onChange={e => setSearchStr(e.target.value)}
                        placeholder="Search term..."
                        className="bg-slate-950 border border-slate-700 rounded-lg px-2.5 py-1.5 text-xs text-slate-200 placeholder-slate-500 focus:outline-none focus:border-emerald-500"
                        title="Search string across all workspace projects"
                      />
                      <input
                        type="text"
                        value={replaceStr}
                        onChange={e => setReplaceStr(e.target.value)}
                        placeholder="Replace with..."
                        className="bg-slate-950 border border-slate-700 rounded-lg px-2.5 py-1.5 text-xs text-slate-200 placeholder-slate-500 focus:outline-none focus:border-emerald-500"
                        title="Replacement target string"
                      />
                    </div>
                    <div className="flex items-center justify-between pt-1">
                      <label className="flex items-center space-x-2 text-xs text-slate-300 cursor-pointer" title="Dry run lists files that will change without altering content on disk.">
                        <input
                          type="checkbox"
                          checked={dryRun}
                          onChange={e => setDryRun(e.target.checked)}
                          className="rounded bg-slate-950 border-slate-700 text-emerald-600 focus:ring-0"
                        />
                        <span>Dry Run (Safe Preview)</span>
                      </label>
                      <button
                        type="submit"
                        disabled={refactorLoading}
                        className="bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white font-semibold px-3 py-1.5 rounded-lg text-xs transition-colors shadow"
                      >
                        {refactorLoading ? "Scanning..." : "Execute"}
                      </button>
                    </div>
                  </form>
                </div>

                {refactorResult && (
                  <div className="mt-3 p-3 bg-slate-950 border border-emerald-500/30 rounded-xl text-xs space-y-1 max-h-24 overflow-y-auto">
                    <p className="text-emerald-300 font-bold">Affected: {refactorResult.modified_files} file(s)</p>
                    {refactorResult.files?.map((f, idx) => (
                      <div key={idx} className="font-mono text-[10px] text-slate-400 truncate">{f}</div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          </section>
        )}

        {/* SECTION 2: ECOSYSTEM MONITOR & SUBMODULE MANAGEMENT */}
        {(activeTab === 'all' || activeTab === 'submodules') && (
          <section className="space-y-4">
            <div className="flex flex-wrap items-center justify-between gap-4">
              <div>
                <h2 className="text-xl font-bold text-white flex items-center gap-2">
                  <span className="w-2 h-2 rounded-full bg-emerald-500" />
                  Ecosystem Repositories & Submodule Synchronization
                </h2>
                <p className="text-xs text-slate-400">Status, drift detection, conflict resolution, and automated CI/CD pipeline triggers.</p>
              </div>

              {/* Submodule Action Buttons */}
              <div className="flex flex-wrap gap-2 text-xs">
                <button
                  onClick={handleCheckDrift}
                  className="bg-slate-800 hover:bg-slate-700 border border-slate-700 px-3 py-1.5 rounded-lg text-slate-200 transition-colors font-medium flex items-center gap-1"
                  title="Identifies submodules that have fallen behind their upstream remote repositories."
                >
                  Detect Drift
                </button>
                <button
                  onClick={handleCheckConflicts}
                  className="bg-slate-800 hover:bg-slate-700 border border-slate-700 px-3 py-1.5 rounded-lg text-slate-200 transition-colors font-medium flex items-center gap-1"
                  title="Scans for unmerged git conflicts and applies feature-preserving resolution routines."
                >
                  Resolve Conflicts
                </button>
                <button
                  onClick={handlePruneSubmodules}
                  className="bg-slate-800 hover:bg-slate-700 border border-slate-700 px-3 py-1.5 rounded-lg text-slate-200 transition-colors font-medium flex items-center gap-1"
                  title="Safely cleans up stale or orphaned submodule pointers."
                >
                  Prune Broken
                </button>
                <button
                  onClick={handleGenerateDashboard}
                  className="bg-blue-600 hover:bg-blue-500 text-white px-3 py-1.5 rounded-lg font-medium transition-colors shadow"
                  title="Rebuilds SUBMODULE_DASHBOARD.md across all registered submodules."
                >
                  Regenerate Dashboard
                </button>
              </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Submodule Table (2 cols wide on desktop) */}
              <div className="lg:col-span-2 bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-lg">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-base font-semibold text-slate-200">Active Submodules ({submodules.length})</h3>
                  <button onClick={fetchSubmodules} className="text-xs text-blue-400 hover:underline">Refresh</button>
                </div>

                <div className="overflow-x-auto max-h-72 overflow-y-auto border border-slate-800 rounded-xl">
                  <table className="w-full text-left text-xs font-mono">
                    <thead className="bg-slate-950 text-slate-400 border-b border-slate-800 sticky top-0">
                      <tr>
                        <th className="px-4 py-2.5">Submodule Path</th>
                        <th className="px-4 py-2.5">Upstream State</th>
                        <th className="px-4 py-2.5 text-right">Actions</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-800/60 bg-slate-900/50">
                      {submodules.length === 0 ? (
                        <tr>
                          <td colSpan="3" className="px-4 py-6 text-center text-slate-500 font-sans">
                            No submodules detected or backend endpoint unreachable.
                          </td>
                        </tr>
                      ) : (
                        submodules.map((sm, idx) => (
                          <tr key={idx} className="hover:bg-slate-800/50 transition-colors">
                            <td className="px-4 py-2.5 font-medium text-slate-200" title={`Path: ${sm.path}`}>
                              {sm.path}
                            </td>
                            <td className="px-4 py-2.5">
                              <span className={`px-2 py-0.5 rounded text-[10px] font-bold ${sm.status === 'Synchronized' ? 'bg-emerald-950 text-emerald-400 border border-emerald-800' : 'bg-rose-950 text-rose-400 border border-rose-800'}`}>
                                {sm.status}
                              </span>
                            </td>
                            <td className="px-4 py-2.5 text-right font-sans">
                              <button
                                onClick={() => setPipelineSubmodule(sm.path)}
                                className="text-blue-400 hover:text-blue-300 text-xs font-semibold mr-2"
                                title="Target this submodule for CI/CD actions"
                              >
                                Select
                              </button>
                            </td>
                          </tr>
                        ))
                      )}
                    </tbody>
                  </table>
                </div>

                {/* Submodule Drift / Conflict Results Panel if triggered */}
                {(driftReport.length > 0 || conflicts.length > 0) && (
                  <div className="mt-4 p-3 bg-slate-950 border border-slate-800 rounded-xl text-xs space-y-2">
                    {driftReport.length > 0 && (
                      <div>
                        <div className="font-bold text-amber-400 mb-1">Drift Scan Results:</div>
                        <div className="grid grid-cols-2 gap-2 max-h-24 overflow-y-auto">
                          {driftReport.map((d, i) => (
                            <div key={i} className="text-slate-300 font-mono text-[11px] truncate">
                              • {d.submodule}: <span className={d.drift_detected ? "text-rose-400" : "text-emerald-400"}>{d.drift_detected ? "Diverged" : "Up-to-date"}</span>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                    {conflicts.length > 0 && (
                      <div className="pt-2 border-t border-slate-800">
                        <div className="font-bold text-rose-400 mb-1">Conflicts Handled:</div>
                        {conflicts.map((c, i) => (
                          <div key={i} className="text-slate-300 text-[11px]">• {c.submodule}: {c.action_taken}</div>
                        ))}
                      </div>
                    )}
                  </div>
                )}
              </div>

              {/* Submodule CI/CD Deployment Pipeline Runner */}
              <div className="bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-lg flex flex-col justify-between">
                <div>
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-base font-semibold text-cyan-400" title="Generates GitHub Actions CI/CD workflows and triggers stage verification.">
                      CI/CD Pipeline Dispatcher
                    </h3>
                    <span className="text-[10px] uppercase font-mono px-2 py-0.5 bg-cyan-950 text-cyan-300 border border-cyan-800 rounded">
                      GitHub Actions
                    </span>
                  </div>
                  <p className="text-xs text-slate-300 mb-4">
                    Deploy or auto-generate workflow specs (`.github/workflows/`) for any registered submodule.
                  </p>

                  <div className="space-y-3">
                    <div>
                      <label className="block text-[11px] text-slate-400 font-medium mb-1">Target Submodule</label>
                      <input
                        type="text"
                        value={pipelineSubmodule}
                        onChange={e => setPipelineSubmodule(e.target.value)}
                        placeholder="e.g., aios or bobmani"
                        className="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-2 text-xs text-slate-200 focus:outline-none focus:border-cyan-500 font-mono"
                        title="Submodule folder name to operate on"
                      />
                    </div>
                    <div>
                      <label className="block text-[11px] text-slate-400 font-medium mb-1">Deployment Environment</label>
                      <select
                        value={pipelineEnv}
                        onChange={e => setPipelineEnv(e.target.value)}
                        className="w-full bg-slate-950 border border-slate-700 rounded-lg px-3 py-2 text-xs text-slate-200 focus:outline-none focus:border-cyan-500"
                      >
                        <option value="staging">Staging</option>
                        <option value="production">Production</option>
                        <option value="local">Local Dry-Run</option>
                      </select>
                    </div>
                  </div>
                </div>

                <div className="mt-4 pt-4 border-t border-slate-800 flex gap-2">
                  <button
                    onClick={() => handlePipelineAction(false)}
                    disabled={pipelineLoading}
                    className="flex-1 bg-cyan-600 hover:bg-cyan-500 disabled:opacity-50 text-slate-950 font-bold py-2 rounded-lg text-xs transition-colors shadow"
                    title="Executes build, test, and deploy stages"
                  >
                    Run Pipeline
                  </button>
                  <button
                    onClick={() => handlePipelineAction(true)}
                    disabled={pipelineLoading}
                    className="bg-slate-800 hover:bg-slate-700 border border-slate-700 text-slate-200 px-3 py-2 rounded-lg text-xs font-medium transition-colors"
                    title="Generates a .github/workflows/<submodule>_pipeline.yml workflow file"
                  >
                    Gen Workflow
                  </button>
                </div>

                {pipelineResult && (
                  <div className={`mt-3 p-2.5 rounded-lg text-xs border ${pipelineResult.success ? 'bg-emerald-950/40 border-emerald-600/40 text-emerald-300' : 'bg-rose-950/40 border-rose-600/40 text-rose-300'}`}>
                    <p className="font-semibold">{pipelineResult.message}</p>
                  </div>
                )}
              </div>
            </div>
          </section>
        )}

        {/* SECTION 3: TELEMETRY & SYSTEM OBSERVABILITY */}
        {(activeTab === 'all' || activeTab === 'telemetry') && (
          <section className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-xl font-bold text-white flex items-center gap-2">
                  <span className="w-2 h-2 rounded-full bg-purple-500" />
                  Realtime Telemetry & Observability
                </h2>
                <p className="text-xs text-slate-400">Database-backed asynchronous task queues, system audit status, and live Go server logs.</p>
              </div>

              <div className="flex gap-2">
                <button
                  onClick={handleRunAudit}
                  className="bg-slate-800 hover:bg-slate-700 border border-slate-700 px-3 py-1.5 rounded-lg text-xs text-slate-300 transition-colors"
                  title="Scans all submodules for missing UI representations (UI Auditor)"
                >
                  Run UI Audit
                </button>
                <button
                  onClick={handleStandardizeTelemetry}
                  className="bg-slate-800 hover:bg-slate-700 border border-slate-700 px-3 py-1.5 rounded-lg text-xs text-slate-300 transition-colors"
                  title="Injects standardized telemetry stubs into repositories"
                >
                  Standardize Telemetry
                </button>
                <button
                  onClick={fetchLogs}
                  className="bg-purple-600/20 hover:bg-purple-600 text-purple-300 hover:text-white border border-purple-500/30 px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors"
                >
                  Refresh Logs
                </button>
              </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Task Queue Telemetry */}
              <div className="bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-lg flex flex-col justify-between">
                <div>
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-base font-semibold text-purple-400" title="GORM SQLite database task queue providing async background job execution.">
                      Task Queue Telemetry
                    </h3>
                    <span className="text-[10px] font-mono px-2 py-0.5 bg-purple-950 text-purple-300 border border-purple-800 rounded">
                      SQLite Queue
                    </span>
                  </div>
                  <p className="text-xs text-slate-300 mb-3">Live task lifecycle persisted to SQLite database.</p>

                  <div className="space-y-2 max-h-56 overflow-y-auto pr-1">
                    {queue.length === 0 ? (
                      <p className="text-xs text-slate-500 italic py-4 text-center">No active tasks in the persistent queue.</p>
                    ) : (
                      queue.map((task, idx) => (
                        <div key={idx} className="p-2.5 bg-slate-950 border border-slate-800 rounded-xl flex items-center justify-between text-xs hover:border-slate-700 transition-colors">
                          <div>
                            <span className="font-mono text-[10px] text-slate-500">#{task.ID}</span>
                            <p className="font-medium text-slate-200">{task.Name}</p>
                          </div>
                          <span className={`text-[10px] font-bold px-2 py-0.5 rounded ${task.Status === 'COMPLETED' ? 'bg-emerald-950 text-emerald-400 border border-emerald-800' : 'bg-amber-950 text-amber-400 border border-amber-800'}`}>
                            {task.Status}
                          </span>
                        </div>
                      ))
                    )}
                  </div>
                </div>

                {auditReport.length > 0 && (
                  <div className="mt-4 pt-3 border-t border-slate-800 text-[11px] text-slate-400 max-h-24 overflow-y-auto">
                    <p className="font-bold text-purple-300 mb-1">UI Auditor Highlights:</p>
                    {auditReport.slice(0, 3).map((a, i) => (
                      <div key={i} className="truncate">• {a.submodule}: Missing {a.missing_ui_elements?.length || 0} UI element(s)</div>
                    ))}
                  </div>
                )}
              </div>

              {/* Live Tailer Backend Logs */}
              <div className="lg:col-span-2 bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-lg flex flex-col justify-between">
                <div>
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-base font-semibold text-blue-400" title="Real-time tail of the Go backend execution logs.">
                      System Execution Logs
                    </h3>
                    <span className="text-[10px] font-mono px-2 py-0.5 bg-blue-950 text-blue-300 border border-blue-800 rounded">
                      Live Tail
                    </span>
                  </div>
                  <pre className="bg-slate-950 p-4 rounded-xl text-xs font-mono text-slate-300 h-64 overflow-y-auto border border-slate-800 shadow-inner">
                    {logs || "Waiting for server execution logs..."}
                  </pre>
                </div>
              </div>
            </div>
          </section>
        )}

      </main>

      {/* Global Compact Footer */}
      <footer className="bg-slate-900 border-t border-slate-800 py-3 px-6 text-center text-xs text-slate-500 flex flex-wrap items-center justify-between gap-2">
        <div>
          Jules Autopilot • Autonomous Orchestration Engine
        </div>
        <div className="flex items-center space-x-4">
          <span className="flex items-center gap-1.5">
            <span className="w-2 h-2 rounded-full bg-emerald-500" />
            Vite SPA + Go Fiber/NetHTTP
          </span>
          <span>Port :8080 (API) / :3006 (Dev)</span>
        </div>
      </footer>
    </div>
  )
}

export default App
