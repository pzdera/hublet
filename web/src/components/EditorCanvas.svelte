<script lang="ts">
  import { flip } from 'svelte/animate';

  import {
    dragHandle,
    dragHandleZone,
    SHADOW_ITEM_MARKER_PROPERTY_NAME
  } from 'svelte-dnd-action';

  import type {
    DndEvent
  } from 'svelte-dnd-action';

  import type {
    HubletConfig,
    Item,
    Section
  } from '../lib/types';

  import type {
    EditorSelection
  } from '../lib/editor';

  import {
    sectionSurfaceStyle
  } from '../lib/section-surface';

  import {
    sectionSortable
  } from '../lib/section-sortable';

  import ServiceIcon from './ServiceIcon.svelte';

  let {
    config,
    selection,
    onSelect,
    onAddSection,
    onAddItem
  }: {
    config: HubletConfig;
    selection: EditorSelection;
    onSelect: (selection: EditorSelection) => void;
    onAddSection: () => void;
    onAddItem: (sectionId: string) => void;
  } = $props();

  const flipDurationMs = 180;

  function isShadowItem(
    item: Item
  ): boolean {
    const record =
      item as unknown as Record<
        string,
        unknown
      >;

    return Boolean(
      record[
        SHADOW_ITEM_MARKER_PROPERTY_NAME
      ]
    );
  }

  function itemDndKey(
    item: Item
  ): string {
    return isShadowItem(item)
      ? `${item.id}-shadow`
      : item.id;
  }

  function handleItemsDnd(
    sectionId: string,
    event: CustomEvent<DndEvent<Item>>
  ) {
    const section = config.sections.find(
      (candidate) =>
        candidate.id === sectionId
    );

    if (!section) {
      return;
    }

    section.items =
      event.detail.items as Item[];
  }

  function reorderSections(
    orderedSectionIDs: string[]
  ) {
    if (
      orderedSectionIDs.length !==
      config.sections.length
    ) {
      return;
    }

    const sectionsByID = new Map(
      config.sections.map(
        (section) => [
          section.id,
          section
        ]
      )
    );

    const reordered: Section[] = [];

    for (
      const sectionID
      of orderedSectionIDs
    ) {
      const section =
        sectionsByID.get(sectionID);

      if (!section) {
        return;
      }

      reordered.push(section);
    }

    config.sections = reordered;
  }

  function sectionSelected(
    sectionId: string
  ): boolean {
    return (
      selection.type === 'section' &&
      selection.sectionId === sectionId
    );
  }

  function itemSelected(
    sectionId: string,
    itemId: string
  ): boolean {
    return (
      selection.type === 'item' &&
      selection.sectionId === sectionId &&
      selection.itemId === itemId
    );
  }

  function selectDashboard() {
    onSelect({
      type: 'dashboard'
    });
  }

  function selectSection(
    sectionId: string
  ) {
    onSelect({
      type: 'section',
      sectionId
    });
  }

  function selectItem(
    sectionId: string,
    itemId: string
  ) {
    onSelect({
      type: 'item',
      sectionId,
      itemId
    });
  }

  function itemInitial(
    item: Item
  ): string {
    return (
      item.name
        .trim()
        .slice(0, 1)
        .toUpperCase() || '?'
    );
  }

  function iconURL(
    item: Item
  ): string {
    if (
      item.icon.type === 'local' &&
      item.icon.value
    ) {
      return `/icons/${item.icon.value}`;
    }

    return '';
  }

  function effectiveCardSize(
    section: Section
  ): string {
    if (
      section.cardSize === 'inherit'
    ) {
      return config.appearance.cards.size;
    }

    return section.cardSize;
  }

  function sectionClass(
    section: Section
  ): string {
    return [
      'editor-canvas-section',
      `canvas-width-${section.width}`,
      `surface-${section.surface}`,
      sectionSelected(section.id)
        ? 'selected'
        : ''
    ]
      .filter(Boolean)
      .join(' ');
  }

  function cardClass(
    section: Section,
    item: Item
  ): string {
    return [
      'canvas-card',
      `card-size-${effectiveCardSize(
        section
      )}`,
      itemSelected(
        section.id,
        item.id
      )
        ? 'selected'
        : ''
    ]
      .filter(Boolean)
      .join(' ');
  }

  function sectionWidthLabel(
    width: Section['width']
  ): string {
    switch (width) {
      case 'narrow':
        return '3/12';

      case 'medium':
        return '4/12';

      case 'wide':
        return '6/12';

      case 'extra-wide':
        return '8/12';

      case 'full':
        return '12/12';
    }
  }
</script>

<div class="editor-canvas-shell">
  <div class="canvas-toolbar">
    <div>
      <span
        class="canvas-status-dot"
      ></span>

      <span>Live preview</span>
    </div>

    <span>
      Drag sections by the handle in their header
      and cards by the handle inside each card
    </span>
  </div>

  <div
    class="editor-canvas-page"
    style={`--canvas-background:${config.appearance.background.color}`}
  >
    <header
      class:selected={
        selection.type === 'dashboard'
      }
      class="canvas-header"
      onclick={(event) => {
        event.stopPropagation();
        selectDashboard();
      }}
    >
      <div class="canvas-brand">
        <span class="canvas-brand-mark">
          H
        </span>

        <div>
          <small>
            SELF-HOSTED DASHBOARD
          </small>

          <h1>
            {config.dashboard.title}
          </h1>
        </div>
      </div>

      <div class="canvas-search">
        <span>⌕</span>

        <span>
          Search services or enter a shortcut…
        </span>
      </div>
    </header>

    <main
      class="editor-canvas-grid"
      aria-label="Dashboard preview"
      use:sectionSortable={{
        enabled: true,
        onReorder: reorderSections
      }}
    >
      {#each config.sections as section (section.id)}
        <article
          class={sectionClass(section)}
          data-section-id={section.id}
          style={sectionSurfaceStyle(section)}
          aria-label={section.title}
          onclick={(event) => {
            event.stopPropagation();

            selectSection(section.id);
          }}
        >
          <header
            class="canvas-section-header"
          >
            <button
              class="canvas-section-drag-handle"
              type="button"
              aria-label={`Move section ${section.title}`}
              title="Drag to move section"
              onclick={(event) => {
                event.stopPropagation();
              }}
            >
              <i></i>
              <i></i>
              <i></i>
              <i></i>
              <i></i>
              <i></i>
            </button>

            <span
              class="canvas-accent-dot"
            ></span>

            <h2>
              {section.title}
            </h2>

            <span
              class="canvas-section-width"
            >
              {sectionWidthLabel(
                section.width
              )}
            </span>

            <span
              class="canvas-item-count"
            >
              {section.items.length}
            </span>
          </header>

          {#if !section.collapsed}
            <div
              class={[
                'canvas-cards',
                `arrangement-${section.layout}`,
                `card-size-${effectiveCardSize(
                  section
                )}`
              ].join(' ')}
              style={`--grid-columns:${section.gridColumns}`}
              aria-label={`${section.title} cards`}
              use:dragHandleZone={{
                items: section.items,
                type: 'hublet-service-cards',
                flipDurationMs,
                delayTouchStart: true,
                useCursorForDetection: true,
                dropTargetClasses: [
                  'dnd-card-zone-active'
                ]
              }}
              onconsider={(event) =>
                handleItemsDnd(
                  section.id,
                  event
                )}
              onfinalize={(event) =>
                handleItemsDnd(
                  section.id,
                  event
                )}
            >
              {#each section.items as item (itemDndKey(item))}
                <div
                  class="canvas-card-wrapper"
                  aria-label={item.name}
                  animate:flip={{
                    duration:
                      flipDurationMs
                  }}
                >
                  {#if isShadowItem(item)}
                    <div
                      class="card-drop-placeholder"
                    >
                      <span></span>

                      <strong>
                        Drop card here
                      </strong>
                    </div>
                  {:else}
                    <button
                      class={cardClass(
                        section,
                        item
                      )}
                      type="button"
                      onclick={(event) => {
                        event.stopPropagation();

                        selectItem(
                          section.id,
                          item.id
                        );
                      }}
                    >
                      <span
                        class="canvas-card-handle"
                        use:dragHandle
                        aria-label={`Move ${item.name}`}
                        title="Move card"
                      >
                        <i></i>
                        <i></i>
                        <i></i>
                        <i></i>
                        <i></i>
                        <i></i>
                      </span>

                      <span
                        class="canvas-card-icon"
                      >
                        <ServiceIcon {item} />
                      </span>

                      <span
                        class="canvas-card-copy"
                      >
                        <strong>
                          {item.name ||
                            'Untitled service'}
                        </strong>

                        {#if (
                          section.layout !==
                          'compact'
                        )}
                          <small>
                            {item.description ||
                              item.url ||
                              'No URL'}
                          </small>
                        {/if}
                      </span>

                      <span
                        class="canvas-card-edit"
                      >
                        ✎
                      </span>
                    </button>
                  {/if}
                </div>
              {/each}
            </div>

            <button
              class="canvas-add-card"
              type="button"
              onclick={(event) => {
                event.stopPropagation();

                onAddItem(section.id);
              }}
            >
              <span>+</span>
              <strong>Add service</strong>
            </button>
          {/if}
        </article>
      {/each}
    </main>

    <button
      class="canvas-add-section"
      type="button"
      onclick={(event) => {
        event.stopPropagation();
        onAddSection();
      }}
    >
      <span>+</span>

      <div>
        <strong>
          Add section
        </strong>

        <small>
          Create another group of cards
        </small>
      </div>
    </button>
  </div>
</div>
