import { useEffect, useState } from 'react'

// Single-page consolidated dashboard. Every backend route is wired here,
// ordered by value: server control first, then submodule health, live
// monitoring, AI routing, pipeline, maintenance tools, and logs last.
// The ⓘ marker carries the tooltip (title) for each control.
const Tip = ({ text }) => <span title={text} style={{ cursor: 'help', color: '#7dd3fc' }}> ⓘ</span>

function Section({ title, desc, tip, children }) {
  return (
    <section style={{ background: '#1e293b', border: '1px solid #334155', borderRadius: 8, padding: 20 }}>
      <h2 style={{ margin: '0 0 4px', fontSize: 20, color: '#e2e8f0' }}>{title}<Tip text={tip} /></h2>
      <p style={{ margin: '0 0 14px', fontSize: 13, color: '#cbd5e1' }}>{desc}</p>
      {children}
    </section>
  )
}

const btn = { background: '#0f766e', color: '#fff', border: 0, borderRadius: 6, padding: '7px 14px', cursor: 'pointer', fontWeight: 600 }
const btnDanger = { ...btn, background: '#b91c1c' }
const input = { width: '100%', boxSizing: 'border-box', background: '#0f172a', border: '1px solid #475569', color: '#e2e8f0', borderRadius: 6, padding: '7px 9px', marginBottom: 8 }
const label = { display: 'block', fontSize: 12, fontWeight: 600, color: '#cbd5e1', marginBottom: 4 }
const out = { background: '#020617', borderRadius: 6, padding: 10, fontSize: 12, overflow: 'auto', maxHeight: 260, whiteSpace: 'pre-wrap' }

export default function App() {
  const [server, setServer] = useState(null)
  const [submodules, setSubmodules] = useState([])
  const [anomaly, setAnomaly] = useState(null)
  const [queue, setQueue] = useState([])
  const [logs, setLogs] = useState('')
  const [apiKey, setApiKey] = useState('')
  const [result, setResult] = useState('')
  const [taskDesc, setTaskDesc] = useState('')
  const [pipeSub, setPipeSub] = useState('borg')
  const [pipeEnv, setPipeEnv] = useState('staging')
  const [search, setSearch] = useState('')
  const [replace, setReplace] = useState('')
  const [dryRun, setDryRun] = useState(true)
  const [targetDir, setTargetDir] = useState('')

  const headers = () => (apiKey ? { 'Content-Type': 'application/json', 'X-API-KEY': apiKey } : { 'Content-Type': 'application/json' })
  const call = async (name, url, opts = {}) => {
    try {
      const r = await fetch(url, opts)
      const t = await r.text()
      setResult(`${name} → HTTP ${r.status}\n${t.slice(0, 4000)}`)
    } catch (e) { setResult(`${name} → FAILED: ${e.message} (is the Go backend running on :8080?)`) }
  }
  const refreshCore = () => {
    fetch('/api/system/server/status').then(r => r.json()).then(setServer).catch(() => setServer(null))
    fetch('/system/status').then(r => r.json()).then(d => setSubmodules(d || [])).catch(() => {})
    fetch('/api/shadow/diff').then(r => r.json()).then(setAnomaly).catch(() => {})
    fetch('/api/system/logs').then(r => r.json()).then(d => setLogs(d.logs || '')).catch(() => {})
  }
  useEffect(() => {
    refreshCore()
    const id = setInterval(() => fetch('/api/queue/telemetry').then(r => r.json()).then(d => setQueue(d || [])).catch(() => {}), 3000)
    return () => clearInterval(id)
  }, [])

  return (
    <div style={{ background: '#0f172a', color: '#f1f5f9', minHeight: '100vh', padding: 24, fontFamily: 'system-ui, sans-serif' }}>
      <h1 style={{ margin: '0 0 4px', fontSize: 26 }}>Bobsaver Command Center</h1>
      <p style={{ color: '#cbd5e1', marginTop: 0 }}>One page for every backend feature. Highest-value controls first; destructive actions need the API key.</p>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit,minmax(340px,1fr))', gap: 16 }}>

        <Section title="1 · Server & System Tray" desc="Start here: is the backend alive, and how do you stop it? The tray icon (scripts/system_tray.py) offers the same Open / Status / Logs / Quit Servers menu." tip="GET /api/system/server/status · POST /api/system/server/shutdown (API key required when JULES_API_KEY is set). Tray: python scripts/system_tray.py">
          <p>{server ? `Online — PID ${server.pid}, uptime ${server.uptime_sec}s, since ${server.started_at}` : 'Offline or unreachable — start: cd geiss/backend-go && go run .'}</p>
          <label style={label}>JULES_API_KEY<Tip text="Required for shutdown, dashboard generation and UI audit when the server was started with JULES_API_KEY set." /></label>
          <input style={input} type="password" value={apiKey} onChange={e => setApiKey(e.target.value)} placeholder="API key for destructive actions" />
          <button style={btn} onClick={refreshCore}>Refresh status</button>{' '}
          <button style={btnDanger} title="Gracefully stops the Go backend (same action as tray Quit Servers & Exit)" onClick={() => { if (confirm('Shut down the backend server?')) call('Shutdown', '/api/system/server/shutdown', { method: 'POST', headers: headers(), body: '{}' }) }}>Quit servers</button>
        </Section>

        <Section title="2 · Submodule Health" desc="Sync state of JWildfire, MilkDrop3, projectm, electricsheep, geiss and apophysis-j, plus drift, prune and coverage scans." tip="GET /system/status · POST /api/system/dashboard · GET /api/system/drift · GET+POST /api/system/prune · GET /api/system/audit">
          {submodules.length === 0 ? <p>No data — backend offline or submodules not initialised.</p> : <ul>{submodules.map((sm, i) => <li key={i}><code title={`Commit ${sm.commit}`}>{sm.path}</code> — {sm.status}</li>)}</ul>}
          <button style={btn} onClick={() => call('Dashboard', '/api/system/dashboard', { method: 'POST', headers: headers(), body: '{}' })}>Generate dashboard</button>{' '}
          <button style={btn} onClick={() => call('Drift', '/api/system/drift')}>Check drift</button>{' '}
          <button style={btn} onClick={() => call('Prune scan', '/api/system/prune')}>Prune scan</button>{' '}
          <button style={btn} onClick={() => call('UI audit', '/api/system/audit', { headers: headers() })}>UI audit</button>
        </Section>

        <Section title="3 · Shadow Pilot & Queue" desc="Watches git diffs for runaway automated changes and shows the live in-memory task queue (polled every 3s)." tip="GET /api/shadow/diff · POST /api/shadow/autofix · GET /api/queue/telemetry">
          {anomaly ? <p>Last scan {new Date(anomaly.timestamp).toLocaleString()} · diff {anomaly.diff_size_bytes} bytes{anomaly.warning ? ` · WARNING: ${anomaly.warning}` : ''}</p> : <p>Loading…</p>}
          {anomaly?.warning && <button style={btnDanger} onClick={() => call('Auto-fix', '/api/shadow/autofix', { method: 'POST', headers: headers() })}>Trigger auto-fix</button>}
          <ul>{queue.map((t, i) => <li key={i}><code>{t.ID}</code> {t.Name} — {t.Status}</li>)}</ul>
        </Section>

        <Section title="4 · AI Task Routing" desc="Routes a task description to Claude (architecture/docs), Gemini (scanning/speed) or GPT (code/tests). Use it to coordinate work across models, alongside Edge CDP 9222 sessions." tip="POST /api/tasks/route with {task_description}. Keywords decide: doc/plan/refactor→Claude, scan/fast→Gemini, code/test/implement→GPT.">
          <label style={label}>Task description<Tip text="Plain English. Include words like 'refactor', 'scan all repos' or 'implement and test' to steer routing." /></label>
          <input style={input} value={taskDesc} onChange={e => setTaskDesc(e.target.value)} placeholder="e.g. refactor the dashboard and write tests" />
          <button style={btn} onClick={() => call('Route task', '/api/tasks/route', { method: 'POST', headers: headers(), body: JSON.stringify({ task_description: taskDesc }) })}>Route task</button>
        </Section>

        <Section title="5 · CI/CD Pipeline" desc="Runs build/test/deploy stages for a submodule, or generates its GitHub Actions workflow file." tip="POST /api/system/pipeline with {submodule, environment, generate}. generate=true writes .github/workflows/<name>_pipeline.yml instead of running stages.">
          <label style={label}>Submodule<Tip text="Nested submodule folder inside geiss/ (the backend resolves paths relative to geiss/): aios, bobmani, fwber or borg. Initialise it first with: git -C geiss submodule update --init <name>." /></label>
          <input style={input} value={pipeSub} onChange={e => setPipeSub(e.target.value)} />
          <label style={label}>Environment<Tip text="Deployment target label passed to the deploy stage. Defaults to staging." /></label>
          <input style={input} value={pipeEnv} onChange={e => setPipeEnv(e.target.value)} />
          <button style={btn} onClick={() => call('Pipeline', '/api/system/pipeline', { method: 'POST', headers: headers(), body: JSON.stringify({ submodule: pipeSub, environment: pipeEnv, generate: false }) })}>Run pipeline</button>{' '}
          <button style={btn} onClick={() => call('Generate workflow', '/api/system/pipeline', { method: 'POST', headers: headers(), body: JSON.stringify({ submodule: pipeSub, environment: pipeEnv, generate: true }) })}>Generate workflow</button>
        </Section>

        <Section title="6 · Refactor & Telemetry" desc="Global search-and-replace across submodules (dry-run first) and injection of standard [OMNI] logging stubs into Python/JS files." tip="POST /api/system/refactor {search_str, replace_str, dry_run} · POST /api/system/telemetry {target_dir}. Telemetry writes files — review the report.">
          <label style={label}>Search / Replace<Tip text="Literal text match across all submodules, skipping .git, node_modules and dist. Dry-run lists files without writing." /></label>
          <input style={input} value={search} onChange={e => setSearch(e.target.value)} placeholder="search text" />
          <input style={input} value={replace} onChange={e => setReplace(e.target.value)} placeholder="replace with" />
          <label style={label}><input type="checkbox" checked={dryRun} onChange={e => setDryRun(e.target.checked)} /> Dry run<Tip text="On: only report matching files. Off: files are rewritten in place." /></label>
          <button style={btn} onClick={() => call('Refactor', '/api/system/refactor', { method: 'POST', headers: headers(), body: JSON.stringify({ search_str: search, replace_str: replace, dry_run: dryRun }) })}>Run refactor</button>
          <label style={label}>Telemetry target dir<Tip text="Submodule-relative folder to scan, e.g. geiss. Empty scans the whole workspace root." /></label>
          <input style={input} value={targetDir} onChange={e => setTargetDir(e.target.value)} placeholder="e.g. geiss (optional)" />
          <button style={btn} onClick={() => call('Telemetry', '/api/system/telemetry', { method: 'POST', headers: headers(), body: JSON.stringify({ target_dir: targetDir }) })}>Standardise telemetry</button>
        </Section>

        <Section title="7 · Merge Conflicts" desc="Scans the workspace for unresolved git conflict markers and summarises each block." tip="GET /api/conflicts/resolve — read-only scan; resolve conflicts in git, then re-run to confirm clean.">
          <button style={btn} onClick={() => call('Conflicts', '/api/conflicts/resolve')}>Scan for conflicts</button>
        </Section>

        <Section title="8 · Action Result" desc="Raw response of the last button you pressed, so nothing happens silently." tip="Every control on this page writes its HTTP status and response body here.">
          <pre style={out}>{result || 'No action yet — press any button above.'}</pre>
        </Section>
      </div>

      <section style={{ background: '#1e293b', border: '1px solid #334155', borderRadius: 8, padding: 20, marginTop: 16 }}>
        <h2 style={{ margin: '0 0 4px', fontSize: 20 }}>Backend Logs<Tip text="GET /api/system/logs — tail of backend log output. Refresh from section 1." /></h2>
        <p style={{ fontSize: 13, color: '#cbd5e1' }}>If this is empty while the server shows Online, the backend has produced no log output yet.</p>
        <pre style={out}>{logs || 'Loading logs…'}</pre>
      </section>
    </div>
  )
}
