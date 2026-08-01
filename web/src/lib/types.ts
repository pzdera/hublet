export type Icon = {
  type: 'auto' | 'local' | 'none';
  value: string;
};

export type ServiceResources = {
  enabled: boolean;
  showStatus: boolean;
  showCpu: boolean;
  showMemory: boolean;
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
  resources: ServiceResources;
};

export type SectionLayout =
  | 'list'
  | 'grid'
  | 'compact'
  | 'featured';

export type SectionWidth =
  | 'narrow'
  | 'medium'
  | 'wide'
  | 'extra-wide'
  | 'full';

export type SectionCardSize =
  | 'inherit'
  | 'small'
  | 'medium'
  | 'large';

export type Section = {
  id: string;
  title: string;
  accent: string;
  layout: SectionLayout;
  width: SectionWidth;
  cardSize: SectionCardSize;
  gridColumns: number;
  fillLastRow: boolean;
  collapsed: boolean;
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

      scale:
        | 'small'
        | 'medium'
        | 'large';
    };

    cards: {
      size:
        | 'small'
        | 'medium'
        | 'large';

      density:
        | 'compact'
        | 'comfortable'
        | 'relaxed';

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
        | 'gradient'
        | 'wallpaper';

      color: string;
      blur: number;
      brightness: number;
      overlay: number;
      wallpaper: string | null;
    };

    animations: boolean;
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
