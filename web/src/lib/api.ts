import type {
  HubletConfig
} from './types';

export type WallpaperFile = {
  filename: string;
  url: string;
  size: number;
  modified: string;
};

async function readError(
  response: Response
): Promise<string> {
  try {
    const body =
      await response.json();

    if (
      typeof body?.error ===
      'string'
    ) {
      return body.error;
    }
  } catch {
    // Ignore invalid error bodies.
  }

  return `HTTP ${response.status}`;
}

export async function loadConfig():
Promise<HubletConfig> {
  const response =
    await fetch(
      '/api/v2/config'
    );

  if (!response.ok) {
    throw new Error(
      `Unable to load configuration: ${
        await readError(response)
      }`
    );
  }

  return response.json();
}

export async function saveConfig(
  config: HubletConfig
): Promise<void> {
  const response =
    await fetch(
      '/api/v2/config',
      {
        method: 'PUT',

        headers: {
          'Content-Type':
            'application/json'
        },

        body:
          JSON.stringify(
            config
          )
      }
    );

  if (!response.ok) {
    throw new Error(
      `Unable to save configuration: ${
        await readError(response)
      }`
    );
  }
}

export async function listWallpapers():
Promise<WallpaperFile[]> {
  const response =
    await fetch(
      '/api/v2/wallpapers'
    );

  if (!response.ok) {
    throw new Error(
      `Unable to load wallpapers: ${
        await readError(response)
      }`
    );
  }

  return response.json();
}

export async function uploadWallpaper(
  file: File
): Promise<WallpaperFile> {
  const form =
    new FormData();

  form.append(
    'file',
    file
  );

  const response =
    await fetch(
      '/api/v2/wallpapers/upload',
      {
        method: 'POST',
        body: form
      }
    );

  if (!response.ok) {
    throw new Error(
      `Unable to upload wallpaper: ${
        await readError(response)
      }`
    );
  }

  return response.json();
}

export async function downloadWallpaper(
  url: string
): Promise<WallpaperFile> {
  const response =
    await fetch(
      '/api/v2/wallpapers/download',
      {
        method: 'POST',

        headers: {
          'Content-Type':
            'application/json'
        },

        body:
          JSON.stringify({
            url
          })
      }
    );

  if (!response.ok) {
    throw new Error(
      `Unable to download wallpaper: ${
        await readError(response)
      }`
    );
  }

  return response.json();
}

export async function deleteWallpaper(
  filename: string
): Promise<void> {
  const response =
    await fetch(
      `/api/v2/wallpapers/${
        encodeURIComponent(
          filename
        )
      }`,
      {
        method: 'DELETE'
      }
    );

  if (!response.ok) {
    throw new Error(
      `Unable to delete wallpaper: ${
        await readError(response)
      }`
    );
  }
}
