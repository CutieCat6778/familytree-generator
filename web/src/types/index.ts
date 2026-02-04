

export interface VisualizationData {
  id: string;
  root_id: string;
  country: string;
  generations: number;
  seed: number;
  reference_year?: number;
  nodes: VisualizationNode[];
  edges: VisualizationEdge[];
  stats: VisualizationStats;
}

export interface VisualizationNode {
  id: string;
  name: string;
  first_name: string;
  last_name: string;
  gender: 'M' | 'F';
  birth_year: number;
  death_year?: number;
  is_alive: boolean;
  generation: number;
  marital_status: string;
  marriage_age?: number;
  number_of_children: number;
  education: string;
  employment: string;
  alcohol_consumption: number;
  tobacco_use: boolean;
  born_outside_marriage: boolean;
  is_single_parent: boolean;
  underweight?: boolean;
  residence?: 'urban' | 'rural';
  gdp_per_capita?: number;
  wealth_index?: number;
  family_wealth?: number;
  is_rich?: boolean;
  country: string;
  current_country?: string;
}

export interface VisualizationEdge {
  source: string;
  target: string;
  type: 'parent' | 'spouse';
}

export interface VisualizationStats {
  total_persons: number;
  total_families: number;
  living_persons: number;
  deceased_persons: number;
  average_age: number;
  oldest_person_age: number;
  total_children: number;
  average_children: number;
  divorce_count: number;
  single_count: number;
  married_count: number;
  male_count: number;
  female_count: number;
  births_outside_marriage: number;
  tertiary_education: number;
  employed_count: number;
  average_gdp_per_capita: number;
  average_wealth_index: number;
  average_family_wealth: number;
  rich_count: number;
}


export interface TreeNode extends VisualizationNode {
  children?: TreeNode[];
  spouse?: VisualizationNode;
  _children?: TreeNode[]; 
}

export interface HierarchyNode {
  data: TreeNode;
  x: number;
  y: number;
  children?: HierarchyNode[];
}


export interface GenerateRequest {
  country: string;
  generations: number;
  seed?: number;
  start_year?: number;
  gender?: 'M' | 'F';
  include_extended?: boolean;
  life_expectancy_mode?: 'total' | 'female' | 'male' | 'by_gender';
}

export interface GenerateResponse {
  success: boolean;
  message?: string;
  tree?: VisualizationData;
  stats?: {
    generation_time: string;
  };
}

export interface CountryInfo {
  slug: string;
  name: string;
  iso_code: string;
  has_name_data: boolean;
  population?: number;
  life_expectancy?: number;
}

export interface CountriesResponse {
  countries: CountryInfo[];
  count: number;
}
