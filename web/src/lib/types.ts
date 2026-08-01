export type Icon = {
  type: 'auto' | 'local' | 'none';
  value: string;
};

export type Item = {
  id: string;
  name: string;
  url: string;
  description: string;
  icon: Icon;
  openInNewTab: boolean;
};

export type Section = {
  id: string;
  title: string;
  accent: string;
  layout: 'grid' | 'list' | 'compact' | 'large';
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
    density: string;
    cardStyle: string;
    radius: string;
  };

  search: {
    autoFocus: boolean;
    openShortcutDirectly: boolean;
    webSearchEnabled: boolean;
    webSearchEngine: 'google' | 'duckduckgo' | 'bing';
  };

  sections: Section[];
  shortcuts: Shortcut[];

  weather: {
    enabled: boolean;
    location: string;
    latitude: number | null;
    longitude: number | null;
  };
};
