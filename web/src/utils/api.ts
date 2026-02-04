import { GenerateRequest, GenerateResponse, CountriesResponse } from '../types';

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const data = await response.json().catch(() => ({}));
    throw new ApiError(response.status, data.error || 'Request failed');
  }
  return response.json();
}

export async function getCountries(): Promise<CountriesResponse> {
  const response = await fetch(`${API_BASE}/api/countries`);
  return handleResponse<CountriesResponse>(response);
}

export async function getCountryStats(slug: string): Promise<unknown> {
  const response = await fetch(`${API_BASE}/api/country/${slug}`);
  return handleResponse(response);
}

export async function generateTree(request: GenerateRequest): Promise<GenerateResponse> {
  const response = await fetch(`${API_BASE}/api/generate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(request),
  });
  return handleResponse<GenerateResponse>(response);
}

export async function checkHealth(): Promise<{ status: string }> {
  const response = await fetch(`${API_BASE}/api/health`);
  return handleResponse(response);
}

// Helper to check if API is available
export async function isApiAvailable(): Promise<boolean> {
  try {
    await checkHealth();
    return true;
  } catch {
    return false;
  }
}
