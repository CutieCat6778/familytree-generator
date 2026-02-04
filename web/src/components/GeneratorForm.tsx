import React, { useState, useEffect } from 'react';
import { CountryInfo, GenerateRequest } from '../types';
import { getCountries, generateTree, isApiAvailable } from '../utils/api';

interface GeneratorFormProps {
  onGenerate: (data: unknown) => void;
  onError: (error: string) => void;
}

const styles: Record<string, React.CSSProperties> = {
  form: {
    backgroundColor: '#fff',
    borderRadius: '8px',
    padding: '20px',
    boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
    marginBottom: '20px',
  },
  title: {
    fontSize: '18px',
    fontWeight: 'bold',
    marginBottom: '16px',
    color: '#333',
  },
  row: {
    display: 'flex',
    gap: '16px',
    marginBottom: '16px',
    flexWrap: 'wrap',
  },
  field: {
    flex: '1',
    minWidth: '150px',
  },
  label: {
    display: 'block',
    fontSize: '13px',
    fontWeight: '500',
    marginBottom: '4px',
    color: '#555',
  },
  input: {
    width: '100%',
    padding: '8px 12px',
    border: '1px solid #ddd',
    borderRadius: '6px',
    fontSize: '14px',
    boxSizing: 'border-box',
  },
  select: {
    width: '100%',
    padding: '8px 12px',
    border: '1px solid #ddd',
    borderRadius: '6px',
    fontSize: '14px',
    backgroundColor: '#fff',
    boxSizing: 'border-box',
  },
  checkbox: {
    display: 'flex',
    alignItems: 'center',
    gap: '8px',
    fontSize: '14px',
    color: '#555',
  },
  button: {
    padding: '10px 24px',
    backgroundColor: '#4a90d9',
    color: '#fff',
    border: 'none',
    borderRadius: '6px',
    fontSize: '14px',
    fontWeight: '500',
    cursor: 'pointer',
  },
  buttonDisabled: {
    opacity: 0.6,
    cursor: 'not-allowed',
  },
  offline: {
    backgroundColor: '#fee2e2',
    color: '#dc2626',
    padding: '12px 16px',
    borderRadius: '6px',
    fontSize: '14px',
    marginBottom: '16px',
  },
  hint: {
    fontSize: '12px',
    color: '#888',
    marginTop: '4px',
  },
};

export const GeneratorForm: React.FC<GeneratorFormProps> = ({ onGenerate, onError }) => {
  const [countries, setCountries] = useState<CountryInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [apiOnline, setApiOnline] = useState<boolean | null>(null);

  const [country, setCountry] = useState('germany');
  const [generations, setGenerations] = useState(3);
  const [seed, setSeed] = useState('');
  const [startYear, setStartYear] = useState(1970);
  const [gender, setGender] = useState('');
  const [extended, setExtended] = useState(false);
  const [lifeExpectancyMode, setLifeExpectancyMode] = useState<'total' | 'female' | 'male' | 'by_gender'>('total');

  useEffect(() => {
    
    const init = async () => {
      const online = await isApiAvailable();
      setApiOnline(online);

      if (online) {
        try {
          const data = await getCountries();
          setCountries(data.countries);
        } catch (err) {
          console.error('Failed to load countries:', err);
        }
      }
    };

    init();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const request: GenerateRequest = {
        country,
        generations,
        start_year: startYear,
        include_extended: extended,
        life_expectancy_mode: lifeExpectancyMode,
      };

      if (seed) {
        request.seed = parseInt(seed, 10);
      }

      if (gender) {
        request.gender = gender as 'M' | 'F';
      }

      const response = await generateTree(request);

      if (response.success && response.tree) {
        onGenerate(response.tree);
      } else {
        onError(response.message || 'Generation failed');
      }
    } catch (err) {
      onError(err instanceof Error ? err.message : 'Request failed');
    } finally {
      setLoading(false);
    }
  };

  if (apiOnline === false) {
    return (
      <div style={styles.form}>
        <div style={styles.offline}>
          <strong>API Server Offline</strong>
          <br />
          Start the server with: <code>make server</code>
          <br />
          Or load a JSON file instead.
        </div>
      </div>
    );
  }

  if (apiOnline === null) {
    return (
      <div style={styles.form}>
        <div>Connecting to API...</div>
      </div>
    );
  }

  return (
    <form style={styles.form} onSubmit={handleSubmit}>
      <div style={styles.title}>Generate Family Tree</div>

      <div style={styles.row}>
        <div style={styles.field}>
          <label style={styles.label}>Country</label>
          <select
            style={styles.select}
            value={country}
            onChange={e => setCountry(e.target.value)}
          >
            {countries.map(c => (
              <option key={c.slug} value={c.slug}>
                {c.name || c.slug.replace(/-/g, ' ')}
              </option>
            ))}
          </select>
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Generations</label>
          <input
            type="number"
            style={styles.input}
            min={1}
            max={10}
            value={generations}
            onChange={e => setGenerations(parseInt(e.target.value, 10))}
          />
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Start Year</label>
          <input
            type="number"
            style={styles.input}
            min={1800}
            max={2024}
            value={startYear}
            onChange={e => setStartYear(parseInt(e.target.value, 10))}
          />
        </div>
      </div>

      <div style={styles.row}>
        <div style={styles.field}>
          <label style={styles.label}>Seed (optional)</label>
          <input
            type="number"
            style={styles.input}
            placeholder="Random"
            value={seed}
            onChange={e => setSeed(e.target.value)}
          />
          <div style={styles.hint}>Use same seed to reproduce tree</div>
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Root Gender</label>
          <select
            style={styles.select}
            value={gender}
            onChange={e => setGender(e.target.value)}
          >
            <option value="">Random</option>
            <option value="M">Male</option>
            <option value="F">Female</option>
          </select>
        </div>

        <div style={styles.field}>
          <label style={styles.label}>&nbsp;</label>
          <label style={styles.checkbox}>
            <input
              type="checkbox"
              checked={extended}
              onChange={e => setExtended(e.target.checked)}
            />
            Include siblings
          </label>
        </div>

        <div style={styles.field}>
          <label style={styles.label}>Life Expectancy</label>
          <select
            style={styles.select}
            value={lifeExpectancyMode}
            onChange={e => setLifeExpectancyMode(e.target.value as 'total' | 'female' | 'male' | 'by_gender')}
          >
            <option value="total">Total (overall)</option>
            <option value="female">Female</option>
            <option value="male">Male</option>
            <option value="by_gender">By gender</option>
          </select>
          <div style={styles.hint}>Controls which life expectancy baseline to use</div>
        </div>
      </div>

      <button
        type="submit"
        style={{
          ...styles.button,
          ...(loading ? styles.buttonDisabled : {})
        }}
        disabled={loading}
      >
        {loading ? 'Generating...' : 'Generate Tree'}
      </button>
    </form>
  );
};
