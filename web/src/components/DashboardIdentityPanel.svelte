<script lang="ts">
  import {
    onMount
  } from 'svelte';

  import type {
    HubletConfig
  } from '../lib/types';

  import {
    listLocalIcons,
    uploadLocalIcon
  } from '../lib/api';

  import type {
    LocalIconFile
  } from '../lib/api';

  import DashboardBrand from './DashboardBrand.svelte';

  let {
    config
  }: {
    config: HubletConfig;
  } = $props();

  let icons =
    $state<LocalIconFile[]>([]);

  let loading =
    $state(false);

  let uploading =
    $state(false);

  let error =
    $state('');

  let fileInput:
    HTMLInputElement;

  onMount(() => {
    loadIcons();
  });

  async function loadIcons() {
    loading = true;
    error = '';

    try {
      icons =
        await listLocalIcons();
    } catch (reason) {
      error =
        reason instanceof Error
          ? reason.message
          : String(reason);
    } finally {
      loading = false;
    }
  }

  function setIconType(
    type:
      HubletConfig['dashboard']['icon']['type']
  ) {
    config.dashboard.icon.type =
      type;

    if (
      type === 'initial' &&
      !config.dashboard.icon.value.trim()
    ) {
      config.dashboard.icon.value =
        config.dashboard.title
          .trim()
          .slice(0, 1)
          .toUpperCase() ||
        'H';
    }
  }

  function selectLocalIcon(
    icon: LocalIconFile
  ) {
    config.dashboard.icon.type =
      'local';

    config.dashboard.icon.value =
      icon.filename;
  }

  async function uploadIcon(
    event: Event
  ) {
    const input =
      event.currentTarget as
        HTMLInputElement;

    const file =
      input.files?.[0];

    if (!file) {
      return;
    }

    uploading = true;
    error = '';

    try {
      const icon =
        await uploadLocalIcon(file);

      icons = [
        icon,
        ...icons.filter(
          (candidate) =>
            candidate.filename !==
            icon.filename
        )
      ];

      selectLocalIcon(icon);
    } catch (reason) {
      error =
        reason instanceof Error
          ? reason.message
          : String(reason);
    } finally {
      uploading = false;
      input.value = '';
    }
  }
</script>

<section class="dashboard-identity-panel">
  <header class="dashboard-identity-heading">
    <div>
      <strong>Dashboard identity</strong>

      <small>
        Customize the title, description, and header icon.
      </small>
    </div>
  </header>

  <div class="dashboard-identity-preview">
    <DashboardBrand
      dashboard={config.dashboard}
      preview={true}
    />
  </div>

  <label class="appearance-field">
    <span>Dashboard title</span>

    <input
      bind:value={config.dashboard.title}
      type="text"
      maxlength="80"
      placeholder="Hublet v2"
    />
  </label>

  <label class="appearance-field">
    <span>Description</span>

    <input
      bind:value={
        config.dashboard.description
      }
      type="text"
      maxlength="160"
      placeholder="My self-hosted dashboard"
    />
  </label>

  <label class="appearance-toggle">
    <span>
      <strong>Show description</strong>

      <small>
        Display the description above the dashboard title.
      </small>
    </span>

    <input
      bind:checked={
        config.dashboard.descriptionVisible
      }
      type="checkbox"
    />
  </label>

  <div class="appearance-label">
    <span>Header icon</span>

    <div class="appearance-segmented">
      <button
        class:active={
          config.dashboard.icon.type ===
          'initial'
        }
        type="button"
        onclick={() =>
          setIconType('initial')}
      >
        Initial
      </button>

      <button
        class:active={
          config.dashboard.icon.type ===
          'local'
        }
        type="button"
        onclick={() =>
          setIconType('local')}
      >
        Local
      </button>

      <button
        class:active={
          config.dashboard.icon.type ===
          'none'
        }
        type="button"
        onclick={() =>
          setIconType('none')}
      >
        None
      </button>
    </div>
  </div>

  {#if config.dashboard.icon.type === 'initial'}
    <label class="appearance-field">
      <span>Initial</span>

      <input
        bind:value={
          config.dashboard.icon.value
        }
        type="text"
        maxlength="2"
        placeholder="H"
      />

      <small>
        One or two characters.
      </small>
    </label>
  {/if}

  {#if config.dashboard.icon.type === 'local'}
    <div class="dashboard-icon-library">
      <div class="dashboard-icon-library-header">
        <div>
          <strong>Local icon library</strong>

          <small>
            Select an existing PNG or upload a new one.
          </small>
        </div>

        <button
          type="button"
          disabled={uploading}
          onclick={() =>
            fileInput.click()}
        >
          {uploading
            ? 'Uploading…'
            : 'Upload PNG'}
        </button>
      </div>

      <input
        class="dashboard-icon-file-input"
        bind:this={fileInput}
        type="file"
        accept="image/png"
        onchange={uploadIcon}
      />

      {#if error}
        <div class="dashboard-icon-error">
          {error}
        </div>
      {/if}

      {#if loading}
        <div class="dashboard-icon-empty">
          Loading icons…
        </div>
      {:else if icons.length === 0}
        <button
          class="dashboard-icon-empty actionable"
          type="button"
          onclick={() =>
            fileInput.click()}
        >
          Upload the first PNG icon
        </button>
      {:else}
        <div class="dashboard-icon-grid">
          {#each icons as icon (icon.filename)}
            <button
              class:active={
                config.dashboard.icon.value ===
                icon.filename
              }
              type="button"
              title={icon.filename}
              onclick={() =>
                selectLocalIcon(icon)}
            >
              <img
                src={icon.url}
                alt=""
                loading="lazy"
              />

              <span>✓</span>
            </button>
          {/each}
        </div>
      {/if}
    </div>
  {/if}

  {#if config.dashboard.icon.type !== 'none'}
    <div class="appearance-label">
      <span>Icon size</span>

      <div class="appearance-segmented">
        {#each [
          'small',
          'medium',
          'large'
        ] as size}
          <button
            class:active={
              config.dashboard.iconSize ===
              size
            }
            type="button"
            onclick={() => {
              config.dashboard.iconSize =
                size as
                  HubletConfig['dashboard']['iconSize'];
            }}
          >
            {size}
          </button>
        {/each}
      </div>
    </div>
  {/if}
</section>
