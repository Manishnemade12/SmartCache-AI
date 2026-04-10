import { useState } from 'react';

interface SubmitFormProps {
  onSubmit: (input: string, isUrl: boolean) => void;
  loading: boolean;
}

export function SubmitForm({ onSubmit, loading }: SubmitFormProps) {
  const [input, setInput] = useState('');
  const [mode, setMode] = useState<'text' | 'url'>('text');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim()) return;
    onSubmit(input.trim(), mode === 'url');
  };

  const placeholder =
    mode === 'text'
      ? 'Paste your article, blog post, or any text you want summarized...'
      : 'https://example.com/article-to-summarize';

  return (
    <form className="submit-form" onSubmit={handleSubmit}>
      <div className="input-group">
        <div className="input-toggle">
          <button
            type="button"
            className={`toggle-btn ${mode === 'text' ? 'active' : ''}`}
            onClick={() => { setMode('text'); setInput(''); }}
          >
            📝 Text
          </button>
          <button
            type="button"
            className={`toggle-btn ${mode === 'url' ? 'active' : ''}`}
            onClick={() => { setMode('url'); setInput(''); }}
          >
            🔗 URL
          </button>
        </div>

        <label className="input-label">
          {mode === 'text' ? 'Enter text to summarize' : 'Enter URL to fetch & summarize'}
        </label>

        <textarea
          className={`text-input ${mode === 'url' ? 'url-input' : ''}`}
          placeholder={placeholder}
          value={input}
          onChange={e => setInput(e.target.value)}
          disabled={loading}
          rows={mode === 'url' ? 1 : 6}
        />

        {mode === 'text' && (
          <div className="char-count">{input.length} characters</div>
        )}
      </div>

      <button
        type="submit"
        className="submit-btn"
        disabled={loading || !input.trim()}
      >
        {loading ? (
          <>
            <div className="spinner" />
            Submitting...
          </>
        ) : (
          <>⚡ Summarize</>
        )}
      </button>
    </form>
  );
}
