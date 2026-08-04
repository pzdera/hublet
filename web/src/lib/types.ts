export type Icon = {
  type: 'auto' | 'local' | 'none';
  value: string;
};

export type ItemType =
  | 'service';

export type Item = {
  id: string;
  type: ItemType;
  name: string;
  url: string;
  description: string;
  icon: Icon;
  openInNewTab: boolean;
};

export type SectionWidth =
  | 'narrow'
  | 'medium'
  | 'wide'
  | 'extra-wide'
  | 'full';

export type SectionSurface =
  | 'solid'
  | 'transparent'
  | 'glass'
  | 'none';

export type Section = {
  id: string;
  title: string;
  accent: string;
  surface: SectionSurface;
  surfaceOpacity: number;
  surfaceBlur: number;
  showBorder: boolean;
  width: SectionWidth;
  gridRow: number;
  gridColumn: number;
  gridRowSpan: number;
  gridColumnSpan: number;
  gridColumns: number;
  items: Item[];
};

export type Shortcut = {
  key: string;
  label: string;
  url: string;
  icon: Icon;
};

export type HubletConfig = {
  version: 2;

  dashboard: {
    title: string;
    description: string;
    descriptionVisible: boolean;

    icon: {
      type:
        | 'initial'
        | 'local'
        | 'none';

      value: string;
    };

    iconSize:
      | 'small'
      | 'medium'
      | 'large';

    theme: string;
    wallpaper: string | null;
  };

  appearance: {
    mode:
      | 'minimal'
      | 'standard'
      | 'advanced';

    font: {
      family:
        | 'system'
        | 'inter'
        | 'geist'
        | 'manrope'
        | 'ibm-plex-sans';
    };

    cards: {
      radius:
        | 'small'
        | 'medium'
        | 'large';

      shadow:
        | 'none'
        | 'soft'
        | 'medium'
        | 'floating';

      border: boolean;
    };

    background: {
      type:
        | 'solid'
        | 'wallpaper';

      color: string;
      blur: number;
      brightness: number;
      overlay: number;
      wallpaper: string | null;
    };
  };

  search: {
    autoFocus: boolean;
    openShortcutDirectly: boolean;
    webSearchEnabled: boolean;

    webSearchEngine:
      | 'google'
      | 'duckduckgo'
      | 'bing';
  };

  sections: Section[];
  shortcuts: Shortcut[];

  modules: {
    weather: {
      enabled: boolean;

      mode:
        | 'current'
        | 'today'
        | 'five-day';

      location: string;
      latitude: number | null;
      longitude: number | null;
    };

    clock: {
      enabled: boolean;

      style:
        | 'minimal'
        | 'digital'
        | 'large';

      timeFormat:
        | '12h'
        | '24h';

      showDate: boolean;
    };

    statistics: {
      enabled: boolean;
    };
  };
};
