<script lang="ts">
  import {
    onMount
  } from 'svelte';

  import type {
    HubletConfig
  } from '../lib/types';

  import {
    deleteWallpaper,
    downloadWallpaper,
    listWallpapers,
    uploadWallpaper
  } from '../lib/api';

  import type {
    WallpaperFile
  } from '../lib/api';

  let {
    config
  }: {
    config: HubletConfig;
  } = $props();

  let wallpapers =
    $state<WallpaperFile[]>([]);

  let remoteURL =
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
    refresh();
  });

  async function refresh() {
    loading = true;
    error = '';

    try {
      wallpapers =
        await listWallpapers();
    } catch (reason) {
      error =
        reason instanceof Error
          ? reason.message
          : String(reason);
    } finally {
      loading = false;
    }
  }

  function selectWallpaper(
    wallpaper: WallpaperFile
  ) {
    config.appearance.background.type =
      'wallpaper';

    config.appearance.background.wallpaper =
      wallpaper.filename;

    if (
      config.appearance.background.overlay <
      10
    ) {
      config.appearance.background.overlay =
        24;
    }
  }

  async function handleFileSelection(
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
      const wallpaper =
        await uploadWallpaper(
          file
        );

      wallpapers = [
        wallpaper,
        ...wallpapers.filter(
          (candidate) =>
            candidate.filename !==
            wallpaper.filename
        )
      ];

      selectWallpaper(
        wallpaper
      );
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

  async function handleURLDownload() {
    const value =
      remoteURL.trim();

    if (!value) {
      error =
        'Paste a direct image URL first.';

      return;
    }

    busy = true;
    error = '';

    try {
      const wallpaper =
        await downloadWallpaper(
          value
        );

      wallpapers = [
        wallpaper,
        ...wallpapers.filter(
          (candidate) =>
            candidate.filename !==
            wallpaper.filename
        )
      ];

      selectWallpaper(
        wallpaper
      );

      remoteURL = '';
    } catch (reason) {
      error =
        reason instanceof Error
          ? reason.message
          : String(reason);
    } finally {
      busy = false;
    }
  }

  async function removeWallpaper(
    wallpaper: WallpaperFile
  ) {
    const confirmed =
      window.confirm(
        `Delete wallpaper "${wallpaper.filename}"?`
      );

    if (!confirmed) {
      return;
    }

    busy = true;
    error = '';

    try {
      await deleteWallpaper(
        wallpaper.filename
      );

      wallpapers =
        wallpapers.filter(
          (candidate) =>
            candidate.filename !==
            wallpaper.filename
        );

      if (
        config.appearance.background.wallpaper ===
        wallpaper.filename
      ) {
        config.appearance.background.wallpaper =
          null;

        config.appearance.background.type =
          'solid';
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

    if (
      bytes <
      1024 * 1024
    ) {
      return `${
        (
          bytes / 1024
        ).toFixed(1)
      } KB`;
    }

    return `${
      (
        bytes /
        1024 /
        1024
      ).toFixed(1)
    } MB`;
  }
</script>

<section class="wallpaper-manager">
  <header class="wallpaper-manager-heading">
    <div>
      <strong>Wallpaper library</strong>

      <small>
        Upload an image from this computer or
        paste a direct image link.
      </small>
    </div>

    <button
      type="button"
      disabled={busy}
      onclick={() =>
        fileInput.click()}
    >
      Upload
    </button>
  </header>

  <input
    class="wallpaper-file-input"
    bind:this={fileInput}
    type="file"
    accept="image/jpeg,image/png,image/webp,image/gif"
    onchange={handleFileSelection}
  />

  <div class="wallpaper-url-control">
    <input
      bind:value={remoteURL}
      type="url"
      placeholder="https://example.com/wallpaper.webp"
      disabled={busy}
      onkeydown={(event) => {
        if (event.key === 'Enter') {
          event.preventDefault();
          handleURLDownload();
        }
      }}
    />

    <button
      type="button"
      disabled={
        busy ||
        !remoteURL.trim()
      }
      onclick={handleURLDownload}
    >
      {busy
        ? 'Working…'
        : 'Download'}
    </button>
  </div>

  <div class="wallpaper-manager-note">
    Use the direct image address obtained through
    “Copy image link”. Maximum file size is 20 MB.
  </div>

  {#if error}
    <div class="wallpaper-manager-error">
      {error}
    </div>
  {/if}

  {#if loading}
    <div class="wallpaper-manager-empty">
      Loading wallpapers…
    </div>
  {:else if wallpapers.length === 0}
    <button
      class="wallpaper-manager-empty actionable"
      type="button"
      onclick={() =>
        fileInput.click()}
    >
      <span>+</span>

      <strong>
        Add the first wallpaper
      </strong>
    </button>
  {:else}
    <div class="wallpaper-gallery">
      {#each wallpapers as wallpaper (wallpaper.filename)}
        <article
          class:selected={
            config.appearance.background.wallpaper ===
            wallpaper.filename
          }
          class="wallpaper-gallery-item"
        >
          <button
            class="wallpaper-gallery-select"
            type="button"
            title={wallpaper.filename}
            onclick={() =>
              selectWallpaper(
                wallpaper
              )}
          >
            <img
              src={wallpaper.url}
              alt=""
              loading="lazy"
            />

            <span
              class="wallpaper-gallery-selected"
            >
              ✓
            </span>
          </button>

          <footer>
            <div>
              <strong>
                {wallpaper.filename}
              </strong>

              <small>
                {readableSize(
                  wallpaper.size
                )}
              </small>
            </div>

            <button
              type="button"
              disabled={busy}
              title="Delete wallpaper"
              aria-label={`Delete ${wallpaper.filename}`}
              onclick={() =>
                removeWallpaper(
                  wallpaper
                )}
            >
              ×
            </button>
          </footer>
        </article>
      {/each}
    </div>
  {/if}
</section>
