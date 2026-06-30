import { useEffect, useState } from 'react'

function App() {
  const [submodules, setSubmodules] = useState([])
  const [anomaly, setAnomaly] = useState(null)
  const [queue, setQueue] = useState([])

  useEffect(() => {
    fetch('/system/status')
      .then(res => res.json())
      .then(data => setSubmodules(data || []))
      .catch(err => console.error("Failed to fetch submodules:", err))

    fetch('/api/shadow/diff')
      .then(res => res.json())
      .then(data => setAnomaly(data))
      .catch(err => console.error("Failed to fetch diff logic:", err))

    // Polling queue telemetry
    const interval = setInterval(() => {
      fetch('/api/queue/telemetry')
        .then(res => res.json())
        .then(data => setQueue(data || []))
        .catch(err => console.error("Failed to fetch queue telemetry:", err))
    }, 2000)

    return () => clearInterval(interval)
  }, [])

  return (
    <div className="p-8 bg-gray-900 text-white min-h-screen">
      <h1 className="text-3xl font-bold mb-6 text-blue-400">Jules Autopilot Dashboard</h1>

      <div className="grid grid-cols-3 gap-8">
        {/* Submodule Status */}
        <div className="bg-gray-800 p-6 rounded-lg shadow-lg">
          <h2 className="text-2xl font-semibold mb-4 text-green-400" title="Live status of all nested git submodules within the Omni-Workspace.">Submodule Status</h2>
          <p className="text-sm text-gray-300 mb-4">Displays the synchronization state of injected projects.</p>
          {submodules.length === 0 ? <p className="text-gray-400">No submodules synced or error loading.</p> : (
            <ul className="space-y-2">
              {submodules.map((sm, idx) => (
                <li key={idx} className="flex justify-between items-center border-b border-gray-700 pb-2 hover:bg-gray-750 transition-colors">
                  <span className="font-mono text-sm font-medium">{sm.path}</span>
                  <span className={`text-xs font-bold px-2 py-1 rounded ${sm.status === 'Synchronized' ? 'bg-green-900 text-green-300' : 'bg-red-900 text-red-300'}`} title={`Current state: ${sm.status}`}>
                    {sm.status}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* Shadow Pilot */}
        <div className="bg-gray-800 p-6 rounded-lg shadow-lg">
          <h2 className="text-2xl font-semibold mb-4 text-yellow-400" title="Actively monitors the repository for massive structural changes or anomalies.">Shadow Pilot</h2>
          {anomaly ? (
             <div>
                <p className="text-sm text-gray-300 mb-2">Monitors git diffs for automated refactoring anomalies.</p>
                <p><strong>Last Scan:</strong> {new Date(anomaly.timestamp).toLocaleString()}</p>
                <p><strong>Diff Size:</strong> {anomaly.diff_size_bytes} bytes</p>
                {anomaly.warning && (
                  <div className="mt-4 p-3 bg-red-900 bg-opacity-50 border border-red-500 rounded">
                    <p className="text-red-400 font-bold">WARNING: {anomaly.warning}</p>
                    <button
                      onClick={() => {
                        fetch('/api/shadow/autofix', { method: 'POST', headers: { 'X-API-KEY': '' } })
                          .then(res => res.json())
                          .then(data => alert(data.message))
                          .catch(err => alert("Auto-fix trigger failed."));
                      }}
                      className="mt-2 bg-red-600 hover:bg-red-700 text-white font-bold py-1 px-3 rounded text-sm transition-colors"
                      title="Trigger the CI Pipeline Auto-Fix routine via the Go backend."
                    >
                      Trigger Auto-Fix
                    </button>
                  </div>
                )}
             </div>
          ) : (
             <p className="text-gray-400">Loading shadow pilot data...</p>
          )}
        </div>

        {/* Queue Telemetry */}
        <div className="bg-gray-800 p-6 rounded-lg shadow-lg">
          <h2 className="text-2xl font-semibold mb-4 text-purple-400" title="Real-time view of the database-backed task execution queue.">Queue Telemetry</h2>
          <p className="text-sm text-gray-300 mb-4">Live execution stream from the Go SQLite Task Queue.</p>
          {queue.length === 0 ? <p className="text-gray-400">No tasks in queue.</p> : (
            <ul className="space-y-2 max-h-96 overflow-y-auto pr-2 custom-scrollbar">
              {queue.map((task, idx) => (
                <li key={idx} className="flex justify-between items-center border-b border-gray-700 pb-3 mb-2 hover:bg-gray-750 transition-colors">
                  <div>
                     <p className="font-mono text-xs text-gray-400 mb-1" title="Unique Task ID">{task.ID}</p>
                     <p className="font-semibold text-sm text-gray-100">{task.Name}</p>
                  </div>
                  <span className={`text-xs font-bold px-2 py-1 rounded shadow-sm ${task.Status === 'COMPLETED' ? 'bg-green-600 text-white' : 'bg-yellow-500 text-gray-900'}`} title={`Task status: ${task.Status}`}>
                    {task.Status}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </div>

      </div>
    </div>
  )
}

export default App
