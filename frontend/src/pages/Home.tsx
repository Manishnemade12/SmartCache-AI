import { useState, useCallback } from 'react';
import { SubmitForm } from '../components/SubmitForm';
import { JobStatusCard } from '../components/JobStatus';
import { ResultCard } from '../components/ResultCard';
import { submitJob, type SubmitResponse, type JobStatus } from '../services/api';

interface ActiveJob {
  id: string;
  status: string;
  cached: boolean;
  result?: JobStatus;
}

export function HomePage() {
  const [loading, setLoading] = useState(false);
  const [activeJob, setActiveJob] = useState<ActiveJob | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (input: string) => {
    setLoading(true);
    setError(null);
    setActiveJob(null);

    try {
      const resp: SubmitResponse = await submitJob(input);

      if (resp.status === 'completed' && resp.cached) {
        // Instant cache hit
        setActiveJob({
          id: resp.job_id,
          status: 'completed',
          cached: true,
          result: {
            job_id: resp.job_id,
            status: 'completed',
            summary: resp.summary,
            tags: resp.tags,
          },
        });
      } else {
        setActiveJob({ id: resp.job_id, status: resp.status, cached: false });
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Something went wrong');
    } finally {
      setLoading(false);
    }
  };

  const handleJobComplete = useCallback((job: JobStatus) => {
    setActiveJob(prev =>
      prev ? { ...prev, status: job.status, result: job } : null
    );
  }, []);

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">⚡ Summarize</h1>
        <p className="page-subtitle">
          Submit text or a URL — workers process it asynchronously, results cached in Valkey
        </p>
      </div>

      <div className="page-content">
        {error && (
          <div className="error-box" style={{ marginBottom: 24 }}>
            ⚠️ {error}
          </div>
        )}

        <div className="home-grid">
          {/* Left: Input */}
          <div className="card">
            <div className="card-title">📝 Input</div>
            <SubmitForm onSubmit={handleSubmit} loading={loading} />
          </div>

          {/* Right: Status / Result */}
          <div>
            {!activeJob && (
              <div className="card">
                <div className="empty-state">
                  <div className="empty-icon">🤖</div>
                  <div className="empty-title">Ready to summarize</div>
                  <div className="empty-desc">
                    Submit text or a URL to get an AI-powered summary
                  </div>
                </div>
              </div>
            )}

            {activeJob && activeJob.status !== 'completed' && activeJob.status !== 'failed' && (
              <div className="card" style={{ marginBottom: 16 }}>
                <div className="card-title">🔄 Job Status</div>
                <JobStatusCard
                  jobId={activeJob.id}
                  initialStatus={activeJob.status}
                  onComplete={handleJobComplete}
                />
              </div>
            )}

            {activeJob?.result && (
              <ResultCard job={activeJob.result} cached={activeJob.cached} />
            )}

            {activeJob?.status === 'failed' && (
              <div className="card">
                <div className="error-box">
                  ❌ Job failed. {activeJob.result?.error || 'Unknown error occurred.'}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
