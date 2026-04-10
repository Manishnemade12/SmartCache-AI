import { useState } from 'react';
import { HomePage } from './pages/Home';
import { AnalyticsPage } from './pages/Analytics';
import './index.css';

type Page = 'home' | 'analytics';

export default function App() {
  const [page, setPage] = useState<Page>('home');

  return (
    <div className="app-layout">
      {/* Sidebar */}
      <aside className="sidebar">
        <div className="sidebar-logo">
          <div className="sidebar-logo-icon">⚡</div>
          <div>
            <div className="sidebar-logo-text">SmartCache AI</div>
            <div className="sidebar-logo-sub">Async Processing Engine</div>
          </div>
        </div>

        <nav className="sidebar-nav">
          <button
            className={`nav-item ${page === 'home' ? 'active' : ''}`}
            onClick={() => setPage('home')}
          >
            <span className="nav-icon">🏠</span>
            <span>Summarize</span>
          </button>
          <button
            className={`nav-item ${page === 'analytics' ? 'active' : ''}`}
            onClick={() => setPage('analytics')}
          >
            <span className="nav-icon">📊</span>
            <span>Analytics</span>
          </button>
        </nav>

        <div className="sidebar-footer">
          Go · Valkey · Gemini
        </div>
      </aside>

      {/* Main */}
      <main className="main-content">
        {page === 'home' && <HomePage />}
        {page === 'analytics' && (
          <div className="page-content" style={{ paddingTop: 32 }}>
            <AnalyticsPage />
          </div>
        )}
      </main>
    </div>
  );
}
