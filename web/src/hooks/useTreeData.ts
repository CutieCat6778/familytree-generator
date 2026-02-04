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


const sampleData: VisualizationData = {
  id: "sample_tree",
  root_id: "P001",
  country: "united-states",
  generations: 3,
  seed: 12345,
  reference_year: 2018,
  nodes: [
    {
      id: "P001",
      name: "John Smith",
      first_name: "John",
      last_name: "Smith",
      gender: "M",
      birth_year: 1990,
      is_alive: true,
      generation: 0,
      marital_status: "married",
      number_of_children: 2,
      education: "tertiary",
      employment: "employed",
      alcohol_consumption: 3.2,
      tobacco_use: false,
      born_outside_marriage: false,
      is_single_parent: false,
      country: "united-states",
      current_country: "united-states",
      gdp_per_capita: 75000,
      residence: "urban",
      underweight: false,
      wealth_index: 1.2,
      family_wealth: 90000,
      is_rich: false,
    },
    {
      id: "P002",
      name: "Jane Smith",
      first_name: "Jane",
      last_name: "Smith",
      gender: "F",
      birth_year: 1992,
      is_alive: true,
      generation: 0,
      marital_status: "married",
      number_of_children: 2,
      education: "secondary",
      employment: "employed",
      alcohol_consumption: 1.1,
      tobacco_use: false,
      born_outside_marriage: false,
      is_single_parent: false,
      country: "united-states",
      current_country: "united-states",
      gdp_per_capita: 75000,
      residence: "urban",
      underweight: false,
      wealth_index: 1.1,
      family_wealth: 82500,
      is_rich: false,
    },
    {
      id: "P003",
      name: "Robert Smith",
      first_name: "Robert",
      last_name: "Smith",
      gender: "M",
      birth_year: 1960,
      death_year: 2020,
      is_alive: false,
      generation: -1,
      marital_status: "married",
      number_of_children: 1,
      education: "secondary",
      employment: "retired",
      alcohol_consumption: 2.1,
      tobacco_use: true,
      born_outside_marriage: false,
      is_single_parent: false,
      country: "united-states",
      current_country: "united-states",
      gdp_per_capita: 75000,
      residence: "rural",
      underweight: false,
      wealth_index: 1.0,
      family_wealth: 75000,
      is_rich: false,
    },
    {
      id: "P004",
      name: "Mary Smith",
      first_name: "Mary",
      last_name: "Smith",
      gender: "F",
      birth_year: 1962,
      is_alive: true,
      generation: -1,
      marital_status: "married",
      number_of_children: 1,
      education: "secondary",
      employment: "retired",
      alcohol_consumption: 0.8,
      tobacco_use: false,
      born_outside_marriage: false,
      is_single_parent: false,
      country: "united-states",
      current_country: "united-states",
      gdp_per_capita: 75000,
      residence: "rural",
      underweight: false,
      wealth_index: 0.9,
      family_wealth: 67500,
      is_rich: false,
    },
    {
      id: "P005",
      name: "James Smith",
      first_name: "James",
      last_name: "Smith",
      gender: "M",
      birth_year: 2015,
      is_alive: true,
      generation: 1,
      marital_status: "single",
      number_of_children: 0,
      education: "primary",
      employment: "child",
      alcohol_consumption: 0.0,
      tobacco_use: false,
      born_outside_marriage: false,
      is_single_parent: false,
      country: "united-states",
      current_country: "united-states",
      gdp_per_capita: 75000,
      residence: "urban",
      underweight: false,
      wealth_index: 1.0,
      family_wealth: 75000,
      is_rich: false,
    },
    {
      id: "P006",
      name: "Emily Smith",
      first_name: "Emily",
      last_name: "Smith",
      gender: "F",
      birth_year: 2018,
      is_alive: true,
      generation: 1,
      marital_status: "single",
      number_of_children: 0,
      education: "primary",
      employment: "child",
      alcohol_consumption: 0.0,
      tobacco_use: false,
      born_outside_marriage: false,
      is_single_parent: false,
      country: "united-states",
      current_country: "united-states",
      gdp_per_capita: 75000,
      residence: "urban",
      underweight: false,
      wealth_index: 1.0,
      family_wealth: 75000,
      is_rich: false,
    },
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
    deceased_persons: 1,
    average_age: 45,
    oldest_person_age: 62,
    total_children: 3,
    average_children: 1.5,
    divorce_count: 0,
    single_count: 2,
    married_count: 4,
    male_count: 3,
    female_count: 3,
    births_outside_marriage: 0,
    tertiary_education: 1,
    employed_count: 2,
    average_gdp_per_capita: 75000,
    average_wealth_index: 1.03,
    average_family_wealth: 77500,
    rich_count: 0,
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
