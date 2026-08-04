<script lang="ts">
  import {
    onDestroy
  } from 'svelte';

  import type {
    HubletConfig
  } from '../lib/types';

  import {
    applyAppearanceVariables,
    captureAppearanceVariables,
    restoreAppearanceVariables
  } from '../lib/appearance';

  import WallpaperManager from './WallpaperManager.svelte';

  let {
    config
  }: {
    config: HubletConfig;
  } = $props();

  const root =
    document.documentElement;

  const originalVariables =
    captureAppearanceVariables(root);

  const fonts: Array<{
    value:
      HubletConfig['appearance']['font']['family'];
    label: string;
    sample: string;
  }> = [
    {
      value: 'system',
      label: 'System',
      sample: 'Fast and familiar'
    },
    {
      value: 'inter',
      label: 'Inter',
      sample: 'Clean and neutral'
    },
    {
      value: 'geist',
      label: 'Geist',
      sample: 'Modern interface'
    },
    {
      value: 'manrope',
      label: 'Manrope',
      sample: 'Friendly geometry'
    },
    {
      value: 'ibm-plex-sans',
      label: 'IBM Plex',
      sample: 'Technical character'
    }
  ];

  const presets: Array<{
    id: HubletConfig['appearance']['mode'];
    name: string;
    description: string;
  }> = [
    {
      id: 'minimal',
      name: 'Minimal',
      description:
        'Subtle cards and almost no visual noise.'
    },
    {
      id: 'standard',
      name: 'Standard',
      description:
        'Balanced depth, spacing, and contrast.'
    },
    {
      id: 'advanced',
      name: 'Glass',
      description:
        'Transparent cards with backdrop blur.'
    }
  ];

  $effect(() => {
    applyAppearanceVariables(
      root,
      config
    );
  });

  onDestroy(() => {
    restoreAppearanceVariables(
      root,
      originalVariables
    );
  });

  function applyPreset(
    preset:
      HubletConfig['appearance']['mode']
  ) {
    config.appearance.mode = preset;

    switch (preset) {
      case 'minimal':
        config.appearance.font.family =
          'system';

        config.appearance.cards.radius =
          'medium';

        config.appearance.cards.shadow =
          'none';

        config.appearance.cards.border =
          true;

        break;

      case 'advanced':
        config.appearance.font.family =
          'manrope';

        config.appearance.cards.radius =
          'large';

        config.appearance.cards.shadow =
          'floating';

        config.appearance.cards.border =
          true;

        if (
          config.appearance.background.type ===
            'wallpaper' &&
          config.appearance.background.overlay <
            18
        ) {
          config.appearance.background.overlay =
            28;
        }

        break;

      case 'standard':
      default:
        config.appearance.font.family =
          'inter';

        config.appearance.cards.radius =
          'large';

        config.appearance.cards.shadow =
          'soft';

        config.appearance.cards.border =
          true;

        break;
    }
  }

  function selectBackgroundType(
    type:
      HubletConfig['appearance']['background']['type']
  ) {
    config.appearance.background.type =
      type;

    if (
      type === 'wallpaper' &&
      !config.appearance.background.wallpaper
    ) {
      config.appearance.background.wallpaper =
        '';
    }
  }

  import DashboardIdentityPanel from './DashboardIdentityPanel.svelte';
</script>

<aside class="appearance-panel">
  <header class="appearance-header">
    <div class="appearance-header-icon">
      ◉
    </div>

    <div>
      <p>APPEARANCE</p>
      <h2>Theme Studio</h2>
    </div>
  </header>

  <div class="appearance-content">
    <DashboardIdentityPanel {config} />
    <section class="appearance-group">
      <header>
        <div>
          <strong>Style preset</strong>

          <small>
            Start with a complete visual style,
            then customize every setting.
          </small>
        </div>
      </header>

      <div class="theme-presets">
        {#each presets as preset}
          <button
            class:active={
              config.appearance.mode ===
              preset.id
            }
            type="button"
            onclick={() =>
              applyPreset(preset.id)}
          >
            <span
              class={`theme-preset-preview ${preset.id}`}
            >
              <i></i>
              <i></i>
              <i></i>
            </span>

            <span>
              <strong>{preset.name}</strong>
              <small>{preset.description}</small>
            </span>
          </button>
        {/each}
      </div>
    </section>

    <section class="appearance-group">
      <header>
        <div>
          <strong>Background</strong>

          <small>
            Choose a solid color or wallpaper.
          </small>
        </div>
      </header>

      <div class="background-type-options">
        <button
          class:active={
            config.appearance.background.type ===
            'solid'
          }
          type="button"
          onclick={() =>
            selectBackgroundType('solid')}
        >
          <span class="background-preview solid">
          </span>
          <strong>Solid</strong>
        </button>

        <button
          class:active={
            config.appearance.background.type ===
            'wallpaper'
          }
          type="button"
          onclick={() =>
            selectBackgroundType('wallpaper')}
        >
          <span class="background-preview wallpaper">
            ▧
          </span>
          <strong>Wallpaper</strong>
        </button>
      </div>

      <label class="appearance-field">
        <span>Base color</span>

        <div class="appearance-color-control">
          <input
            bind:value={
              config.appearance.background.color
            }
            type="color"
          />

          <input
            bind:value={
              config.appearance.background.color
            }
            type="text"
            maxlength="7"
            placeholder="#090c12"
          />
        </div>
      </label>

      {#if (
        config.appearance.background.type ===
        'wallpaper'
      )}
        <WallpaperManager {config} />

        <label class="appearance-range">
          <span>
            <strong>Brightness</strong>
            <output>
              {config.appearance.background.brightness}%
            </output>
          </span>

          <input
            bind:value={
              config.appearance.background.brightness
            }
            type="range"
            min="20"
            max="150"
            step="1"
          />
        </label>

        <label class="appearance-range">
          <span>
            <strong>Dark overlay</strong>
            <output>
              {config.appearance.background.overlay}%
            </output>
          </span>

          <input
            bind:value={
              config.appearance.background.overlay
            }
            type="range"
            min="0"
            max="100"
            step="1"
          />
        </label>

        <label class="appearance-range">
          <span>
            <strong>Background blur</strong>
            <output>
              {config.appearance.background.blur}px
            </output>
          </span>

          <input
            bind:value={
              config.appearance.background.blur
            }
            type="range"
            min="0"
            max="40"
            step="1"
          />
        </label>
      {/if}
    </section>

    <section class="appearance-group">
      <header>
        <div>
          <strong>Typography</strong>

          <small>
            Font settings apply to the complete
            dashboard.
          </small>
        </div>
      </header>

      <div class="font-options">
        {#each fonts as font}
          <button
            class:active={
              config.appearance.font.family ===
              font.value
            }
            class={`font-${font.value}`}
            type="button"
            onclick={() => {
              config.appearance.font.family =
                font.value;
            }}
          >
            <span>Aa</span>

            <span>
              <strong>{font.label}</strong>
              <small>{font.sample}</small>
            </span>
          </button>
        {/each}
      </div>

    </section>

    <section class="appearance-group">
      <header>
        <div>
          <strong>Cards</strong>

          <small>
            Shared styling for every service card.
          </small>
        </div>
      </header>

      <label class="appearance-label">
        <span>Corner radius</span>

        <div class="radius-options">
          {#each ['small', 'medium', 'large'] as radius}
            <button
              class:active={
                config.appearance.cards.radius ===
                radius
              }
              type="button"
              onclick={() => {
                config.appearance.cards.radius =
                  radius as
                    HubletConfig['appearance']['cards']['radius'];
              }}
            >
              <span class={`radius-preview ${radius}`}>
              </span>
              <strong>{radius}</strong>
            </button>
          {/each}
        </div>
      </label>

      <label class="appearance-label">
        <span>Shadow</span>

        <div class="shadow-options">
          {#each [
            'none',
            'soft',
            'medium',
            'floating'
          ] as shadow}
            <button
              class:active={
                config.appearance.cards.shadow ===
                shadow
              }
              type="button"
              onclick={() => {
                config.appearance.cards.shadow =
                  shadow as
                    HubletConfig['appearance']['cards']['shadow'];
              }}
            >
              <span class={`shadow-preview ${shadow}`}>
              </span>
              <strong>{shadow}</strong>
            </button>
          {/each}
        </div>
      </label>

      <label class="appearance-toggle">
        <span>
          <strong>Card borders</strong>

          <small>
            Draw a subtle outline around cards.
          </small>
        </span>

        <input
          bind:checked={
            config.appearance.cards.border
          }
          type="checkbox"
        />
      </label>
    </section>

  </div>
</aside>
