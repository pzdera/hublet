<script lang="ts">
  import {
    onMount
  } from 'svelte';

  import type {
    Item
  } from '../lib/types';

  import {
    deleteLocalIcon,
    downloadDashboardLocalIcon,
    listLocalIcons,
    uploadLocalIcon
  } from '../lib/api';

  import type {
    LocalIconFile
  } from '../lib/api';

  let {
    item
  }: {
    item: Item;
  } = $props();

  let icons =
    $state<LocalIconFile[]>([]);

  let dashboardValue =
    $state('');

  let loading =
    $state(true);

  let busy =
    $state(false);

  let error =
    $state('');

  let fileInput:
    HTMLInputElement;

  onMount(() => {
    refreshIcons();
  });

  async function refreshIcons() {
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

  function selectIcon(
    icon: LocalIconFile
  ) {
    item.icon.type = 'local';
    item.icon.value = icon.filename;
  }

  async function downloadIcon() {
    const value =
      dashboardValue.trim();

    if (!value) {
      error =
        'Enter an icon name or Dashboard Icons link.';

      return;
    }

    busy = true;
    error = '';

    try {
      const icon =
        await downloadDashboardLocalIcon(
          value
        );

      icons = [
        icon,
        ...icons.filter(
          (candidate) =>
            candidate.filename !==
            icon.filename
        )
      ];

      selectIcon(icon);
      dashboardValue = '';
    } catch (reason) {
      error =
        reason instanceof Error
          ? reason.message
          : String(reason);
    } finally {
      busy = false;
    }
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

    busy = true;
    error = '';

    try {
      const icon =
        await uploadLocalIcon(
          file
        );

      icons = [
        icon,
        ...icons
      ];

      selectIcon(icon);
    } catch (reason) {
      error =
        reason instanceof Error
          ? reason.message
          : String(reason);
    } finally {
      busy = false;
      input.value = '';
    }
  }

  async function removeIcon(
    icon: LocalIconFile
  ) {
    const confirmed =
      window.confirm(
        `Delete icon "${icon.filename}"?`
      );

    if (!confirmed) {
      return;
    }

    busy = true;
    error = '';

    try {
      await deleteLocalIcon(
        icon.filename
      );

      icons =
        icons.filter(
          (candidate) =>
            candidate.filename !==
            icon.filename
        );

      if (
        item.icon.type === 'local' &&
        item.icon.value ===
          icon.filename
      ) {
        item.icon.type = 'auto';
        item.icon.value = '';
      }
    } catch (reason) {
      error =
        reason instanceof Error
          ? reason.message
          : String(reason);
    } finally {
      busy = false;
    }
  }

  function readableSize(
    bytes: number
  ): string {
    if (bytes < 1024) {
      return `${bytes} B`;
    }

    return `${
      (
        bytes / 1024
      ).toFixed(1)
    } KB`;
  }
</script>

<section class="local-icon-manager">
  <div class="local-icon-source">
    <header>
      <div>
        <strong>Dashboard Icons</strong>

        <small>
          Enter an icon name or paste its
          dashboardicons.com link.
        </small>
      </div>
    </header>

    <div class="local-icon-download">
      <input
        bind:value={dashboardValue}
        type="text"
        disabled={busy}
        placeholder="proxmox or dashboardicons.com/icons/proxmox"
        onkeydown={(event) => {
          if (event.key === 'Enter') {
            event.preventDefault();
            downloadIcon();
          }
        }}
      />

      <button
        type="button"
        disabled={
          busy ||
          !dashboardValue.trim()
        }
        onclick={downloadIcon}
      >
        {busy
          ? 'Working…'
          : 'Download'}
      </button>
    </div>
  </div>

  <div class="local-icon-divider">
    <span>or</span>
  </div>

  <div class="local-icon-upload">
    <input
      bind:this={fileInput}
      type="file"
      accept="image/png"
      onchange={uploadIcon}
    />

    <button
      type="button"
      disabled={busy}
      onclick={() =>
        fileInput.click()}
    >
      <span>+</span>

      <div>
        <strong>Upload PNG</strong>

        <small>
          Choose a PNG icon from this computer.
          Maximum size is 2 MB.
        </small>
      </div>
    </button>
  </div>

  {#if error}
    <div class="local-icon-error">
      {error}
    </div>
  {/if}

  <div class="local-icon-library-heading">
    <strong>Local library</strong>

    <small>
      {icons.length}
      {icons.length === 1
        ? ' icon'
        : ' icons'}
    </small>
  </div>

  {#if loading}
    <div class="local-icon-empty">
      Loading icons…
    </div>
  {:else if icons.length === 0}
    <div class="local-icon-empty">
      No locally stored icons yet.
    </div>
  {:else}
    <div class="local-icon-gallery">
      {#each icons as icon (icon.filename)}
        <article
          class:selected={
            item.icon.type === 'local' &&
            item.icon.value ===
              icon.filename
          }
          class="local-icon-item"
        >
          <button
            class="local-icon-select"
            type="button"
            title={icon.filename}
            onclick={() =>
              selectIcon(icon)}
          >
            <img
              src={icon.url}
              alt=""
              loading="lazy"
            />

            <span>✓</span>
          </button>

          <footer>
            <div>
              <strong>
                {icon.filename}
              </strong>

              <small>
                {readableSize(icon.size)}
              </small>
            </div>

            <button
              type="button"
              disabled={busy}
              title="Delete icon"
              onclick={() =>
                removeIcon(icon)}
            >
              ×
            </button>
          </footer>
        </article>
      {/each}
    </div>
  {/if}
</section>
