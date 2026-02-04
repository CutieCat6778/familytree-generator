import React from 'react';
import { VisualizationNode } from '../types';
import { getPersonColor } from '../utils/treeLayout';

interface PersonCardProps {
  person: VisualizationNode;
  isSelected: boolean;
  isRoot: boolean;
  referenceYear?: number;
  onClick: () => void;
}

const styles: Record<string, React.CSSProperties> = {
  card: {
    padding: '8px 12px',
    borderRadius: '8px',
    cursor: 'pointer',
    textAlign: 'center',
    minWidth: '120px',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    transition: 'all 0.2s ease',
    border: '2px solid transparent',
  },
  name: {
    fontWeight: 'bold',
    fontSize: '14px',
    marginBottom: '4px',
    color: '#fff',
  },
  years: {
    fontSize: '12px',
    color: 'rgba(255,255,255,0.9)',
  },
  deceased: {
    opacity: 0.8,
  },
  root: {
    border: '2px solid gold',
  },
  selected: {
    border: '2px solid #333',
    transform: 'scale(1.05)',
  }
};

export const PersonCard: React.FC<PersonCardProps> = ({ person, isSelected, isRoot, referenceYear, onClick }) => {
  const bgColor = getPersonColor(person.gender, person.is_alive);

  const years = person.death_year
    ? `${person.birth_year} - ${person.death_year}`
    : `b. ${person.birth_year}`;

  const age = person.death_year
    ? person.death_year - person.birth_year
    : Math.max(0, (referenceYear ?? new Date().getFullYear()) - person.birth_year);

  return (
    <div
      style={{
        ...styles.card,
        backgroundColor: bgColor,
        ...(isRoot ? styles.root : {}),
        ...(isSelected ? styles.selected : {}),
        ...(!person.is_alive ? styles.deceased : {}),
      }}
      onClick={onClick}
      title={`${person.name}\nAge: ${age}\n${person.is_alive ? 'Living' : 'Deceased'}`}
    >
      <div style={styles.name}>{person.name}</div>
      <div style={styles.years}>{years}</div>
    </div>
  );
};
