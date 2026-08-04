<script lang="ts">
  import type {
    HubletConfig,
    Shortcut
  } from '../lib/types';

  let {
    config
  }: {
    config: HubletConfig;
  } = $props();

  let key = $state('');
  let label = $state('');
  let url = $state('');
  let addError = $state('');

  const normalizedKey = $derived(
    key.trim().toLowerCase()
  );

  const keyValid = $derived(
    /^[a-z0-9_-]{1,24}$/.test(
      normalizedKey
    )
  );

  const duplicateKey = $derived(
    config.shortcuts.some(
      (shortcut) =>
        shortcut.key.trim().toLowerCase() ===
        normalizedKey
    )
  );

  const canAdd = $derived(
    keyValid &&
    !duplicateKey &&
    Boolean(label.trim()) &&
    Boolean(url.trim())
  );

  function addShortcut() {
    addError = '';

    if (!keyValid) {
      addError =
        'Use 1–24 letters, numbers, dashes, or underscores.';
      return;
    }

    if (duplicateKey) {
      addError = 'This shortcut key already exists.';
      return;
    }

    if (!label.trim() || !url.trim()) {
      addError = 'Label and URL are required.';
      return;
    }

    const shortcut: Shortcut = {
      key: normalizedKey,
      label: label.trim(),
      url: url.trim(),
      icon: {
        type: 'auto',
        value: ''
      }
    };

    config.shortcuts.push(shortcut);

    key = '';
    label = '';
    url = '';
  }

  function deleteShortcut(
    index: number
  ) {
    config.shortcuts.splice(index, 1);
  }
</script>

<div class="shortcuts-panel">
  <header class="shortcuts-header">
    <div class="shortcuts-header-icon">
      ↗
    </div>

    <div>
      <p>SEARCH</p>
      <h2>Shortcuts</h2>
      <small>
        Open a URL by entering its key in search.
      </small>
    </div>
  </header>

  <div class="shortcuts-content">
    <section class="shortcut-add-card">
      <header>
        <strong>Add shortcut</strong>
        <small>
          Example: rd opens Reddit.
        </small>
      </header>

      <div class="shortcut-key-label-row">
        <label>
          <span>Key</span>
          <input
            bind:value={key}
            type="text"
            maxlength="24"
            autocomplete="off"
            placeholder="rd"
          />
        </label>

        <label>
          <span>Label</span>
          <input
            bind:value={label}
            type="text"
            maxlength="100"
            placeholder="Reddit"
          />
        </label>
      </div>

      <label>
        <span>URL</span>
        <input
          bind:value={url}
          type="text"
          placeholder="https://reddit.com"
        />
      </label>

      {#if addError}
        <p class="shortcut-form-error">
          {addError}
        </p>
      {/if}

      <button
        class="shortcut-add-button"
        type="button"
        disabled={!canAdd}
        onclick={addShortcut}
      >
        + Add shortcut
      </button>
    </section>

    <section class="shortcut-list-section">
      <header>
        <strong>Saved shortcuts</strong>
        <span>{config.shortcuts.length}</span>
      </header>

      {#if config.shortcuts.length === 0}
        <div class="shortcut-empty-state">
          No shortcuts configured yet.
        </div>
      {:else}
        <div class="shortcut-list">
          {#each config.shortcuts as shortcut, index}
            <article class="shortcut-card">
              <div class="shortcut-card-heading">
                <span>
                  {shortcut.key || '—'}
                </span>

                <button
                  type="button"
                  title="Delete shortcut"
                  aria-label={`Delete ${shortcut.label || 'shortcut'}`}
                  onclick={() =>
                    deleteShortcut(index)}
                >
                  ×
                </button>
              </div>

              <div class="shortcut-key-label-row">
                <label>
                  <span>Key</span>
                  <input
                    bind:value={shortcut.key}
                    type="text"
                    maxlength="24"
                    autocomplete="off"
                  />
                </label>

                <label>
                  <span>Label</span>
                  <input
                    bind:value={shortcut.label}
                    type="text"
                    maxlength="100"
                  />
                </label>
              </div>

              <label>
                <span>URL</span>
                <input
                  bind:value={shortcut.url}
                  type="text"
                />
              </label>
            </article>
          {/each}
        </div>
      {/if}
    </section>

    <p class="shortcut-help">
      Type the exact key in search and press Enter.
      Shortcuts open in a new tab.
    </p>
  </div>
</div>
