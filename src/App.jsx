import { useEffect, useState } from 'react'

function App() {
  const [submodules, setSubmodules] = useState([])
  const [anomaly, setAnomaly] = useState(null)

  useEffect(() => {
    fetch('/system/status')
      .then(res => res.json())
      .then(data => setSubmodules(data || []))
      .catch(err => console.error("Failed to fetch submodules:", err))

    fetch('/api/shadow/diff')
      .then(res => res.json())
      .then(data => setAnomaly(data))
      .catch(err => console.error("Failed to fetch diff logic:", err))
  }, [])

  return (
    <div className="p-8 bg-gray-900 text-white min-h-screen">
      <h1 className="text-3xl font-bold mb-6 text-blue-400">Jules Autopilot Dashboard</h1>

      <div className="grid grid-cols-2 gap-8">
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

        <div className="bg-gray-800 p-6 rounded-lg shadow-lg">
          <h2 className="text-2xl font-semibold mb-4 text-yellow-400">Shadow Pilot Anomaly Detection</h2>
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
      </div>
    </div>
  )
}

export default App
