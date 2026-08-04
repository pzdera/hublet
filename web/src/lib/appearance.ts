import type {
  HubletConfig
} from './types';

export type AppearanceVariables =
  Record<string, string>;

const variableNames = [
  '--hublet-background-color',
  '--hublet-background-image',
  '--hublet-background-size',
  '--hublet-background-position',
  '--hublet-background-transform-scale',
  '--hublet-background-blur',
  '--hublet-background-brightness',
  '--hublet-background-overlay',
  '--hublet-font-family',
  '--hublet-card-radius',
  '--hublet-card-shadow',
  '--hublet-card-border-width',
  '--hublet-card-padding',
  '--hublet-card-min-height',
  '--hublet-card-gap',
  '--hublet-card-backdrop-blur',
  '--hublet-card-background-alpha'
] as const;

function fontFamily(
  family: HubletConfig['appearance']['font']['family']
): string {
  switch (family) {
    case 'inter':
      return [
        'Inter Variable',
        'ui-sans-serif',
        'system-ui',
        'sans-serif'
      ].join(', ');

    case 'geist':
      return [
        'Geist Variable',
        'Inter Variable',
        'ui-sans-serif',
        'system-ui',
        'sans-serif'
      ].join(', ');

    case 'manrope':
      return [
        'Manrope Variable',
        'Inter Variable',
        'ui-sans-serif',
        'system-ui',
        'sans-serif'
      ].join(', ');

    case 'ibm-plex-sans':
      return [
        '"IBM Plex Sans Variable"',
        'Inter Variable',
        'ui-sans-serif',
        'system-ui',
        'sans-serif'
      ].join(', ');

    case 'system':
    default:
      return [
        'ui-sans-serif',
        'system-ui',
        '-apple-system',
        'BlinkMacSystemFont',
        '"Segoe UI"',
        'sans-serif'
      ].join(', ');
  }
}

function cardRadius(
  radius: HubletConfig['appearance']['cards']['radius']
): string {
  switch (radius) {
    case 'small':
      return '9px';

    case 'medium':
      return '14px';

    case 'large':
    default:
      return '20px';
  }
}

function cardShadow(
  shadow: HubletConfig['appearance']['cards']['shadow']
): string {
  switch (shadow) {
    case 'none':
      return 'none';

    case 'medium':
      return [
        '0 14px 36px',
        'rgba(0, 0, 0, 0.28)'
      ].join(' ');

    case 'floating':
      return [
        '0 22px 60px',
        'rgba(0, 0, 0, 0.42)'
      ].join(' ');

    case 'soft':
    default:
      return [
        '0 8px 24px',
        'rgba(0, 0, 0, 0.18)'
      ].join(' ');
  }
}

function backgroundImage(
  config: HubletConfig
): string {
  const background =
    config.appearance.background;

  if (
    background.type === 'wallpaper' &&
    background.wallpaper
  ) {
    const filename =
      encodeURIComponent(
        background.wallpaper
      );

    return `url("/wallpapers/${filename}")`;
  }

  return 'none';
}

export function appearanceVariables(
  config: HubletConfig
): AppearanceVariables {
  const appearance =
    config.appearance;

  const background =
    appearance.background;

  const advanced =
    appearance.mode === 'advanced';

  const wallpaperActive =
    background.type === 'wallpaper' &&
    Boolean(background.wallpaper);

  return {
    '--hublet-background-color':
      background.color,

    '--hublet-background-image':
      backgroundImage(config),

    '--hublet-background-size':
      background.type === 'wallpaper'
        ? background.fit === 'contain'
          ? 'contain'
          : 'cover'
        : 'auto',

    '--hublet-background-position':
      'center center',

    '--hublet-background-transform-scale':
      background.type === 'wallpaper' &&
      background.fit !== 'contain'
        ? '1.035'
        : '1',

    '--hublet-background-blur':
      wallpaperActive
        ? `${background.blur}px`
        : '0px',

    '--hublet-background-brightness':
      wallpaperActive
        ? String(
            background.brightness / 100
          )
        : '1',

    '--hublet-background-overlay':
      wallpaperActive
        ? String(
            background.overlay / 100
          )
        : '0',

    '--hublet-font-family':
      fontFamily(
        appearance.font.family
      ),

    '--hublet-card-radius':
      cardRadius(
        appearance.cards.radius
      ),

    '--hublet-card-shadow':
      cardShadow(
        appearance.cards.shadow
      ),

    '--hublet-card-border-width':
      appearance.cards.border
        ? '1px'
        : '0px',

    '--hublet-card-padding':
      '10px 12px',

    '--hublet-card-min-height':
      '65px',

    '--hublet-card-gap':
      '9px',

    '--hublet-card-backdrop-blur':
      advanced
        ? '18px'
        : '0px',

    '--hublet-card-background-alpha':
      advanced
        ? '0.09'
        : appearance.mode === 'standard'
          ? '0.065'
          : '0.045'
  };
}

export function appearanceStyle(
  config: HubletConfig
): string {
  return Object.entries(
    appearanceVariables(config)
  )
    .map(
      ([name, value]) =>
        `${name}:${value}`
    )
    .join(';');
}

export function applyAppearanceVariables(
  target: HTMLElement,
  config: HubletConfig
): void {
  const variables =
    appearanceVariables(config);

  for (
    const [name, value]
    of Object.entries(variables)
  ) {
    target.style.setProperty(
      name,
      value
    );
  }
}

export function captureAppearanceVariables(
  target: HTMLElement
): Map<string, string> {
  return new Map(
    variableNames.map(
      (name) => [
        name,
        target.style.getPropertyValue(
          name
        )
      ]
    )
  );
}

export function restoreAppearanceVariables(
  target: HTMLElement,
  snapshot: Map<string, string>
): void {
  for (
    const name
    of variableNames
  ) {
    const value =
      snapshot.get(name) ?? '';

    if (value) {
      target.style.setProperty(
        name,
        value
      );
    } else {
      target.style.removeProperty(
        name
      );
    }
  }
}
