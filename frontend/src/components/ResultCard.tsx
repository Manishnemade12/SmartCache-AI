import type { JobStatus } from '../services/api';

interface ResultCardProps {
  job: JobStatus;
  cached?: boolean;
}

export function ResultCard({ job, cached }: ResultCardProps) {
  if (!job.summary) return null;

  return (
    <div className="card result-card">
      <div className="card-title">
        🧠 Summary
        {cached && (
          <span className="status-badge status-completed" style={{ fontSize: 10, marginLeft: 'auto' }}>
            ⚡ Cached
          </span>
        )}
      </div>

      <p className="result-summary">{job.summary}</p>

      {job.tags && job.tags.length > 0 && (
        <div className="result-tags">
          {job.tags.map((tag, i) => (
            <span key={i} className="tag">#{tag}</span>
          ))}
        </div>
      )}

      <div className="divider" />

      <div className="result-meta-row">
        {job.duration_ms !== undefined && (
          <span className="meta-chip">
            ⏱ {job.duration_ms}ms
          </span>
        )}
        {job.completed_at && (
          <span className="meta-chip">
            🕐 {new Date(job.completed_at).toLocaleTimeString()}
          </span>
        )}
        {cached && (
          <span className="meta-chip cached">
            ✅ Served from Valkey cache
          </span>
        )}
      </div>
    </div>
  );
}
