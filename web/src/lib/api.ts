import type { HubletConfig } from './types';

async function readError(response: Response): Promise<string> {
  try {
    const body = await response.json();

    if (typeof body?.error === 'string') {
      return body.error;
    }
  } catch {
    // Ignore invalid error response bodies.
  }

  return `HTTP ${response.status}`;
}

export async function loadConfig(): Promise<HubletConfig> {
  const response = await fetch('/api/v2/config');

  if (!response.ok) {
    throw new Error(
      `Unable to load configuration: ${await readError(response)}`
    );
  }

  return response.json();
}

export async function saveConfig(
  config: HubletConfig
): Promise<void> {
  const response = await fetch('/api/v2/config', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(config)
  });

  if (!response.ok) {
    throw new Error(
      `Unable to save configuration: ${await readError(response)}`
    );
  }
}
