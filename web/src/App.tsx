import React, { useState } from 'react';
import { TreeView } from './components/TreeView';
import { StatsPanel } from './components/StatsPanel';
import { PersonDetail } from './components/PersonDetail';
import { FileLoader } from './components/FileLoader';
import { GeneratorForm } from './components/GeneratorForm';
import { useTreeData } from './hooks/useTreeData';
import { VisualizationNode, VisualizationData } from './types';

const styles: Record<string, React.CSSProperties> = {
  app: {
    minHeight: '100vh',
    display: 'flex',
    flexDirection: 'column',
  },
  header: {
    backgroundColor: '#fff',
    padding: '16px 24px',
    borderBottom: '1px solid #eee',
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  title: {
    fontSize: '24px',
    fontWeight: 'bold',
    color: '#333',
  },
  subtitle: {
    fontSize: '14px',
    color: '#666',
    marginTop: '4px',
  },
  main: {
    flex: 1,
    display: 'flex',
    padding: '20px',
    gap: '20px',
    overflow: 'hidden',
  },
  content: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    minWidth: 0,
  },
  treeContainer: {
    flex: 1,
    backgroundColor: '#fff',
    borderRadius: '12px',
    boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
    overflowX: 'auto',
    overflowY: 'hidden',
    minWidth: 0,
  },
  sidebar: {
    width: '280px',
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',
  },
  welcome: {
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: '#fff',
    borderRadius: '12px',
    padding: '40px',
    textAlign: 'center',
  },
  welcomeTitle: {
    fontSize: '28px',
    fontWeight: 'bold',
    color: '#333',
    marginBottom: '16px',
  },
  welcomeText: {
    fontSize: '16px',
    color: '#666',
    marginBottom: '32px',
    maxWidth: '500px',
    lineHeight: 1.6,
  },
  error: {
    backgroundColor: '#fee2e2',
    color: '#dc2626',
    padding: '12px 16px',
    borderRadius: '8px',
    marginTop: '16px',
  },
  legend: {
    backgroundColor: '#fff',
    borderRadius: '8px',
    padding: '16px',
    boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
  },
  legendTitle: {
    fontSize: '14px',
    fontWeight: 'bold',
    marginBottom: '12px',
    color: '#333',
  },
  legendItem: {
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    marginBottom: '8px',
    fontSize: '13px',
    color: '#666',
  },
  legendColor: {
    width: '16px',
    height: '16px',
    borderRadius: '4px',
  },
  instructions: {
    backgroundColor: '#fff',
    borderRadius: '8px',
    padding: '16px',
    boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
    fontSize: '13px',
    color: '#666',
    lineHeight: 1.6,
  },
  tabs: {
    display: 'flex',
    gap: '8px',
    marginBottom: '16px',
  },
  tab: {
    padding: '8px 16px',
    backgroundColor: '#e5e7eb',
    border: 'none',
    borderRadius: '6px',
    cursor: 'pointer',
    fontSize: '14px',
    color: '#666',
  },
  tabActive: {
    backgroundColor: '#4a90d9',
    color: '#fff',
  },
};

export const App: React.FC = () => {
  const { data, loading, error, loadFromFile, loadSample, setData, setError } = useTreeData();
  const [selectedPerson, setSelectedPerson] = useState<VisualizationNode | null>(null);
  const [mode, setMode] = useState<'file' | 'api'>('api');
  const referenceYear = data
    ? data.reference_year ?? data.nodes.reduce((max, node) => Math.max(max, node.birth_year), 0)
    : new Date().getFullYear();

  const handlePersonClick = (person: VisualizationNode) => {
    setSelectedPerson(person);
  };

  const handleApiGenerate = (tree: unknown) => {
    setData(tree as VisualizationData);
  };

  const handleApiError = (err: string) => {
    setError(err);
  };

  return (
    <div style={styles.app}>
      <header style={styles.header}>
        <div>
          <div style={styles.title}>Family Tree Generator</div>
          <div style={styles.subtitle}>Generate realistic family trees based on demographic data</div>
        </div>
        <div style={styles.tabs}>
          <button
            style={{ ...styles.tab, ...(mode === 'api' ? styles.tabActive : {}) }}
            onClick={() => setMode('api')}
          >
            Generate
          </button>
          <button
            style={{ ...styles.tab, ...(mode === 'file' ? styles.tabActive : {}) }}
            onClick={() => setMode('file')}
          >
            Load File
          </button>
        </div>
      </header>

      <main style={styles.main}>
        {data ? (
          <>
            <div style={styles.content}>
              {mode === 'api' && (
                <GeneratorForm onGenerate={handleApiGenerate} onError={handleApiError} />
              )}
              {mode === 'file' && (
                <div style={{ marginBottom: '16px' }}>
                  <FileLoader
                    onFileLoad={loadFromFile}
                    onSampleLoad={loadSample}
                    loading={loading}
                  />
                </div>
              )}
              <div style={styles.treeContainer}>
                <TreeView
                  data={data}
                  onPersonClick={handlePersonClick}
                  selectedPersonId={selectedPerson?.id || null}
                />
              </div>
            </div>
            <div style={styles.sidebar}>
              <StatsPanel
                stats={data.stats}
                country={data.country}
                seed={data.seed}
                referenceYear={referenceYear}
              />
              <div style={styles.legend}>
                <div style={styles.legendTitle}>Legend</div>
                <div style={styles.legendItem}>
                  <div style={{ ...styles.legendColor, backgroundColor: '#4a90d9' }} />
                  <span>Male (living)</span>
                </div>
                <div style={styles.legendItem}>
                  <div style={{ ...styles.legendColor, backgroundColor: '#d94a7b' }} />
                  <span>Female (living)</span>
                </div>
                <div style={styles.legendItem}>
                  <div style={{ ...styles.legendColor, backgroundColor: '#7a9cc6' }} />
                  <span>Male (deceased)</span>
                </div>
                <div style={styles.legendItem}>
                  <div style={{ ...styles.legendColor, backgroundColor: '#c6a7b8' }} />
                  <span>Female (deceased)</span>
                </div>
                <div style={styles.legendItem}>
                  <div style={{ ...styles.legendColor, border: '2px solid gold', backgroundColor: 'transparent' }} />
                  <span>Root person</span>
                </div>
              </div>
              <div style={styles.instructions}>
                <strong>Controls:</strong><br />
                - Click person to view details<br />
                - Scroll to zoom in/out<br />
                - Drag to pan the view<br />
                - Dashed lines = spouse links
              </div>
            </div>
          </>
        ) : (
          <div style={styles.content}>
            {mode === 'api' ? (
              <GeneratorForm onGenerate={handleApiGenerate} onError={handleApiError} />
            ) : (
              <div style={styles.welcome}>
                <div style={styles.welcomeTitle}>Load Family Tree</div>
                <div style={styles.welcomeText}>
                  Load a family tree JSON file generated by the CLI tool, or click "Load Sample Data" to see a demo.
                  <br /><br />
                  Generate with CLI:<br />
                  <code style={{ backgroundColor: '#f0f0f0', padding: '4px 8px', borderRadius: '4px' }}>
                    ./familytree -format json -output tree.json
                  </code>
                  <br /><br />
                  Then load the <code>tree_viz.json</code> file.
                </div>
                <FileLoader
                  onFileLoad={loadFromFile}
                  onSampleLoad={loadSample}
                  loading={loading}
                />
                {error && <div style={styles.error}>{error}</div>}
              </div>
            )}
          </div>
        )}
      </main>

      {selectedPerson && (
        <PersonDetail
          person={selectedPerson}
          referenceYear={referenceYear}
          onClose={() => setSelectedPerson(null)}
        />
      )}
    </div>
  );
};

export default App;
