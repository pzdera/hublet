import type {
  Item
} from './types';

const dashboardIconsBase =
  'https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png';

const knownAliases: Record<string, string> = {
  adguardhome: 'adguard-home',
  adguard: 'adguard-home',

  autobrr: 'autobrr',

  bazarr: 'bazarr',

  caddy: 'caddy',

  dockhand: 'docker',
  docker: 'docker',

  frigate: 'frigate',

  github: 'github',
  gitlab: 'gitlab',

  gotify: 'gotify',

  homeassistant: 'home-assistant',
  'home-assistant': 'home-assistant',

  homepage: 'homepage',

  immich: 'immich',

  jellyfin: 'jellyfin',

  nextcloud: 'nextcloud',

  paperless: 'paperless-ngx',
  'paperless-ngx': 'paperless-ngx',

  plex: 'plex',

  prowlarr: 'prowlarr',

  proxmox: 'proxmox',
  pve: 'proxmox',
  pbs: 'proxmox-backup-server',

  qbittorrent: 'qbittorrent',
  qbit: 'qbittorrent',

  radarr: 'radarr',

  sonarr: 'sonarr',

  tailscale: 'tailscale',

  tdarr: 'tdarr',

  truenas: 'truenas',
  'truenas-scale': 'truenas-scale',

  uptimekuma: 'uptime-kuma',
  'uptime-kuma': 'uptime-kuma',

  vaultwarden: 'vaultwarden',
  bitwarden: 'bitwarden'
};

function normalizeCandidate(
  value: string
): string {
  return value
    .trim()
    .toLowerCase()
    .replace(/^www\./, '')
    .replace(/\.[a-z]{2,}$/i, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
}

function compactCandidate(
  value: string
): string {
  return normalizeCandidate(value)
    .replace(/-/g, '');
}

function hostnameFromURL(
  value: string
): string {
  const trimmed =
    value.trim();

  if (!trimmed) {
    return '';
  }

  try {
    const normalized =
      /^https?:\/\//i.test(trimmed)
        ? trimmed
        : `http://${trimmed}`;

    return new URL(
      normalized
    ).hostname;
  } catch {
    return '';
  }
}

function hostnameParts(
  hostname: string
): string[] {
  if (!hostname) {
    return [];
  }

  const firstPart =
    hostname
      .split('.')[0]
      .trim();

  return firstPart
    ? [firstPart]
    : [];
}

export function iconCandidates(
  item: Item
): string[] {
  const hostname =
    hostnameFromURL(
      item.url
    );

  const rawCandidates = [
    item.icon.value,
    item.name,
    ...hostnameParts(hostname)
  ];

  const candidates =
    new Set<string>();

  for (
    const rawCandidate
    of rawCandidates
  ) {
    const normalized =
      normalizeCandidate(
        rawCandidate
      );

    if (!normalized) {
      continue;
    }

    const compact =
      compactCandidate(
        normalized
      );

    const known =
      knownAliases[normalized] ??
      knownAliases[compact];

    if (known) {
      candidates.add(known);
    }

    candidates.add(normalized);

    const words =
      normalized.split('-');

    if (words.length > 1) {
      candidates.add(
        words[0]
      );
    }
  }

  return [
    ...candidates
  ];
}

export function dashboardIconURL(
  slug: string
): string {
  return `${
    dashboardIconsBase
  }/${
    encodeURIComponent(slug)
  }.png`;
}

export function serviceFaviconURL(
  value: string
): string {
  const trimmed =
    value.trim();

  if (!trimmed) {
    return '';
  }

  try {
    const normalized =
      /^https?:\/\//i.test(trimmed)
        ? trimmed
        : `http://${trimmed}`;

    const parsed =
      new URL(normalized);

    return `${
      parsed.origin
    }/favicon.ico`;
  } catch {
    return '';
  }
}

export function localIconURL(
  item: Item
): string {
  if (
    item.icon.type !== 'local' ||
    !item.icon.value.trim()
  ) {
    return '';
  }

  return `/icons/${
    encodeURIComponent(
      item.icon.value.trim()
    )
  }`;
}

export function itemInitial(
  item: Item
): string {
  return (
    item.name
      .trim()
      .slice(0, 1)
      .toUpperCase() ||
    '?'
  );
}
