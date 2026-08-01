import type { HubletConfig } from './types';

export async function loadConfig(): Promise<HubletConfig> {
  const response = await fetch('/api/v2/config');

  if (!response.ok) {
    throw new Error(`Unable to load config: HTTP ${response.status}`);
  }

  return response.json();
}
