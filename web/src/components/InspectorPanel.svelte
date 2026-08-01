<script lang="ts">
  import type {
    HubletConfig,
    Item,
    Section,
    SectionCardSize,
    SectionLayout,
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
    value: SectionWidth;
    label: string;
    description: string;
  }> = [
    {
      value: 'narrow',
      label: 'Narrow',
      description: '3 / 12'
    },
    {
      value: 'medium',
      label: 'Medium',
      description: '4 / 12'
    },
    {
      value: 'wide',
      label: 'Wide',
      description: '6 / 12'
    },
    {
      value: 'extra-wide',
      label: 'Extra wide',
      description: '8 / 12'
    },
    {
      value: 'full',
      label: 'Full',
      description: '12 / 12'
    }
  ];

  const cardSizes: Array<{
    value: SectionCardSize;
    label: string;
  }> = [
    {
      value: 'inherit',
      label: 'Default'
    },
    {
      value: 'small',
      label: 'Small'
    },
    {
      value: 'medium',
      label: 'Medium'
    },
    {
      value: 'large',
      label: 'Large'
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

  function selectLayout(
    layout: SectionLayout
  ) {
    if (!selectedSection) {
      return;
    }

    selectedSection.layout = layout;

    if (
      layout === 'list' ||
      layout === 'featured'
    ) {
      selectedSection.gridColumns = 1;
    }

    if (
      layout === 'grid' &&
      selectedSection.gridColumns < 2
    ) {
      selectedSection.gridColumns = 2;
    }

    if (
      layout === 'compact' &&
      selectedSection.gridColumns < 2
    ) {
      selectedSection.gridColumns = 3;
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
          placeholder="Hublet"
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

      <label class="toggle-control">
        <span>
          <strong>Animations</strong>
          <small>
            Use subtle movement and transitions.
          </small>
        </span>

        <input
          bind:checked={
            config.appearance.animations
          }
          type="checkbox"
        />
      </label>
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

    <div class="inspector-content">
      <label class="inspector-field">
        <span>Title</span>

        <input
          bind:value={selectedSection.title}
          type="text"
          maxlength="80"
          placeholder="Media"
        />
      </label>

      <div class="inspector-group">
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

        <details class="advanced-color">
          <summary>Advanced color</summary>

          <label class="color-value-control">
            <input
              bind:value={selectedSection.accent}
              type="color"
            />

            <input
              bind:value={selectedSection.accent}
              type="text"
              maxlength="7"
            />
          </label>
        </details>
      </div>

      <div class="inspector-group">
        <div class="inspector-group-heading">
          <div>
            <strong>Card arrangement</strong>
            <small>
              Controls how cards flow inside this section.
            </small>
          </div>
        </div>

        <div class="arrangement-options">
          <button
            class:active={
              selectedSection.layout === 'list'
            }
            type="button"
            onclick={() =>
              selectLayout('list')}
          >
            <span class="arrangement-preview preview-list">
              <i></i>
              <i></i>
              <i></i>
            </span>

            <span>
              <strong>List</strong>
              <small>One card per row</small>
            </span>
          </button>

          <button
            class:active={
              selectedSection.layout === 'grid'
            }
            type="button"
            onclick={() =>
              selectLayout('grid')}
          >
            <span class="arrangement-preview preview-grid">
              <i></i>
              <i></i>
              <i></i>
              <i></i>
            </span>

            <span>
              <strong>Grid</strong>
              <small>Multiple columns</small>
            </span>
          </button>

          <button
            class:active={
              selectedSection.layout === 'compact'
            }
            type="button"
            onclick={() =>
              selectLayout('compact')}
          >
            <span class="arrangement-preview preview-compact">
              <i></i>
              <i></i>
              <i></i>
              <i></i>
              <i></i>
              <i></i>
            </span>

            <span>
              <strong>Compact</strong>
              <small>Icon and name only</small>
            </span>
          </button>

          <button
            class:active={
              selectedSection.layout === 'featured'
            }
            type="button"
            onclick={() =>
              selectLayout('featured')}
          >
            <span class="arrangement-preview preview-featured">
              <i></i>
              <i></i>
            </span>

            <span>
              <strong>Featured</strong>
              <small>Large prominent cards</small>
            </span>
          </button>
        </div>
      </div>

      {#if (
        selectedSection.layout === 'grid' ||
        selectedSection.layout === 'compact'
      )}
        <div class="inspector-group">
          <div class="inspector-group-heading">
            <div>
              <strong>Columns</strong>
              <small>
                Maximum cards displayed in each row.
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
      {/if}

      <div class="inspector-group">
        <div class="inspector-group-heading">
          <div>
            <strong>Section width</strong>
            <small>
              Space occupied by the complete section.
            </small>
          </div>
        </div>

        <div class="section-width-options">
          {#each sectionWidths as option}
            <button
              class:active={
                selectedSection.width ===
                option.value
              }
              type="button"
              onclick={() => {
                selectedSection.width =
                  option.value;
              }}
            >
              <span
                class={`section-width-preview ${option.value}`}
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
      </div>

      <div class="inspector-group">
        <div class="inspector-group-heading">
          <div>
            <strong>Card size</strong>
            <small>
              Overrides the dashboard default in this section.
            </small>
          </div>
        </div>

        <div class="card-size-options">
          {#each cardSizes as option}
            <button
              class:active={
                selectedSection.cardSize ===
                option.value
              }
              type="button"
              onclick={() => {
                selectedSection.cardSize =
                  option.value;
              }}
            >
              <span
                class={`card-size-preview ${option.value}`}
              ></span>

              <strong>
                {option.label}
              </strong>
            </button>
          {/each}
        </div>
      </div>

      <label class="toggle-control">
        <span>
          <strong>Collapsed by default</strong>
          <small>
            Show only the section header at startup.
          </small>
        </span>

        <input
          bind:checked={selectedSection.collapsed}
          type="checkbox"
        />
      </label>

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
          <label class="inspector-field">
            <span>Filename</span>

            <input
              bind:value={selectedItem.icon.value}
              type="text"
              placeholder="proxmox.png"
            />
          </label>
        {/if}
      </div>

      <label class="toggle-control">
        <span>
          <strong>Open in new tab</strong>
          <small>
            Keep Hublet open in the current tab.
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
