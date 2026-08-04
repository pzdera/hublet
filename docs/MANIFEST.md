# Hublet Manifest

## Simple by default. Powerful when needed.

Hublet is a modern self-hosted dashboard designed around clarity, speed, and graphical configuration.

## Principles

### GUI first

Normal configuration must be possible through the graphical interface.

Users should not need to edit JSON, YAML, source code, or Docker configuration to manage their dashboard.

### Minimal by default

A fresh Hublet installation should contain only:

- Dashboard title
- Search
- Empty content area
- Edit control

Widgets, wallpaper, statistics, weather, status indicators, and decorative effects are optional.

### Everything optional

Additional features must be opt-in.

Disabled features must not consume visible space or distract the user.

### Progressive disclosure

Basic settings should show only commonly used controls.

Advanced settings appear only when Advanced mode is enabled or when a specific feature is selected.

### Responsive from day one

Desktop, tablet, and mobile are equal interfaces.

Mobile must not be a compressed desktop layout.

### Fast over flashy

Animations should be short, subtle, and purposeful.

Performance, clarity, and accessibility take priority over visual effects.

### Beautiful defaults

The default theme, typography, spacing, and cards should look polished without customization.

### Safe editing

Dashboard changes are made in a draft and applied only when the user chooses Save.

Cancel must discard the entire draft.

### Predictable persistence

Container updates must never replace user configuration, icons, wallpapers, or uploads.

### One obvious way

Each task should have one clear workflow.

The application should avoid multiple overlapping systems that solve the same problem.

## Product boundaries

Hublet is a dashboard and launcher.

It is not intended to become:

- A monitoring platform
- A container manager
- A reverse proxy
- A secrets manager
- A low-code application builder
- A replacement for Grafana
- An unrestricted plugin platform
