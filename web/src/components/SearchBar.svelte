<script lang="ts">
  import { onMount } from 'svelte';
  import type { Item, Shortcut } from '../lib/types';
  import ServiceIcon from './ServiceIcon.svelte';

  let {
    items,
    shortcuts,
    autoFocus,
    openShortcutDirectly,
    webSearchEnabled,
    webSearchEngine,
    query = $bindable()
  }: {
    items: Item[];
    shortcuts: Shortcut[];
    autoFocus: boolean;
    openShortcutDirectly: boolean;
    webSearchEnabled: boolean;
    webSearchEngine: string;
    query: string;
  } = $props();

  let input: HTMLInputElement;

  const normalizedQuery = $derived(query.trim().toLowerCase());

  const matchingItems = $derived(
    normalizedQuery
      ? items.filter((item) => {
          const content = [
            item.name,
            item.description,
            item.url
          ]
            .join(' ')
            .toLowerCase();

          return content.includes(normalizedQuery);
        })
      : []
  );

  const exactShortcut = $derived(
    shortcuts.find(
      (shortcut) =>
        shortcut.key.toLowerCase() === normalizedQuery
    )
  );

  const exactShortcutItem = $derived.by(
    (): Item | null => {
      if (!exactShortcut) {
        return null;
      }

      return {
        id: `shortcut-${exactShortcut.key}`,
        type: 'service',
        name: exactShortcut.label,
        url: exactShortcut.url,
        description: '',
        icon: exactShortcut.icon,
        openInNewTab: true
      };
    }
  );

  function openURL(url: string, newTab = true) {
    query = '';

    if (newTab) {
      window.open(url, '_blank', 'noopener,noreferrer');
      return;
    }

    window.location.href = url;
  }

  function webSearchURL(value: string) {
    const encoded = encodeURIComponent(value);

    switch (webSearchEngine) {
      case 'duckduckgo':
        return `https://duckduckgo.com/?q=${encoded}`;
      case 'bing':
        return `https://www.bing.com/search?q=${encoded}`;
      default:
        return `https://www.google.com/search?q=${encoded}`;
    }
  }

  function submit() {
    if (!normalizedQuery) return;

    if (exactShortcut && openShortcutDirectly) {
      openURL(exactShortcut.url);
      return;
    }

    if (matchingItems.length === 1) {
      const item = matchingItems[0];
      openURL(item.url, item.openInNewTab);
      return;
    }

    if (webSearchEnabled && matchingItems.length === 0) {
      openURL(webSearchURL(query));
    }
  }

  function handleGlobalKeydown(event: KeyboardEvent) {
    const target = event.target as HTMLElement | null;
    const isTyping =
      target instanceof HTMLInputElement ||
      target instanceof HTMLTextAreaElement ||
      target?.isContentEditable;

    if (
      event.key === '/' &&
      !isTyping
    ) {
      event.preventDefault();
      input?.focus();
      return;
    }

    if (
      (event.ctrlKey || event.metaKey) &&
      event.key.toLowerCase() === 'k'
    ) {
      event.preventDefault();
      input?.focus();
      input?.select();
      return;
    }

    if (event.key === 'Escape') {
      query = '';
      input?.focus();
    }
  }

  onMount(() => {
    window.addEventListener('keydown', handleGlobalKeydown);

    if (autoFocus) {
      requestAnimationFrame(() => input?.focus());
    }

    return () => {
      window.removeEventListener(
        'keydown',
        handleGlobalKeydown
      );
    };
  });
</script>

<div class="search-shell">
  <form
    class="search-box"
    onsubmit={(event) => {
      event.preventDefault();
      submit();
    }}
  >
    <span class="search-icon" aria-hidden="true">⌕</span>

    <input
      bind:this={input}
      bind:value={query}
      type="search"
      autocomplete="off"
      spellcheck="false"
      placeholder="Search services or enter a shortcut…"
      aria-label="Search services and shortcuts"
    />

    <kbd>Ctrl K</kbd>
  </form>

  {#if normalizedQuery}
    <div class="results">
      {#if exactShortcut}
        <button
          type="button"
          class="result shortcut"
          onclick={() => openURL(exactShortcut.url)}
        >
          <span class="result-mark result-icon">
            {#if exactShortcutItem}
              <ServiceIcon
                item={exactShortcutItem}
              />
            {/if}
          </span>

          <span>
            <strong>{exactShortcut.label}</strong>
            <small>
              Shortcut · {exactShortcut.key}
            </small>
          </span>

          <span class="enter">↵</span>
        </button>
      {/if}

      {#each matchingItems.slice(0, 8) as item (item.id)}
        <button
          type="button"
          class="result"
          onclick={() =>
            openURL(item.url, item.openInNewTab)}
        >
          <span class="result-mark result-icon">
            <ServiceIcon {item} />
          </span>

          <span>
            <strong>{item.name}</strong>
            <small>
              {item.description || item.url}
            </small>
          </span>
        </button>
      {/each}

      {#if !exactShortcut &&
        matchingItems.length === 0 &&
        webSearchEnabled}
        <button
          type="button"
          class="result"
          onclick={() => openURL(webSearchURL(query))}
        >
          <span class="result-mark">↗</span>

          <span>
            <strong>Search the web</strong>
            <small>{query}</small>
          </span>

          <span class="enter">↵</span>
        </button>
      {/if}
    </div>
  {/if}
</div>
