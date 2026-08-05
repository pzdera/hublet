<script lang="ts">
  import {
    untrack
  } from 'svelte';

  import type {
    HubletConfig,
    Item,
    Section
  } from '../lib/types';

  import type {
    EditorSelection
  } from '../lib/editor';

  import { createID } from '../lib/id';

  import {
    firstAvailablePlacement
  } from '../lib/section-placement';

  import EditorCanvas from './EditorCanvas.svelte';
  import InspectorPanel from './InspectorPanel.svelte';
  import DashboardSettingsPanel from './DashboardSettingsPanel.svelte';

  let {
    config,
    saving,
    error,
    onSave,
    onCancel
  }: {
    config: HubletConfig;
    saving: boolean;
    error: string;
    onSave: () => void;
    onCancel: () => void;
  } = $props();

  let selection = $state<EditorSelection>({
    type: 'dashboard'
  });

  let layoutError = $state('');

  const originalSnapshot = JSON.stringify(
    $state.snapshot(
      untrack(() => config)
    )
  );

  const currentSnapshot = $derived(
    JSON.stringify(
      $state.snapshot(config)
    )
  );

  const dirty = $derived(
    currentSnapshot !== originalSnapshot
  );

  function addSection() {
    layoutError = '';

    const section: Section = {
      id: createID(
        'section',
        'new-section'
      ),
      title: 'New section',
      accent: '#4f8cff',
      surface: 'solid',
      surfaceOpacity: 82,
      surfaceBlur: 16,
      showBorder: true,
      width: 'narrow',
      gridRow: 0,
      gridColumn: 0,
      gridRowSpan: 10,
      gridColumnSpan: 6,
      gridColumns: 1,
      items: []
    };

    const candidateSections = [
      ...config.sections,
      section
    ];

    const placement = firstAvailablePlacement(
      candidateSections,
      section.id
    );

    if (!placement) {
      layoutError =
        'Desktop grid is full. Move or delete a section before adding another one.';
      return;
    }

    section.gridRow = placement.row;
    section.gridColumn = placement.column;

    config.sections.push(section);

    selection = {
      type: 'section',
      sectionId: section.id
    };
  }

  function addItem(
    sectionId: string
  ) {
    const section =
      config.sections.find(
        (candidate) =>
          candidate.id === sectionId
      );

    if (!section) {
      return;
    }

    const item: Item = {
      id: createID(
        'service',
        'new-service'
      ),
      type: 'service',
      name: 'New service',
      url: '',
      description: '',
      icon: {
        type: 'auto',
        value: ''
      }
    };

    section.items.push(item);

    selection = {
      type: 'item',
      sectionId,
      itemId: item.id
    };
  }

  function deleteSection(
    sectionId: string
  ) {
    layoutError = '';

    const sectionIndex =
      config.sections.findIndex(
        (section) =>
          section.id === sectionId
      );

    if (sectionIndex === -1) {
      return;
    }

    config.sections.splice(
      sectionIndex,
      1
    );

    const replacement =
      config.sections[sectionIndex] ??
      config.sections[sectionIndex - 1];

    selection = replacement
      ? {
          type: 'section',
          sectionId: replacement.id
        }
      : {
          type: 'dashboard'
        };
  }

  function deleteItem(
    sectionId: string,
    itemId: string
  ) {
    const section =
      config.sections.find(
        (candidate) =>
          candidate.id === sectionId
      );

    if (!section) {
      return;
    }

    const itemIndex =
      section.items.findIndex(
        (item) =>
          item.id === itemId
      );

    if (itemIndex === -1) {
      return;
    }

    section.items.splice(
      itemIndex,
      1
    );

    const replacement =
      section.items[itemIndex] ??
      section.items[itemIndex - 1];

    selection = replacement
      ? {
          type: 'item',
          sectionId,
          itemId: replacement.id
        }
      : {
          type: 'section',
          sectionId
        };
  }
</script>

<div class="editor-workspace">
  <header class="workspace-toolbar">
    <div class="workspace-brand">
      <span class="workspace-logo">
        H
      </span>

      <div>
        <p>HUBLET EDITOR</p>

        <h2>
          {config.dashboard.title}
        </h2>
      </div>
    </div>

    <div class="workspace-center-status">
      {#if dirty}
        <span class="unsaved-dot"></span>
        <span>Unsaved changes</span>
      {:else}
        <span class="saved-dot"></span>
        <span>All changes saved</span>
      {/if}
    </div>

    <div class="workspace-toolbar-actions">
      <button
        class="workspace-cancel"
        type="button"
        disabled={saving}
        onclick={onCancel}
      >
        Cancel
      </button>

      <button
        class="workspace-save"
        type="button"
        disabled={saving || !dirty}
        onclick={onSave}
      >
        {saving ? 'Saving…' : 'Save'}
      </button>
    </div>
  </header>

  {#if error || layoutError}
    <div class="workspace-error">
      {error || layoutError}
    </div>
  {/if}

  <div class="workspace-body">
    <EditorCanvas
      {config}
      {selection}
      onSelect={(value) => {
        selection = value;
      }}
      onAddSection={addSection}
      onAddItem={addItem}
    />

    {#if selection.type === 'dashboard'}
      <DashboardSettingsPanel {config} />
    {:else}
      <InspectorPanel
        {config}
        {selection}
        onDeleteSection={deleteSection}
        onDeleteItem={deleteItem}
      />
    {/if}
  </div>

  {#if dirty}
    <footer class="unsaved-bar">
      <div>
        <span class="unsaved-dot"></span>

        <div>
          <strong>
            Unsaved changes
          </strong>

          <small>
            Review the live preview before saving.
          </small>
        </div>
      </div>

      <div>
        <button
          class="workspace-cancel"
          type="button"
          disabled={saving}
          onclick={onCancel}
        >
          Discard
        </button>

        <button
          class="workspace-save"
          type="button"
          disabled={saving}
          onclick={onSave}
        >
          {saving
            ? 'Saving…'
            : 'Save changes'}
        </button>
      </div>
    </footer>
  {/if}
</div>
