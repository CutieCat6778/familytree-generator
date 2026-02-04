import { useState, useCallback } from 'react';
import type { Dispatch, SetStateAction } from 'react';
import { VisualizationData } from '../types';

interface UseTreeDataReturn {
  data: VisualizationData | null;
  loading: boolean;
  error: string | null;
  loadFromFile: (file: File) => Promise<void>;
  loadFromUrl: (url: string) => Promise<void>;
  loadSample: () => void;
  setData: Dispatch<SetStateAction<VisualizationData | null>>;
  setError: Dispatch<SetStateAction<string | null>>;
}

// Sample data for demonstration
const sampleData: VisualizationData = {
  id: "sample_tree",
  root_id: "P001",
  country: "united-states",
  generations: 3,
  seed: 12345,
  reference_year: 2018,
  nodes: [
    { id: "P001", name: "John Smith", gender: "M", birth_year: 1990, is_alive: true, generation: 0 },
    { id: "P002", name: "Jane Smith", gender: "F", birth_year: 1992, is_alive: true, generation: 0 },
    { id: "P003", name: "Robert Smith", gender: "M", birth_year: 1960, death_year: 2020, is_alive: false, generation: -1 },
    { id: "P004", name: "Mary Smith", gender: "F", birth_year: 1962, is_alive: true, generation: -1 },
    { id: "P005", name: "James Smith", gender: "M", birth_year: 2015, is_alive: true, generation: 1 },
    { id: "P006", name: "Emily Smith", gender: "F", birth_year: 2018, is_alive: true, generation: 1 },
  ],
  edges: [
    { source: "P001", target: "P002", type: "spouse" },
    { source: "P003", target: "P001", type: "parent" },
    { source: "P004", target: "P001", type: "parent" },
    { source: "P001", target: "P005", type: "parent" },
    { source: "P002", target: "P005", type: "parent" },
    { source: "P001", target: "P006", type: "parent" },
    { source: "P002", target: "P006", type: "parent" },
  ],
  stats: {
    total_persons: 6,
    total_families: 2,
    living_persons: 5,
    average_age: 45,
    oldest_person_age: 62,
    total_children: 3
  }
};

export function useTreeData(): UseTreeDataReturn {
  const [data, setData] = useState<VisualizationData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadFromFile = useCallback(async (file: File) => {
    setLoading(true);
    setError(null);

    try {
      const text = await file.text();
      const json = JSON.parse(text);

      // Validate the data structure
      if (!json.nodes || !json.edges || !json.root_id) {
        throw new Error('Invalid file format. Expected visualization JSON with nodes, edges, and root_id.');
      }

      setData(json as VisualizationData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load file');
    } finally {
      setLoading(false);
    }
  }, []);

  const loadFromUrl = useCallback(async (url: string) => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const json = await response.json();

      if (!json.nodes || !json.edges || !json.root_id) {
        throw new Error('Invalid data format');
      }

      setData(json as VisualizationData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load from URL');
    } finally {
      setLoading(false);
    }
  }, []);

  const loadSample = useCallback(() => {
    setData(sampleData);
    setError(null);
  }, []);

  return {
    data,
    loading,
    error,
    loadFromFile,
    loadFromUrl,
    loadSample,
    setData,
    setError,
  };
}
