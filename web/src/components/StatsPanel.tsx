import React from 'react';
import { VisualizationStats } from '../types';

interface StatsPanelProps {
  stats: VisualizationStats;
  country: string;
  seed: number;
  referenceYear?: number;
}

const styles: Record<string, React.CSSProperties> = {
  panel: {
    backgroundColor: '#fff',
    borderRadius: '8px',
    padding: '16px',
    boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
    minWidth: '220px',
  },
  title: {
    fontSize: '16px',
    fontWeight: 'bold',
    marginBottom: '12px',
    color: '#333',
    borderBottom: '1px solid #eee',
    paddingBottom: '8px',
  },
  stat: {
    display: 'flex',
    justifyContent: 'space-between',
    marginBottom: '8px',
    fontSize: '13px',
  },
  label: {
    color: '#666',
  },
  value: {
    fontWeight: 'bold',
    color: '#333',
  },
  section: {
    marginTop: '12px',
    paddingTop: '12px',
    borderTop: '1px solid #eee',
  },
  sectionTitle: {
    fontSize: '12px',
    fontWeight: 'bold',
    color: '#888',
    marginBottom: '8px',
    textTransform: 'uppercase',
  },
  bar: {
    height: '8px',
    borderRadius: '4px',
    backgroundColor: '#e5e7eb',
    overflow: 'hidden',
    marginTop: '4px',
    marginBottom: '8px',
  },
  barFill: {
    height: '100%',
    borderRadius: '4px',
  },
};

export const StatsPanel: React.FC<StatsPanelProps> = ({ stats, country, seed, referenceYear }) => {
  const livingPercent = stats.total_persons > 0
    ? (stats.living_persons / stats.total_persons) * 100
    : 0;

  const malePercent = stats.total_persons > 0
    ? (stats.male_count / stats.total_persons) * 100
    : 0;

  const marriedPercent = stats.total_persons > 0
    ? (stats.married_count / stats.total_persons) * 100
    : 0;

  const tertiaryPercent = stats.total_persons > 0
    ? (stats.tertiary_education / stats.total_persons) * 100
    : 0;

  const richPercent = stats.total_persons > 0
    ? (stats.rich_count / stats.total_persons) * 100
    : 0;

  const formatCurrency = (value: number) => new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  }).format(value);

  return (
    <div style={styles.panel}>
      <div style={styles.title}>Tree Statistics</div>

      <div style={styles.stat}>
        <span style={styles.label}>Country:</span>
        <span style={styles.value}>{country.replace(/-/g, ' ')}</span>
      </div>

      <div style={styles.stat}>
        <span style={styles.label}>Seed:</span>
        <span style={styles.value}>{seed}</span>
      </div>

      <div style={styles.stat}>
        <span style={styles.label}>Current Year:</span>
        <span style={styles.value}>{referenceYear ?? new Date().getFullYear()}</span>
      </div>

      { }
      <div style={styles.section}>
        <div style={styles.sectionTitle}>Population</div>
        <div style={styles.stat}>
          <span style={styles.label}>Total Persons:</span>
          <span style={styles.value}>{stats.total_persons}</span>
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Families:</span>
          <span style={styles.value}>{stats.total_families}</span>
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Living / Deceased:</span>
          <span style={styles.value}>{stats.living_persons} / {stats.deceased_persons}</span>
        </div>
        <div style={styles.bar}>
          <div style={{ ...styles.barFill, width: `${livingPercent}%`, backgroundColor: '#22c55e' }} />
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Male / Female:</span>
          <span style={styles.value}>{stats.male_count} / {stats.female_count}</span>
        </div>
        <div style={styles.bar}>
          <div style={{ ...styles.barFill, width: `${malePercent}%`, backgroundColor: '#4a90d9' }} />
        </div>
      </div>

      { }
      <div style={styles.section}>
        <div style={styles.sectionTitle}>Age</div>
        <div style={styles.stat}>
          <span style={styles.label}>Average Age:</span>
          <span style={styles.value}>{Math.round(stats.average_age)} years</span>
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Oldest Person:</span>
          <span style={styles.value}>{stats.oldest_person_age} years</span>
        </div>
      </div>

      { }
      <div style={styles.section}>
        <div style={styles.sectionTitle}>Family</div>
        <div style={styles.stat}>
          <span style={styles.label}>Total Children:</span>
          <span style={styles.value}>{stats.total_children}</span>
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Avg Children/Family:</span>
          <span style={styles.value}>{stats.average_children.toFixed(1)}</span>
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Married:</span>
          <span style={styles.value}>{stats.married_count} ({marriedPercent.toFixed(0)}%)</span>
        </div>
        <div style={styles.bar}>
          <div style={{ ...styles.barFill, width: `${marriedPercent}%`, backgroundColor: '#f59e0b' }} />
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Single:</span>
          <span style={styles.value}>{stats.single_count}</span>
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Divorces:</span>
          <span style={styles.value}>{stats.divorce_count}</span>
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Born Outside Marriage:</span>
          <span style={styles.value}>{stats.births_outside_marriage}</span>
        </div>
      </div>

      { }
      <div style={styles.section}>
        <div style={styles.sectionTitle}>Education & Work</div>
        <div style={styles.stat}>
          <span style={styles.label}>University Educated:</span>
          <span style={styles.value}>{stats.tertiary_education} ({tertiaryPercent.toFixed(0)}%)</span>
        </div>
        <div style={styles.bar}>
          <div style={{ ...styles.barFill, width: `${tertiaryPercent}%`, backgroundColor: '#8b5cf6' }} />
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Employed:</span>
          <span style={styles.value}>{stats.employed_count}</span>
        </div>
      </div>

      <div style={styles.section}>
        <div style={styles.sectionTitle}>Economy</div>
        <div style={styles.stat}>
          <span style={styles.label}>Avg GDP/Capita:</span>
          <span style={styles.value}>{formatCurrency(stats.average_gdp_per_capita)}</span>
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Avg Wealth Index:</span>
          <span style={styles.value}>{stats.average_wealth_index.toFixed(2)}</span>
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Avg Family Wealth:</span>
          <span style={styles.value}>{formatCurrency(stats.average_family_wealth)}</span>
        </div>

        <div style={styles.stat}>
          <span style={styles.label}>Rich:</span>
          <span style={styles.value}>{stats.rich_count} ({richPercent.toFixed(0)}%)</span>
        </div>
        <div style={styles.bar}>
          <div style={{ ...styles.barFill, width: `${richPercent}%`, backgroundColor: '#10b981' }} />
        </div>
      </div>
    </div>
  );
};
