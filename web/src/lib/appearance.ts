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
  '--hublet-background-blur',
  '--hublet-background-brightness',
  '--hublet-background-overlay',
  '--hublet-font-family',
  '--hublet-font-scale',
  '--hublet-card-radius',
  '--hublet-card-shadow',
  '--hublet-card-border-width',
  '--hublet-card-padding',
  '--hublet-card-min-height',
  '--hublet-card-gap',
  '--hublet-card-backdrop-blur',
  '--hublet-card-background-alpha'
] as const;

function clamp(
  value: number,
  minimum: number,
  maximum: number
): number {
  return Math.min(
    maximum,
    Math.max(minimum, value)
  );
}

function parseHex(
  value: string
): {
  red: number;
  green: number;
  blue: number;
} {
  const normalized =
    /^#[0-9a-fA-F]{6}$/.test(value)
      ? value.slice(1)
      : '090c12';

  return {
    red: Number.parseInt(
      normalized.slice(0, 2),
      16
    ),

    green: Number.parseInt(
      normalized.slice(2, 4),
      16
    ),

    blue: Number.parseInt(
      normalized.slice(4, 6),
      16
    )
  };
}

function componentToHex(
  value: number
): string {
  return clamp(
    Math.round(value),
    0,
    255
  )
    .toString(16)
    .padStart(2, '0');
}

function shiftHex(
  value: string,
  amount: number
): string {
  const color = parseHex(value);

  return [
    '#',
    componentToHex(color.red + amount),
    componentToHex(color.green + amount),
    componentToHex(color.blue + amount)
  ].join('');
}

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

function fontScale(
  scale: HubletConfig['appearance']['font']['scale']
): string {
  switch (scale) {
    case 'small':
      return '0.92';

    case 'large':
      return '1.09';

    case 'medium':
    default:
      return '1';
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

function cardDensity(
  density: HubletConfig['appearance']['cards']['density']
): {
  padding: string;
  minimumHeight: string;
  gap: string;
} {
  switch (density) {
    case 'compact':
      return {
        padding: '7px 9px',
        minimumHeight: '52px',
        gap: '7px'
      };

    case 'relaxed':
      return {
        padding: '15px 16px',
        minimumHeight: '82px',
        gap: '12px'
      };

    case 'comfortable':
    default:
      return {
        padding: '10px 12px',
        minimumHeight: '65px',
        gap: '9px'
      };
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

  if (
    background.type === 'gradient'
  ) {
    const lighter =
      shiftHex(
        background.color,
        26
      );

    const darker =
      shiftHex(
        background.color,
        -34
      );

    return [
      'radial-gradient(',
      'circle at 18% 0%,',
      `${lighter} 0%,`,
      'transparent 42%',
      '),',
      'linear-gradient(',
      '145deg,',
      `${background.color} 0%,`,
      `${darker} 100%`,
      ')'
    ].join(' ');
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

  const density =
    cardDensity(
      appearance.cards.density
    );

  const advanced =
    appearance.mode === 'advanced';

  return {
    '--hublet-background-color':
      background.color,

    '--hublet-background-image':
      backgroundImage(config),

    '--hublet-background-size':
      background.type === 'wallpaper'
        ? 'cover'
        : 'auto',

    '--hublet-background-position':
      'center center',

    '--hublet-background-blur':
      `${background.blur}px`,

    '--hublet-background-brightness':
      String(
        background.brightness / 100
      ),

    '--hublet-background-overlay':
      String(
        background.overlay / 100
      ),

    '--hublet-font-family':
      fontFamily(
        appearance.font.family
      ),

    '--hublet-font-scale':
      fontScale(
        appearance.font.scale
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
      density.padding,

    '--hublet-card-min-height':
      density.minimumHeight,

    '--hublet-card-gap':
      density.gap,

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
