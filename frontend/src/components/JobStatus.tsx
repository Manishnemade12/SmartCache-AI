import { useEffect, useState } from 'react';
import { getJobStatus, type JobStatus } from '../services/api';

interface JobStatusProps {
  jobId: string;
  initialStatus: string;
  onComplete: (job: JobStatus) => void;
}

export function JobStatusCard({ jobId, initialStatus, onComplete }: JobStatusProps) {
  const [status, setStatus] = useState(initialStatus);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (status === 'completed' || status === 'failed') return;

    const poll = async () => {
      try {
        const job = await getJobStatus(jobId);
        setStatus(job.status);
        if (job.status === 'completed' || job.status === 'failed') {
          onComplete(job);
        }
      } catch (e) {
        setError(e instanceof Error ? e.message : 'Polling error');
      }
    };

    const interval = setInterval(poll, 2000);
    poll(); // immediate first check
    return () => clearInterval(interval);
  }, [jobId, status, onComplete]);

  const statusClass = `status-badge status-${status}`;

  const statusIcon: Record<string, string> = {
    pending: '⏳',
    processing: '⚙️',
    completed: '✅',
    failed: '❌',
  };

  return (
    <div className="job-status-card">
      <div className="job-meta">
        <span className={statusClass}>
          {statusIcon[status] || '•'} {status}
        </span>
        <span className="job-id-text">ID: {jobId.slice(0, 8)}...</span>
      </div>

      {(status === 'pending' || status === 'processing') && (
        <div className="progress-bar-wrap">
          <div className="progress-bar" />
        </div>
      )}

      {status === 'pending' && (
        <p style={{ fontSize: 13, color: 'var(--text-secondary)' }}>
          Job queued — a worker will pick it up shortly
        </p>
      )}

      {status === 'processing' && (
        <p style={{ fontSize: 13, color: 'var(--text-secondary)' }}>
          AI is generating your summary...
        </p>
      )}

      {error && (
        <div className="error-box">⚠️ {error}</div>
      )}
    </div>
  );
}
