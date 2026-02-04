import React, { useRef } from 'react';

interface FileLoaderProps {
  onFileLoad: (file: File) => void;
  onSampleLoad: () => void;
  loading: boolean;
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: 'flex',
    gap: '12px',
    alignItems: 'center',
  },
  button: {
    padding: '10px 20px',
    backgroundColor: '#4a90d9',
    color: '#fff',
    border: 'none',
    borderRadius: '6px',
    fontSize: '14px',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
  },
  secondaryButton: {
    padding: '10px 20px',
    backgroundColor: '#fff',
    color: '#4a90d9',
    border: '2px solid #4a90d9',
    borderRadius: '6px',
    fontSize: '14px',
    cursor: 'pointer',
  },
  hiddenInput: {
    display: 'none',
  },
  loading: {
    opacity: 0.6,
    cursor: 'not-allowed',
  }
};

export const FileLoader: React.FC<FileLoaderProps> = ({ onFileLoad, onSampleLoad, loading }) => {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      onFileLoad(file);
    }
  };

  return (
    <div style={styles.container}>
      <input
        ref={fileInputRef}
        type="file"
        accept=".json"
        onChange={handleFileChange}
        style={styles.hiddenInput}
      />
      <button
        style={{ ...styles.button, ...(loading ? styles.loading : {}) }}
        onClick={handleClick}
        disabled={loading}
      >
        {loading ? 'Loading...' : 'Load JSON File'}
      </button>
      <button
        style={{ ...styles.secondaryButton, ...(loading ? styles.loading : {}) }}
        onClick={onSampleLoad}
        disabled={loading}
      >
        Load Sample Data
      </button>
    </div>
  );
};
