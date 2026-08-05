<script lang="ts">
  import {
    onMount,
    tick
  } from 'svelte';

  import type {
    Item,
    Shortcut
  } from '../lib/types';

  import ServiceIcon from './ServiceIcon.svelte';

  const historyStorageKey =
    'hublet.commandHistory.v1';

  const historyLimit = 30;

  let {
    items,
    shortcuts,
    autoFocus,
    openShortcutDirectly,
    webSearchEnabled,
    webSearchEngine
  }: {
    items: Item[];
    shortcuts: Shortcut[];
    autoFocus: boolean;
    openShortcutDirectly: boolean;
    webSearchEnabled: boolean;
    webSearchEngine: string;
  } = $props();

  let input: HTMLInputElement;
  let query = $state('');
  let history = $state<string[]>([]);
  let historyIndex = $state(-1);
  let historyDraft = $state('');

  const normalizedQuery = $derived(
    query.trim().toLowerCase()
  );

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

  const exactItem = $derived(
    items.find(
      (item) =>
        item.name.trim().toLowerCase() ===
          normalizedQuery ||
        item.url.trim().toLowerCase() ===
          normalizedQuery
    )
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
        icon: exactShortcut.icon
      };
    }
  );

  const historyCompletion = $derived.by(
    (): string => {
      if (!normalizedQuery) {
        return '';
      }

      return history.find((command) => {
        const normalizedCommand =
          command.toLowerCase();

        return normalizedCommand.startsWith(
          normalizedQuery
        ) && normalizedCommand !== normalizedQuery;
      }) ?? '';
    }
  );

  const directQueryURL = $derived(
    normalizedQuery
      ? directURL(query)
      : null
  );

  function loadHistory(): void {
    try {
      const stored = JSON.parse(
        localStorage.getItem(historyStorageKey) ?? '[]'
      );

      if (!Array.isArray(stored)) {
        return;
      }

      history = stored
        .filter(
          (value): value is string =>
            typeof value === 'string' &&
            value.trim() !== ''
        )
        .slice(0, historyLimit);
    } catch {
      history = [];
    }
  }

  function saveHistory(): void {
    try {
      localStorage.setItem(
        historyStorageKey,
        JSON.stringify(history)
      );
    } catch {
      // Search remains usable when browser storage is blocked.
    }
  }

  function rememberCommand(value: string): void {
    const command = value.trim();

    if (!command) {
      return;
    }

    history = [
      command,
      ...history.filter(
        (entry) =>
          entry.toLowerCase() !== command.toLowerCase()
      )
    ].slice(0, historyLimit);

    historyIndex = -1;
    historyDraft = '';
    saveHistory();
  }

  function openURL(
    url: string,
    command = query
  ): void {
    rememberCommand(command);
    query = '';

    window.open(
      url,
      '_blank',
      'noopener,noreferrer'
    );
  }

  function webSearchURL(value: string): string {
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

  function directURL(value: string): string | null {
    const trimmed = value.trim();

    if (/^https?:\/\//i.test(trimmed)) {
      try {
        return new URL(trimmed).toString();
      } catch {
        return null;
      }
    }

    const candidate = trimmed.split('/')[0];
    const host = candidate.split(':')[0].toLowerCase();

    const looksLikeHost =
      host === 'localhost' ||
      /^\d{1,3}(\.\d{1,3}){3}$/.test(host) ||
      /^[a-z0-9-]+(\.[a-z0-9-]+)+$/i.test(host);

    if (!looksLikeHost || /\s/.test(trimmed)) {
      return null;
    }

    const localHost =
      host === 'localhost' ||
      host === '127.0.0.1' ||
      /^10\./.test(host) ||
      /^192\.168\./.test(host) ||
      /^172\.(1[6-9]|2[0-9]|3[0-1])\./.test(host) ||
      host.endsWith('.local') ||
      host.endsWith('.lan') ||
      host.endsWith('.home') ||
      host.endsWith('.internal');

    return `${localHost ? 'http' : 'https'}://${trimmed}`;
  }

  function submitCommand(value = query): void {
    const command = value.trim();

    if (!command) {
      return;
    }

    const normalized = command.toLowerCase();

    const shortcut = shortcuts.find(
      (entry) => entry.key.toLowerCase() === normalized
    );

    if (shortcut && openShortcutDirectly) {
      openURL(shortcut.url, command);
      return;
    }

    const item = items.find(
      (entry) =>
        entry.name.trim().toLowerCase() === normalized ||
        entry.url.trim().toLowerCase() === normalized
    );

    if (item) {
      openURL(item.url, command);
      return;
    }

    const direct = directURL(command);

    if (direct) {
      openURL(direct, command);
      return;
    }

    const matching = items.filter((entry) => {
      return [
        entry.name,
        entry.description,
        entry.url
      ]
        .join(' ')
        .toLowerCase()
        .includes(normalized);
    });

    if (matching.length === 1) {
      openURL(matching[0].url, command);
      return;
    }

    if (webSearchEnabled) {
      openURL(webSearchURL(command), command);
    }
  }

  async function placeCursorAtEnd(): Promise<void> {
    await tick();
    const end = query.length;
    input?.setSelectionRange(end, end);
  }

  function previousCommand(): void {
    if (history.length === 0) {
      return;
    }

    if (historyIndex === -1) {
      historyDraft = query;
    }

    historyIndex = Math.min(
      historyIndex + 1,
      history.length - 1
    );

    query = history[historyIndex];
    void placeCursorAtEnd();
  }

  function nextCommand(): void {
    if (historyIndex === -1) {
      return;
    }

    if (historyIndex === 0) {
      historyIndex = -1;
      query = historyDraft;
    } else {
      historyIndex -= 1;
      query = history[historyIndex];
    }

    void placeCursorAtEnd();
  }

  function completeFromHistory(): void {
    if (!historyCompletion) {
      return;
    }

    query = historyCompletion;
    historyIndex = -1;
    historyDraft = '';
    void placeCursorAtEnd();
  }

  function handleInputKeydown(event: KeyboardEvent): void {
    switch (event.key) {
      case 'ArrowUp':
        event.preventDefault();
        previousCommand();
        break;
      case 'ArrowDown':
        event.preventDefault();
        nextCommand();
        break;
      case 'Tab':
        if (historyCompletion) {
          event.preventDefault();
          completeFromHistory();
        }
        break;
      case 'Escape':
        event.preventDefault();
        query = '';
        historyIndex = -1;
        historyDraft = '';
        break;
    }
  }

  function handleGlobalKeydown(event: KeyboardEvent): void {
    const target = event.target as HTMLElement | null;
    const isTyping =
      target instanceof HTMLInputElement ||
      target instanceof HTMLTextAreaElement ||
      target?.isContentEditable;

    if (event.key === '/' && !isTyping) {
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

    if (event.key === 'Escape' && !isTyping) {
      query = '';
    }
  }

  onMount(() => {
    loadHistory();
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

<div class="search-shell command-search">
  <form
    class="search-box"
    onsubmit={(event) => {
      event.preventDefault();
      submitCommand();
    }}
  >
    <span class="search-icon" aria-hidden="true">›_</span>

    <input
      bind:this={input}
      bind:value={query}
      type="search"
      autocomplete="off"
      spellcheck="false"
      placeholder="Type a service, shortcut, URL or search…"
      aria-label="Hublet command search"
      oninput={() => {
        historyIndex = -1;
        historyDraft = '';
      }}
      onkeydown={handleInputKeydown}
    />

    {#if historyCompletion}
      <kbd title="Complete from command history">Tab ↹</kbd>
    {:else}
      <kbd>Ctrl K</kbd>
    {/if}
  </form>

  {#if normalizedQuery}
    <div class="results command-results">
      {#if historyCompletion}
        <button
          type="button"
          class="result history-result"
          onclick={completeFromHistory}
        >
          <span class="result-mark">↥</span>

          <span>
            <strong>{historyCompletion}</strong>
            <small>Recent command</small>
          </span>

          <span class="enter">Tab</span>
        </button>
      {/if}

      {#if exactShortcut}
        <button
          type="button"
          class="result shortcut"
          onclick={() => openURL(
            exactShortcut.url,
            query
          )}
        >
          <span class="result-mark result-icon">
            {#if exactShortcutItem}
              <ServiceIcon item={exactShortcutItem} />
            {/if}
          </span>

          <span>
            <strong>{exactShortcut.label}</strong>
            <small>Shortcut · {exactShortcut.key}</small>
          </span>

          <span class="enter">↵</span>
        </button>
      {/if}

      {#each matchingItems.slice(0, 8) as item (item.id)}
        <button
          type="button"
          class="result"
          class:exact-result={item.id === exactItem?.id}
          onclick={() => openURL(item.url, query)}
        >
          <span class="result-mark result-icon">
            <ServiceIcon {item} />
          </span>

          <span>
            <strong>{item.name}</strong>
            <small>{item.description || item.url}</small>
          </span>

          {#if item.id === exactItem?.id}
            <span class="enter">↵</span>
          {/if}
        </button>
      {/each}

      {#if directQueryURL && matchingItems.length === 0}
        <button
          type="button"
          class="result exact-result"
          onclick={() => openURL(directQueryURL, query)}
        >
          <span class="result-mark">↗</span>

          <span>
            <strong>Open address</strong>
            <small>{directQueryURL}</small>
          </span>

          <span class="enter">↵</span>
        </button>
      {:else if matchingItems.length === 0 && webSearchEnabled}
        <button
          type="button"
          class="result"
          onclick={() => openURL(
            webSearchURL(query),
            query
          )}
        >
          <span class="result-mark">⌕</span>

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
