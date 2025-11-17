import { useState } from 'react'
import './App.css'
import Dashboard from './components/Dashboard'
import ConfigEditor from './components/ConfigEditor'

function App() {
  const [activeTab, setActiveTab] = useState<'dashboard' | 'config'>('dashboard')

  return (
    <div className="app">
      <header className="app-header">
        <h1>Cloudflared GUI</h1>
        <nav className="tab-navigation">
          <button
            className={activeTab === 'dashboard' ? 'active' : ''}
            onClick={() => setActiveTab('dashboard')}
          >
            Dashboard
          </button>
          <button
            className={activeTab === 'config' ? 'active' : ''}
            onClick={() => setActiveTab('config')}
          >
            Configuration
          </button>
        </nav>
      </header>
      <main className="app-content">
        {activeTab === 'dashboard' ? <Dashboard /> : <ConfigEditor />}
      </main>
    </div>
  )
}

export default App

