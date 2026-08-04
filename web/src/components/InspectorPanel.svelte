<script lang="ts">
  import type {
    HubletConfig,
    Item,
    Section,
    SectionSurface,
    SectionWidth
  } from '../lib/types';

  import type {
    EditorSelection
  } from '../lib/editor';

  let {
    config,
    selection,
    onDeleteSection,
    onDeleteItem
  }: {
    config: HubletConfig;
    selection: EditorSelection;
    onDeleteSection: (sectionId: string) => void;
    onDeleteItem: (
      sectionId: string,
      itemId: string
    ) => void;
  } = $props();

  const accentSwatches = [
    '#4f8cff',
    '#8b5cf6',
    '#d946ef',
    '#ec4899',
    '#ef4444',
    '#f97316',
    '#eab308',
    '#22c55e',
    '#14b8a6',
    '#06b6d4',
    '#64748b',
    '#a3a3a3'
  ];

  const sectionWidths: Array<{
    value: number;
    legacyWidth: SectionWidth;
    label: string;
    description: string;
  }> = [
    {
      value: 6,
      legacyWidth: 'narrow',
      label: 'Narrow',
      description: '6 / 24'
    },
    {
      value: 8,
      legacyWidth: 'medium',
      label: 'Medium',
      description: '8 / 24'
    },
    {
      value: 12,
      legacyWidth: 'wide',
      label: 'Wide',
      description: '12 / 24'
    },
    {
      value: 16,
      legacyWidth: 'extra-wide',
      label: 'Extra wide',
      description: '16 / 24'
    },
    {
      value: 24,
      legacyWidth: 'full',
      label: 'Full',
      description: '24 / 24'
    }
  ];

  const sectionSurfaces: Array<{
    value: SectionSurface;
    label: string;
    description: string;
  }> = [
    {
      value: 'solid',
      label: 'Solid',
      description: 'Strong dark panel'
    },
    {
      value: 'transparent',
      label: 'Transparent',
      description: 'Subtle tinted surface'
    },
    {
      value: 'glass',
      label: 'Glass',
      description: 'Transparent with blur'
    },
    {
      value: 'none',
      label: 'None',
      description: 'No panel background'
    }
  ];

  const selectedSection = $derived.by(
    (): Section | null => {
      if (selection.type === 'dashboard') {
        return null;
      }

      return (
        config.sections.find(
          (section) =>
            section.id === selection.sectionId
        ) ?? null
      );
    }
  );

  const selectedItem = $derived.by(
    (): Item | null => {
      if (
        selection.type !== 'item' ||
        !selectedSection
      ) {
        return null;
      }

      return (
        selectedSection.items.find(
          (item) =>
            item.id === selection.itemId
        ) ?? null
      );
    }
  );

  let sectionDeleteArmed = $state(false);
  let itemDeleteArmed = $state(false);

  function setSectionColumnSpan(
    span: number
  ) {
    if (!selectedSection) {
      return;
    }

    const normalized = Math.max(
      4,
      Math.min(24, Math.round(span))
    );

    selectedSection.gridColumnSpan = normalized;

    if (normalized <= 6) {
      selectedSection.width = 'narrow';
    } else if (normalized <= 9) {
      selectedSection.width = 'medium';
    } else if (normalized <= 14) {
      selectedSection.width = 'wide';
    } else if (normalized <= 20) {
      selectedSection.width = 'extra-wide';
    } else {
      selectedSection.width = 'full';
    }
  }

  function requestSectionDelete() {
    if (
      selection.type !== 'section' &&
      selection.type !== 'item'
    ) {
      return;
    }

    if (!sectionDeleteArmed) {
      sectionDeleteArmed = true;

      window.setTimeout(() => {
        sectionDeleteArmed = false;
      }, 3000);

      return;
    }

    onDeleteSection(selection.sectionId);
  }

  function requestItemDelete() {
    if (selection.type !== 'item') {
      return;
    }

    if (!itemDeleteArmed) {
      itemDeleteArmed = true;

      window.setTimeout(() => {
        itemDeleteArmed = false;
      }, 3000);

      return;
    }

    onDeleteItem(
      selection.sectionId,
      selection.itemId
    );
  }

  import ServiceIcon from './ServiceIcon.svelte';
  import LocalIconManager from './LocalIconManager.svelte';

</script>

<aside class="inspector-panel">
  {#if selection.type === 'dashboard'}
    <header class="inspector-header">
      <div class="inspector-icon">
        ◫
      </div>

      <div>
        <p>DASHBOARD</p>
        <h2>Dashboard settings</h2>
      </div>
    </header>

    <div class="inspector-content">
      <label class="inspector-field">
        <span>Dashboard title</span>

        <input
          bind:value={config.dashboard.title}
          type="text"
          maxlength="80"
          placeholder="Hublet v2"
        />
      </label>

      <div class="inspector-group">
        <div class="inspector-group-heading">
          <div>
            <strong>Background</strong>

            <small>
              Choose the dashboard base color.
            </small>
          </div>
        </div>

        <label class="color-value-control">
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
          />
        </label>
      </div>

    </div>
  {:else if selectedSection && selection.type === 'section'}
    <header class="inspector-header">
      <div
        class="inspector-icon accent"
        style={`--inspector-accent:${selectedSection.accent}`}
      >
        ◧
      </div>

      <div>
        <p>SECTION</p>
        <h2>{selectedSection.title}</h2>
      </div>
    </header>

    <div class="inspector-content section-inspector-content">
      <label class="inspector-field">
        <span>Title</span>

        <input
          bind:value={selectedSection.title}
          type="text"
          maxlength="80"
          placeholder="Media"
        />
      </label>

      <div class="inspector-group accent-settings">
        <div class="inspector-group-heading">
          <div>
            <strong>Accent</strong>

            <small>
              Used for the section border and indicator.
            </small>
          </div>
        </div>

        <div class="accent-swatches">
          {#each accentSwatches as swatch}
            <button
              class:active={
                selectedSection.accent === swatch
              }
              type="button"
              style={`--swatch:${swatch}`}
              aria-label={`Use color ${swatch}`}
              onclick={() => {
                selectedSection.accent = swatch;
              }}
            >
              <span></span>
            </button>
          {/each}
        </div>

        <label class="accent-value-control">
          <span class="accent-color-input">
            <input
              bind:value={selectedSection.accent}
              type="color"
              aria-label="Choose accent color"
            />
          </span>

          <span class="accent-hex-input">
            <small>Hex color</small>

            <input
              bind:value={selectedSection.accent}
              type="text"
              maxlength="7"
              aria-label="Accent hex color"
              spellcheck="false"
            />
          </span>
        </label>
      </div>


      <div class="inspector-group section-surface-settings">
        <div class="inspector-group-heading">
          <div>
            <strong>Section surface</strong>

            <small>
              Controls the panel behind this section.
            </small>
          </div>
        </div>

        <div class="section-surface-options">
          {#each sectionSurfaces as option}
            <button
              class:active={
                selectedSection.surface === option.value
              }
              type="button"
              onclick={() => {
                selectedSection.surface = option.value;
              }}
            >
              <span
                class={`section-surface-preview ${option.value}`}
              >
                <i></i>
                <i></i>
              </span>

              <span>
                <strong>{option.label}</strong>
                <small>{option.description}</small>
              </span>
            </button>
          {/each}
        </div>

        {#if selectedSection.surface !== 'none'}
          <label class="section-surface-range">
            <span>
              <strong>Opacity</strong>

              <output>
                {selectedSection.surfaceOpacity}%
              </output>
            </span>

            <input
              bind:value={selectedSection.surfaceOpacity}
              type="range"
              min="0"
              max="100"
              step="1"
            />
          </label>
        {/if}

        {#if selectedSection.surface === 'glass'}
          <label class="section-surface-range">
            <span>
              <strong>Backdrop blur</strong>

              <output>
                {selectedSection.surfaceBlur}px
              </output>
            </span>

            <input
              bind:value={selectedSection.surfaceBlur}
              type="range"
              min="0"
              max="40"
              step="1"
            />
          </label>
        {/if}

        <label class="toggle-control section-toggle section-border-toggle">
          <span>
            <strong>Section border</strong>

            <small>
              Draw the accent outline around this panel.
            </small>
          </span>

          <input
            bind:checked={selectedSection.showBorder}
            type="checkbox"
          />
        </label>
      </div>

      <div class="inspector-group">
        <div class="inspector-group-heading">
          <div>
            <strong>Cards per row</strong>

            <small>
              Number of equal-size cards displayed on desktop.
            </small>
          </div>

          <span class="setting-value">
            {selectedSection.gridColumns}
          </span>
        </div>

        <div class="column-options">
          {#each [1, 2, 3, 4, 5, 6] as columns}
            <button
              class:active={
                selectedSection.gridColumns ===
                columns
              }
              type="button"
              aria-label={`${columns} cards per row`}
              onclick={() => {
                selectedSection.gridColumns =
                  columns;
              }}
            >
              {columns}
            </button>
          {/each}
        </div>
      </div>

      <div class="inspector-group">
        <div class="inspector-group-heading">
          <div>
            <strong>Section width</strong>

            <small>
              {selectedSection.gridColumnSpan} of 24 desktop columns.
            </small>
          </div>
        </div>

        <div class="section-width-options">
          {#each sectionWidths as option}
            <button
              class:active={
                selectedSection.gridColumnSpan ===
                option.value
              }
              type="button"
              onclick={() => {
                setSectionColumnSpan(option.value);
              }}
            >
              <span
                class={`section-width-preview ${option.legacyWidth}`}
              ></span>

              <span>
                <strong>{option.label}</strong>

                <small>
                  {option.description}
                </small>
              </span>
            </button>
          {/each}
        </div>

        <label class="section-width-range">
          <span>
            <strong>Custom width</strong>
            <small>{selectedSection.gridColumnSpan} / 24</small>
          </span>

          <input
            type="range"
            min="4"
            max="24"
            step="1"
            value={selectedSection.gridColumnSpan}
            oninput={(event) => {
              setSectionColumnSpan(
                Number(event.currentTarget.value)
              );
            }}
          />
        </label>
      </div>

      <div class="inspector-danger-zone">
        <button
          class:armed={sectionDeleteArmed}
          type="button"
          onclick={requestSectionDelete}
        >
          {sectionDeleteArmed
            ? 'Confirm section deletion'
            : 'Delete section'}
        </button>
      </div>
    </div>
  {:else if selectedSection && selectedItem && selection.type === 'item'}
    <header class="inspector-header">
      <div class="inspector-icon item">
        {selectedItem.name
          .trim()
          .slice(0, 1)
          .toUpperCase() || '?'}
      </div>

      <div>
        <p>SERVICE</p>

        <h2>
          {selectedItem.name || 'Untitled service'}
        </h2>
      </div>
    </header>

    <div class="inspector-content">
      <div class="service-type-badge">
        <span>Service card</span>

        <small>
          Opens a self-hosted service or website.
        </small>
      </div>

      <label class="inspector-field primary">
        <span>URL</span>

        <input
          bind:value={selectedItem.url}
          type="text"
          placeholder="proxmox.lan:8006"
        />

        <small>
          Protocol will be added automatically when saving.
        </small>
      </label>

      <label class="inspector-field">
        <span>Name</span>

        <input
          bind:value={selectedItem.name}
          type="text"
          maxlength="100"
          placeholder="Proxmox"
        />
      </label>

      <label class="inspector-field">
        <span>Description</span>

        <input
          bind:value={selectedItem.description}
          type="text"
          maxlength="240"
          placeholder="Virtualization server"
        />
      </label>

      <div class="inspector-group">
        <div class="inspector-group-heading">
          <div>
            <strong>Icon</strong>

            <small>
              Automatic icon discovery comes next.
            </small>
          </div>
        </div>

        <div class="segmented-control">
          <button
            class:active={
              selectedItem.icon.type === 'auto'
            }
            type="button"
            onclick={() => {
              selectedItem.icon.type = 'auto';
              selectedItem.icon.value = '';
            }}
          >
            Auto
          </button>

          <button
            class:active={
              selectedItem.icon.type === 'local'
            }
            type="button"
            onclick={() => {
              selectedItem.icon.type = 'local';
            }}
          >
            Local
          </button>

          <button
            class:active={
              selectedItem.icon.type === 'none'
            }
            type="button"
            onclick={() => {
              selectedItem.icon.type = 'none';
              selectedItem.icon.value = '';
            }}
          >
            None
          </button>
        </div>

        {#if selectedItem.icon.type === 'local'}
          <LocalIconManager item={selectedItem} />
        {/if}
      </div>

      <label class="toggle-control">
        <span>
          <strong>Open in new tab</strong>

          <small>
            Keep Hublet v2 open in the current tab.
          </small>
        </span>

        <input
          bind:checked={
            selectedItem.openInNewTab
          }
          type="checkbox"
        />
      </label>

      <div class="inspector-danger-zone">
        <button
          class:armed={itemDeleteArmed}
          type="button"
          onclick={requestItemDelete}
        >
          {itemDeleteArmed
            ? 'Confirm service deletion'
            : 'Delete service'}
        </button>
      </div>
    </div>
  {/if}
</aside>
