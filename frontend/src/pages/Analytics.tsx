import { useCallback, useEffect, useState } from 'react';
import { getAnalytics, type Analytics } from '../services/api';

export function AnalyticsPage() {
  const [metrics, setMetrics] = useState<Analytics | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchMetrics = useCallback(async () => {
    try {
      const data = await getAnalytics();
      setMetrics(data);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load analytics');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchMetrics();
    const interval = setInterval(fetchMetrics, 10000);
    return () => clearInterval(interval);
  }, [fetchMetrics]);

  const hitRate =
    metrics && metrics.total_requests > 0
      ? Math.round((metrics.cache_hits / metrics.total_requests) * 100)
      : 0;

  return (
    <div>
      <div className="analytics-header">
        <div>
          <h1 className="page-title">📊 Analytics</h1>
          <p className="page-subtitle">Live system metrics — auto-refreshes every 10s</p>
        </div>
        <button className="refresh-btn" onClick={fetchMetrics}>
          🔄 Refresh
        </button>
      </div>

      {error && <div className="error-box">⚠️ {error}</div>}

      {loading && !metrics && (
        <div className="empty-state">
          <div className="spinner" style={{ margin: '0 auto' }} />
          <p style={{ marginTop: 12 }}>Loading metrics...</p>
        </div>
      )}

      {metrics && (
        <>
          <div className="metrics-grid">
            <div className="metric-card">
              <div className="metric-label"><span className="metric-icon">📥</span> Total Requests</div>
              <div className="metric-value">{metrics.total_requests.toLocaleString()}</div>
              <div className="metric-sub">All-time submissions</div>
            </div>

            <div className="metric-card">
              <div className="metric-label"><span className="metric-icon">⚡</span> Cache Hits</div>
              <div className="metric-value" style={{ color: 'var(--success)' }}>
                {metrics.cache_hits.toLocaleString()}
              </div>
              <div className="metric-sub">Served from Valkey</div>
            </div>

            <div className="metric-card">
              <div className="metric-label"><span className="metric-icon">🔄</span> Cache Misses</div>
              <div className="metric-value" style={{ color: 'var(--warning)' }}>
                {metrics.cache_misses.toLocaleString()}
              </div>
              <div className="metric-sub">Sent to worker queue</div>
            </div>

            <div className="metric-card">
              <div className="metric-label"><span className="metric-icon">📋</span> Queue Size</div>
              <div className="metric-value" style={{ color: metrics.queue_size > 5 ? 'var(--danger)' : 'var(--text-primary)' }}>
                {metrics.queue_size}
              </div>
              <div className="metric-sub">Jobs waiting</div>
            </div>

            <div className="metric-card">
              <div className="metric-label"><span className="metric-icon">⏱</span> Avg Processing</div>
              <div className="metric-value">{Math.round(metrics.avg_processing_time_ms)}</div>
              <div className="metric-sub">Milliseconds per job</div>
            </div>

            <div className="metric-card">
              <div className="metric-label"><span className="metric-icon">❌</span> Failed Jobs</div>
              <div className="metric-value" style={{ color: metrics.failed_jobs > 0 ? 'var(--danger)' : 'var(--success)' }}>
                {metrics.failed_jobs}
              </div>
              <div className="metric-sub">Processing errors</div>
            </div>
          </div>

          <div className="card">
            <div className="card-title">📈 Cache Hit Rate</div>
            <div className="hit-rate-section">
              <div className="hit-rate-label">
                <span>Cache Efficiency</span>
                <strong style={{ color: hitRate > 60 ? 'var(--success)' : 'var(--warning)' }}>
                  {hitRate}%
                </strong>
              </div>
              <div className="hit-rate-bar">
                <div className="hit-rate-fill" style={{ width: `${hitRate}%` }} />
              </div>
              <div style={{ marginTop: 8, fontSize: 12, color: 'var(--text-muted)' }}>
                {hitRate > 60
                  ? '✅ Excellent — Valkey is saving significant AI API calls'
                  : hitRate > 30
                  ? '⚠️ Moderate — Cache warming up as more requests come in'
                  : 'ℹ️ Low — most requests are new and being processed by workers'}
              </div>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
