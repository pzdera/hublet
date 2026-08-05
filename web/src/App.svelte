<script lang="ts">
  import {
    loadConfig,
    saveConfig
  } from './lib/api';

  import {
    appearanceStyle
  } from './lib/appearance';

  import type {
    HubletConfig,
    Section
  } from './lib/types';

  import {
    sectionSurfaceStyle
  } from './lib/section-surface';

  import {
    sectionPlacementStyle
  } from './lib/section-placement';

  import {
    measureSection
  } from './lib/section-measure';

  import SearchBar from './components/SearchBar.svelte';
  import EditorPanel from './components/EditorPanel.svelte';
  import ServiceIcon from './components/ServiceIcon.svelte';
  import DashboardBrand from './components/DashboardBrand.svelte';
  import WeatherWidget from './components/WeatherWidget.svelte';

  let config =
    $state<HubletConfig | null>(null);

  let draft =
    $state<HubletConfig | null>(null);

  let loading = $state(true);
  let saving = $state(false);
  let editing = $state(false);

  let error = $state('');
  let editorError = $state('');

  loadConfig()
    .then((value) => {
      config = value;
    })
    .catch((reason) => {
      error =
        reason instanceof Error
          ? reason.message
          : String(reason);
    })
    .finally(() => {
      loading = false;
    });

  const allItems = $derived(
    config
      ? config.sections.flatMap(
          (section) =>
            section.items
        )
      : []
  );

  function cloneConfig(
    value: HubletConfig
  ): HubletConfig {
    return JSON.parse(
      JSON.stringify(
        $state.snapshot(value)
      )
    ) as HubletConfig;
  }

  function normalizeURL(
    value: string
  ): string {
    const trimmed =
      value.trim();

    if (!trimmed) {
      return trimmed;
    }

    if (
      trimmed.startsWith('http://') ||
      trimmed.startsWith('https://')
    ) {
      return trimmed;
    }

    const hostPart = trimmed
      .split('/')[0]
      .split(':')[0]
      .toLowerCase();

    const isLocalhost =
      hostPart === 'localhost' ||
      hostPart === '127.0.0.1' ||
      hostPart === '::1';

    const isLocalDomain =
      hostPart.endsWith('.local') ||
      hostPart.endsWith('.lan') ||
      hostPart.endsWith('.home') ||
      hostPart.endsWith('.internal');

    const isPrivateIPv4 =
      /^10\./.test(hostPart) ||
      /^192\.168\./.test(hostPart) ||
      /^172\.(1[6-9]|2[0-9]|3[0-1])\./
        .test(hostPart);

    if (
      isLocalhost ||
      isLocalDomain ||
      isPrivateIPv4
    ) {
      return `http://${trimmed}`;
    }

    return `https://${trimmed}`;
  }

  function normalizeConfigURLs(
    value: HubletConfig
  ): HubletConfig {
    const normalized =
      cloneConfig(value);

    for (
      const section
      of normalized.sections
    ) {
      for (
        const item
        of section.items
      ) {
        item.url =
          normalizeURL(item.url);
      }
    }

    for (
      const shortcut
      of normalized.shortcuts
    ) {
      shortcut.url =
        normalizeURL(
          shortcut.url
        );
    }

    return normalized;
  }

  function startEditing() {
    if (!config) {
      return;
    }

    draft =
      cloneConfig(config);

    editorError = '';
    editing = true;

    document.body
      .classList
      .add('editor-open');
  }

  function cancelEditing() {
    draft = null;
    editorError = '';
    editing = false;

    document.body
      .classList
      .remove('editor-open');
  }

  async function commitEditing() {
    if (!draft) {
      return;
    }

    saving = true;
    editorError = '';

    try {
      const cleanDraft =
        normalizeConfigURLs(
          draft
        );

      await saveConfig(
        cleanDraft
      );

      config = cleanDraft;
      draft = null;
      editing = false;

      document.body
        .classList
        .remove('editor-open');
    } catch (reason) {
      editorError =
        reason instanceof Error
          ? reason.message
          : String(reason);
    } finally {
      saving = false;
    }
  }
</script>

{#if loading}
  <main class="state">
    Loading Hublet…
  </main>
{:else if error || !config}
  <main class="state error">
    {error ||
      'Hublet configuration is unavailable.'}
  </main>
{:else}
  <div
    class="dashboard-root"
    style={appearanceStyle(config)}
  >
    <div
      class="dashboard-background"
      aria-hidden="true"
    ></div>

    <div
      class="dashboard-overlay"
      aria-hidden="true"
    ></div>

    <div class="page">
      <header class="header">
        <DashboardBrand
          dashboard={config.dashboard}
        />

        <SearchBar
          items={allItems}
          shortcuts={config.shortcuts}
          autoFocus={config.search.autoFocus}
          openShortcutDirectly={
            config.search.openShortcutDirectly
          }
          webSearchEnabled={
            config.search.webSearchEnabled
          }
          webSearchEngine={
            config.search.webSearchEngine
          }
        />

        <div class="header-actions">
          {#if config.modules.weather.enabled}
            <WeatherWidget />
          {/if}

          <button
            class="edit-button"
            type="button"
            title="Open editor"
            aria-label="Open editor"
            onclick={startEditing}
          >
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
            <span aria-hidden="true"></span>
          </button>
        </div>
      </header>

      <main class="sections">
        {#if config.sections.length === 0}
          <section class="empty-dashboard">
            <div class="empty-dashboard-icon">
              +
            </div>

            <h2>
              Your dashboard is empty
            </h2>

            <p>
              Add a section and start organizing
              your services and favorite websites.
            </p>

            <button
              class="save-button"
              type="button"
              onclick={startEditing}
            >
              Build your dashboard
            </button>
          </section>
        {:else}
          {#each config.sections as section (section.id)}
            <section
                class={[
                  'section',
                  `width-${section.width}`,
                  `surface-${section.surface}`
                ].join(' ')}
                style={[
                  sectionSurfaceStyle(section),
                  sectionPlacementStyle(section)
                ].filter(Boolean).join(';')}
                use:measureSection={{
                  section,
                  sections: config.sections,
                  desktopOnly: true
                }}
              >
                <header class="section-header">
                  <div>
                    <span class="accent-dot"></span>
                    <h2>{section.title}</h2>
                  </div>

                  <span>{section.items.length}</span>
                </header>

                <div
                  class="cards"
                  style={`--grid-columns:${section.gridColumns}`}
                >
                    {#each section.items as item (item.id)}
                      <a
                        class="card"
                        href={item.url}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        <span class="card-icon">
                          <ServiceIcon {item} />
                        </span>

                        <span class="card-copy">
                          <strong>{item.name}</strong>

                          {#if item.description}
                            <small>
                              {item.description}
                            </small>
                          {/if}
                        </span>

                      </a>
                    {/each}

                    {#if section.items.length === 0}
                      <div class="empty-section">
                        No services in this section.
                      </div>
                    {/if}
                </div>
            </section>
          {/each}
        {/if}
      </main>
    </div>
  </div>

  {#if editing && draft}
    <EditorPanel
      config={draft}
      saving={saving}
      error={editorError}
      onSave={commitEditing}
      onCancel={cancelEditing}
    />
  {/if}
{/if}
