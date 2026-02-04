import React from 'react';
import { VisualizationNode } from '../types';
import { getPersonColor } from '../utils/treeLayout';

interface PersonDetailProps {
  person: VisualizationNode;
  referenceYear?: number;
  onClose: () => void;
}

const styles: Record<string, React.CSSProperties> = {
  overlay: {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0,0,0,0.5)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
  },
  modal: {
    backgroundColor: '#fff',
    borderRadius: '12px',
    padding: '24px',
    minWidth: '380px',
    maxWidth: '500px',
    maxHeight: '80vh',
    overflow: 'auto',
    boxShadow: '0 4px 20px rgba(0,0,0,0.2)',
  },
  header: {
    display: 'flex',
    alignItems: 'center',
    marginBottom: '20px',
    gap: '16px',
  },
  avatar: {
    width: '60px',
    height: '60px',
    borderRadius: '50%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    color: '#fff',
    fontSize: '24px',
    fontWeight: 'bold',
  },
  name: {
    fontSize: '20px',
    fontWeight: 'bold',
    color: '#333',
  },
  status: {
    fontSize: '14px',
    color: '#666',
    marginTop: '4px',
  },
  section: {
    marginTop: '16px',
    paddingTop: '16px',
    borderTop: '1px solid #eee',
  },
  sectionTitle: {
    fontSize: '14px',
    fontWeight: 'bold',
    color: '#333',
    marginBottom: '12px',
  },
  row: {
    display: 'flex',
    justifyContent: 'space-between',
    marginBottom: '8px',
    fontSize: '14px',
  },
  label: {
    color: '#666',
  },
  value: {
    fontWeight: '500',
    color: '#333',
  },
  badge: {
    display: 'inline-block',
    padding: '2px 8px',
    borderRadius: '12px',
    fontSize: '12px',
    fontWeight: '500',
  },
  badgeGreen: {
    backgroundColor: '#dcfce7',
    color: '#166534',
  },
  badgeRed: {
    backgroundColor: '#fee2e2',
    color: '#dc2626',
  },
  badgeBlue: {
    backgroundColor: '#dbeafe',
    color: '#1d4ed8',
  },
  badgeYellow: {
    backgroundColor: '#fef3c7',
    color: '#92400e',
  },
  closeBtn: {
    marginTop: '20px',
    width: '100%',
    padding: '10px',
    backgroundColor: '#4a90d9',
    color: '#fff',
    border: 'none',
    borderRadius: '6px',
    fontSize: '14px',
    cursor: 'pointer',
  }
};

const formatMaritalStatus = (status: string): string => {
  const statusMap: Record<string, string> = {
    'single': 'Single',
    'married': 'Married',
    'divorced': 'Divorced',
    'widowed': 'Widowed',
    'remarried': 'Remarried',
  };
  return statusMap[status] || status;
};

const formatEducation = (education: string): string => {
  const educationMap: Record<string, string> = {
    'none': 'No Formal Education',
    'primary': 'Primary School',
    'secondary': 'Secondary School',
    'tertiary': 'University/College',
  };
  return educationMap[education] || education;
};

const formatEmployment = (employment: string): string => {
  const employmentMap: Record<string, string> = {
    'employed': 'Employed',
    'unemployed': 'Unemployed',
    'retired': 'Retired',
    'student': 'Student',
    'child': 'Child',
  };
  return employmentMap[employment] || employment;
};

const formatCountry = (slug: string | undefined): string => {
  if (!slug) return 'Unknown';
  return slug.replace(/-/g, ' ');
};

const formatCurrency = (value: number | undefined): string => {
  if (value === undefined || Number.isNaN(value)) return 'N/A';
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  }).format(value);
};

export const PersonDetail: React.FC<PersonDetailProps> = ({ person, referenceYear, onClose }) => {
  const bgColor = getPersonColor(person.gender, person.is_alive);
  const currentYear = referenceYear ?? new Date().getFullYear();
  const age = person.death_year
    ? person.death_year - person.birth_year
    : Math.max(0, currentYear - person.birth_year);
  const currentCountry = person.current_country || person.country;
  const isEmigrant = currentCountry !== person.country;
  const wealthRatio = person.wealth_index !== undefined
    ? person.wealth_index
    : (person.family_wealth !== undefined && person.gdp_per_capita)
      ? person.family_wealth / person.gdp_per_capita
      : undefined;

  const initial = person.first_name.charAt(0).toUpperCase();

  return (
    <div style={styles.overlay} onClick={onClose}>
      <div style={styles.modal} onClick={e => e.stopPropagation()}>
        <div style={styles.header}>
          <div style={{ ...styles.avatar, backgroundColor: bgColor }}>
            {initial}
          </div>
          <div>
            <div style={styles.name}>{person.name}</div>
            <div style={styles.status}>
              <span style={{
                ...styles.badge,
                ...(person.is_alive ? styles.badgeGreen : styles.badgeRed)
              }}>
                {person.is_alive ? 'Living' : 'Deceased'}
              </span>
              {' '}
              <span style={{ ...styles.badge, ...styles.badgeBlue }}>
                {person.gender === 'M' ? 'Male' : 'Female'}
              </span>
            </div>
          </div>
        </div>

        {             }
        <div style={styles.section}>
          <div style={styles.sectionTitle}>Basic Information</div>
          <div style={styles.row}>
            <span style={styles.label}>Birth Year:</span>
            <span style={styles.value}>{person.birth_year}</span>
          </div>

          {person.death_year && (
            <div style={styles.row}>
              <span style={styles.label}>Death Year:</span>
              <span style={styles.value}>{person.death_year}</span>
            </div>
          )}

          <div style={styles.row}>
            <span style={styles.label}>Age:</span>
            <span style={styles.value}>
              {age} years{person.is_alive ? ' old' : ' (at death)'}
            </span>
          </div>

          <div style={styles.row}>
            <span style={styles.label}>Generation:</span>
            <span style={styles.value}>
              {person.generation === 0 ? 'Root' :
               person.generation > 0 ? `+${person.generation} (descendant)` :
               `${person.generation} (ancestor)`}
            </span>
          </div>

          <div style={styles.row}>
            <span style={styles.label}>Birth Country:</span>
            <span style={styles.value}>{formatCountry(person.country)}</span>
          </div>
        </div>

        {                     }
        <div style={styles.section}>
          <div style={styles.sectionTitle}>Location & Economy</div>
          <div style={styles.row}>
            <span style={styles.label}>Current Country:</span>
            <span style={styles.value}>{formatCountry(currentCountry)}</span>
          </div>
          <div style={styles.row}>
            <span style={styles.label}>Emigrant:</span>
            <span style={styles.value}>
              {isEmigrant ? `Yes (from ${formatCountry(person.country)})` : 'No'}
            </span>
          </div>
          <div style={styles.row}>
            <span style={styles.label}>Residence:</span>
            <span style={styles.value}>
              {person.residence ? `${person.residence.charAt(0).toUpperCase()}${person.residence.slice(1)}` : 'Unknown'}
            </span>
          </div>
          <div style={styles.row}>
            <span style={styles.label}>GDP per Capita:</span>
            <span style={styles.value}>{formatCurrency(person.gdp_per_capita)}</span>
          </div>
          <div style={styles.row}>
            <span style={styles.label}>Family Wealth:</span>
            <span style={{
              ...styles.badge,
              ...(person.is_rich ? styles.badgeGreen : styles.badgeYellow)
            }}>
              {person.is_rich ? 'Rich' : 'Not Rich'}
            </span>
          </div>
          <div style={styles.row}>
            <span style={styles.label}>Wealth vs GDP:</span>
            <span style={styles.value}>
              {wealthRatio !== undefined ? `${wealthRatio.toFixed(2)}x` : 'N/A'}
            </span>
          </div>
        </div>

        {              }
        <div style={styles.section}>
          <div style={styles.sectionTitle}>Family</div>
          <div style={styles.row}>
            <span style={styles.label}>Marital Status:</span>
            <span style={styles.value}>{formatMaritalStatus(person.marital_status)}</span>
          </div>

          {person.marital_status === 'divorced' && (
            <div style={styles.row}>
              <span style={styles.label}>Divorced:</span>
              <span style={{ ...styles.badge, ...styles.badgeRed }}>Yes</span>
            </div>
          )}

          {(person.marriage_age ?? 0) > 0 && (
            <div style={styles.row}>
              <span style={styles.label}>Married at Age:</span>
              <span style={styles.value}>{person.marriage_age}</span>
            </div>
          )}

          <div style={styles.row}>
            <span style={styles.label}>Children:</span>
            <span style={styles.value}>{person.number_of_children}</span>
          </div>

          {person.born_outside_marriage && (
            <div style={styles.row}>
              <span style={styles.label}>Born Outside Marriage:</span>
              <span style={{ ...styles.badge, ...styles.badgeYellow }}>Yes</span>
            </div>
          )}

          {person.is_single_parent && (
            <div style={styles.row}>
              <span style={styles.label}>Single Parent:</span>
              <span style={{ ...styles.badge, ...styles.badgeYellow }}>Yes</span>
            </div>
          )}
        </div>

        {                         }
        <div style={styles.section}>
          <div style={styles.sectionTitle}>Education & Employment</div>
          <div style={styles.row}>
            <span style={styles.label}>Education:</span>
            <span style={styles.value}>{formatEducation(person.education)}</span>
          </div>

          <div style={styles.row}>
            <span style={styles.label}>Employment:</span>
            <span style={styles.value}>{formatEmployment(person.employment)}</span>
          </div>
        </div>

        {         }
        <div style={styles.section}>
          <div style={styles.sectionTitle}>Health Factors</div>
          <div style={styles.row}>
            <span style={styles.label}>Alcohol (L/year):</span>
            <span style={styles.value}>{person.alcohol_consumption.toFixed(1)}</span>
          </div>

          <div style={styles.row}>
            <span style={styles.label}>Tobacco Use:</span>
            <span style={{
              ...styles.badge,
              ...(person.tobacco_use ? styles.badgeRed : styles.badgeGreen)
            }}>
              {person.tobacco_use ? 'Yes' : 'No'}
            </span>
          </div>

          <div style={styles.row}>
            <span style={styles.label}>Underweight (U5):</span>
            <span style={{
              ...styles.badge,
              ...(person.underweight ? styles.badgeRed : styles.badgeGreen)
            }}>
              {person.underweight ? 'Yes' : 'No'}
            </span>
          </div>
        </div>

        <button style={styles.closeBtn} onClick={onClose}>
          Close
        </button>
      </div>
    </div>
  );
};
