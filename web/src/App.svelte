<script lang="ts">
  import { loadConfig } from './lib/api';
  import type { HubletConfig, Item } from './lib/types';
  import SearchBar from './components/SearchBar.svelte';

  let config = $state<HubletConfig | null>(null);
  let loading = $state(true);
  let error = $state('');
  let query = $state('');

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
      ? config.sections.flatMap((section) => section.items)
      : []
  );

  const normalizedQuery = $derived(
    query.trim().toLowerCase()
  );

  function sectionItems(items: Item[]) {
    if (!normalizedQuery) return items;

    return items.filter((item) => {
      return [
        item.name,
        item.description,
        item.url
      ]
        .join(' ')
        .toLowerCase()
        .includes(normalizedQuery);
    });
  }
</script>

{#if loading}
  <main class="state">Loading Hublet…</main>
{:else if error || !config}
  <main class="state error">
    {error || 'Hublet configuration is unavailable.'}
  </main>
{:else}
  <div class="page">
    <header class="header">
      <div class="brand">
        <span class="brand-mark">H</span>

        <div>
          <p>SELF-HOSTED DASHBOARD</p>
          <h1>{config.dashboard.title}</h1>
        </div>
      </div>

      <button class="edit-button" type="button">
        Edit
      </button>
    </header>

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
      bind:query
    />

    <main class="sections">
      {#each config.sections as section (section.id)}
        {@const visibleItems = sectionItems(section.items)}

        {#if visibleItems.length > 0 || !normalizedQuery}
          <section
            class="section"
            style={`--section-accent:${section.accent}`}
          >
            <header class="section-header">
              <div>
                <span class="accent-dot"></span>
                <h2>{section.title}</h2>
              </div>

              <span>{visibleItems.length}</span>
            </header>

            <div class={`cards ${section.layout}`}>
              {#each visibleItems as item (item.id)}
                <a
                  class="card"
                  href={item.url}
                  target={
                    item.openInNewTab ? '_blank' : '_self'
                  }
                  rel="noreferrer"
                >
                  <span class="card-icon">
                    {item.name.slice(0, 1).toUpperCase()}
                  </span>

                  <span class="card-copy">
                    <strong>{item.name}</strong>

                    {#if item.description}
                      <small>{item.description}</small>
                    {/if}
                  </span>

                  <span class="card-arrow">↗</span>
                </a>
              {/each}
            </div>
          </section>
        {/if}
      {/each}
    </main>
  </div>
{/if}
