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
          <h2 className="text-2xl font-semibold mb-4 text-green-400">Submodule Status</h2>
          {submodules.length === 0 ? <p className="text-gray-400">No submodules synced or error loading.</p> : (
            <ul className="space-y-2">
              {submodules.map((sm, idx) => (
                <li key={idx} className="flex justify-between border-b border-gray-700 pb-2">
                  <span className="font-mono text-sm">{sm.path}</span>
                  <span className={`text-sm font-bold ${sm.status === 'Synchronized' ? 'text-green-500' : 'text-red-500'}`}>{sm.status}</span>
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* Shadow Pilot */}
        <div className="bg-gray-800 p-6 rounded-lg shadow-lg">
          <h2 className="text-2xl font-semibold mb-4 text-yellow-400">Shadow Pilot</h2>
          {anomaly ? (
             <div>
                <p><strong>Last Scan:</strong> {new Date(anomaly.timestamp).toLocaleString()}</p>
                <p><strong>Diff Size:</strong> {anomaly.diff_size_bytes} bytes</p>
                {anomaly.warning && <p className="text-red-400 font-bold mt-2">WARNING: {anomaly.warning}</p>}
             </div>
          ) : (
             <p className="text-gray-400">Loading shadow pilot data...</p>
          )}
        </div>

        {/* Queue Telemetry */}
        <div className="bg-gray-800 p-6 rounded-lg shadow-lg">
          <h2 className="text-2xl font-semibold mb-4 text-purple-400">Queue Telemetry</h2>
          {queue.length === 0 ? <p className="text-gray-400">No tasks in queue.</p> : (
            <ul className="space-y-2">
              {queue.map((task, idx) => (
                <li key={idx} className="flex justify-between items-center border-b border-gray-700 pb-2">
                  <div>
                     <p className="font-mono text-sm text-gray-300">{task.ID}</p>
                     <p className="font-semibold text-md">{task.Name}</p>
                  </div>
                  <span className={`text-sm font-bold px-2 py-1 rounded ${task.Status === 'COMPLETED' ? 'bg-green-600 text-white' : 'bg-yellow-600 text-white'}`}>{task.Status}</span>
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
