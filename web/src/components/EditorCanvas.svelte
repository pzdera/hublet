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
    dashboardGridColumns,
    dashboardGridRowStep,
    maximumGridBottom,
    normalizeGridLayout,
    placementCollisions,
    placementFits,
    sectionColumnSpan,
    sectionPlacementStyle,
    sectionRowSpan,
    sectionSwapPlan
  } from '../lib/section-placement';

  import {
    measureSection
  } from '../lib/section-measure';

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

  type SectionDrag = {
    sectionId: string;
    pointerId: number;
    x: number;
    y: number;
    columnOffset: number;
    rowOffset: number;
    columnSpan: number;
    rowSpan: number;
    targetRow: number;
    targetColumn: number;
    valid: boolean;
    swapTargetId: string;
    swapTargetName: string;
    swapTargetRow: number;
    swapTargetColumn: number;
    swapTargetColumnSpan: number;
    swapTargetRowSpan: number;
  };

  let sectionDrag = $state<SectionDrag | null>(null);

  const maximumGridRow = $derived(
    maximumGridBottom(config.sections)
  );

  const layoutSignature = $derived(
    config.sections
      .map((section) => [
        section.id,
        section.width,
        section.gridRow,
        section.gridColumn,
        section.gridRowSpan,
        section.gridColumnSpan
      ].join(':'))
      .join('|')
  );

  $effect(() => {
    layoutSignature;
    normalizeGridLayout(config.sections);
  });

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

  function startSectionDrag(
    event: PointerEvent,
    section: Section
  ) {
    if (event.button !== 0) {
      return;
    }

    const article = (
      event.currentTarget as HTMLElement
    ).closest<HTMLElement>(
      '.editor-canvas-section'
    );

    if (!article) {
      return;
    }

    event.preventDefault();
    event.stopPropagation();

    const rect = article.getBoundingClientRect();
    const columnSpan = sectionColumnSpan(section);
    const rowSpan = sectionRowSpan(section);
    const relativeX = Math.max(
      0,
      Math.min(rect.width - 1, event.clientX - rect.left)
    );

    sectionDrag = {
      sectionId: section.id,
      pointerId: event.pointerId,
      x: event.clientX,
      y: event.clientY,
      columnOffset: Math.min(
        columnSpan - 1,
        Math.floor(
          relativeX / (rect.width / columnSpan)
        )
      ),
      rowOffset: Math.min(
        rowSpan - 1,
        Math.max(
          0,
          Math.floor(
            (event.clientY - rect.top) /
            dashboardGridRowStep
          )
        )
      ),
      columnSpan,
      rowSpan,
      targetRow: section.gridRow,
      targetColumn: section.gridColumn,
      valid: true,
      swapTargetId: '',
      swapTargetName: '',
      swapTargetRow: 0,
      swapTargetColumn: 0,
      swapTargetColumnSpan: 0,
      swapTargetRowSpan: 0
    };

    document.body.classList.add(
      'section-grid-dragging'
    );
  }

  function updateSectionDrag(
    event: PointerEvent
  ) {
    if (
      !sectionDrag ||
      event.pointerId !== sectionDrag.pointerId
    ) {
      return;
    }

    const grid = document.querySelector<HTMLElement>(
      '.editor-canvas-grid'
    );

    const section = config.sections.find(
      (candidate) => candidate.id === sectionDrag?.sectionId
    );

    if (!grid || !section) {
      return;
    }

    event.preventDefault();

    const rect = grid.getBoundingClientRect();
    const gap = 16;
    const columnWidth = (
      rect.width - gap * (dashboardGridColumns - 1)
    ) / dashboardGridColumns;
    const step = columnWidth + gap;
    const hoveredColumn = Math.max(
      1,
      Math.min(
        dashboardGridColumns,
        Math.floor((event.clientX - rect.left) / step) + 1
      )
    );
    const columnSpan = sectionColumnSpan(section);
    const maximumColumn =
      dashboardGridColumns - columnSpan + 1;
    const targetColumn = Math.max(
      1,
      Math.min(
        maximumColumn,
        hoveredColumn - sectionDrag.columnOffset
      )
    );
    const targetRow = Math.max(
      1,
      Math.floor(
        (event.clientY - rect.top) /
        dashboardGridRowStep
      ) + 1 - sectionDrag.rowOffset
    );

    const collisions = placementCollisions(
      config.sections,
      section.id,
      targetRow,
      targetColumn
    );
    let resolvedRow = targetRow;
    let resolvedColumn = targetColumn;
    let valid = placementFits(
      config.sections,
      section.id,
      targetRow,
      targetColumn
    );
    let swapTargetId = '';
    let swapTargetName = '';
    let swapTargetRow = 0;
    let swapTargetColumn = 0;
    let swapTargetColumnSpan = 0;
    let swapTargetRowSpan = 0;

    if (collisions.length === 1) {
      const target = collisions[0];
      const swapPlan = sectionSwapPlan(
        config.sections,
        section.id,
        target.id,
        targetRow,
        targetColumn
      );

      swapTargetName = target.title;

      if (swapPlan) {
        valid = true;
        resolvedRow = swapPlan.movingRow;
        resolvedColumn = swapPlan.movingColumn;
        swapTargetId = target.id;
        swapTargetRow = swapPlan.targetRow;
        swapTargetColumn = swapPlan.targetColumn;
        swapTargetColumnSpan = sectionColumnSpan(target);
        swapTargetRowSpan = sectionRowSpan(target);
      }
    } else if (collisions.length > 1) {
      valid = false;
    }

    sectionDrag.x = event.clientX;
    sectionDrag.y = event.clientY;
    sectionDrag.targetRow = resolvedRow;
    sectionDrag.targetColumn = resolvedColumn;
    sectionDrag.valid = valid;
    sectionDrag.swapTargetId = swapTargetId;
    sectionDrag.swapTargetName = swapTargetName;
    sectionDrag.swapTargetRow = swapTargetRow;
    sectionDrag.swapTargetColumn = swapTargetColumn;
    sectionDrag.swapTargetColumnSpan = swapTargetColumnSpan;
    sectionDrag.swapTargetRowSpan = swapTargetRowSpan;
  }

  function finishSectionDrag(
    event: PointerEvent,
    cancelled = false
  ) {
    if (
      !sectionDrag ||
      event.pointerId !== sectionDrag.pointerId
    ) {
      return;
    }

    const completedDrag = sectionDrag;
    sectionDrag = null;

    document.body.classList.remove(
      'section-grid-dragging'
    );

    if (cancelled || !completedDrag.valid) {
      return;
    }

    const section = config.sections.find(
      (candidate) => candidate.id === completedDrag.sectionId
    );

    if (!section) {
      return;
    }

    if (completedDrag.swapTargetId) {
      const target = config.sections.find(
        (candidate) =>
          candidate.id === completedDrag.swapTargetId
      );
      const swapPlan = sectionSwapPlan(
        config.sections,
        section.id,
        completedDrag.swapTargetId,
        completedDrag.targetRow,
        completedDrag.targetColumn
      );

      if (!target || !swapPlan) {
        return;
      }

      section.gridRow = swapPlan.movingRow;
      section.gridColumn = swapPlan.movingColumn;
      target.gridRow = swapPlan.targetRow;
      target.gridColumn = swapPlan.targetColumn;
    } else {
      section.gridRow = completedDrag.targetRow;
      section.gridColumn = completedDrag.targetColumn;
    }

    config.sections.sort((left, right) => {
      return (
        left.gridRow - right.gridRow ||
        left.gridColumn - right.gridColumn
      );
    });
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
    section: Section
  ): string {
    return `${sectionColumnSpan(section)}/24`;
  }
</script>

<svelte:window
  onpointermove={updateSectionDrag}
  onpointerup={(event) => finishSectionDrag(event)}
  onpointercancel={(event) => finishSectionDrag(event, true)}
/>

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
      class:drag-active={sectionDrag !== null}
      class="editor-canvas-grid"
      aria-label="Dashboard preview"
    >
      {#each config.sections as section (section.id)}
        <article
          class={sectionClass(section)}
          data-section-id={section.id}
          data-grid-row={section.gridRow}
          use:measureSection={{
            section,
            sections: config.sections
          }}
          style={[
            sectionSurfaceStyle(section),
            sectionPlacementStyle(section)
          ].filter(Boolean).join(';')}
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
              onpointerdown={(event) =>
                startSectionDrag(event, section)}
              onclick={(event) => event.stopPropagation()}
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
                section
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

      {#if sectionDrag}
        <div
          class="section-grid-new-row"
          style={`grid-row:${maximumGridRow + 1} / span 10`}
          aria-hidden="true"
        ></div>

        <div
          class:invalid={!sectionDrag.valid}
          class:swap={Boolean(sectionDrag.swapTargetId)}
          class="section-grid-drop-preview"
          style={[
            `grid-row:${sectionDrag.targetRow}`,
            `grid-column:${sectionDrag.targetColumn} / span ${sectionDrag.columnSpan}`,
            `grid-row-end:span ${sectionDrag.rowSpan}`
          ].join(';')}
          aria-hidden="true"
        >
          <span>
            {sectionDrag.swapTargetId
              ? `Swap with ${sectionDrag.swapTargetName}`
              : sectionDrag.valid
                ? 'Drop section here'
                : sectionDrag.swapTargetName
                  ? 'Swap unavailable'
                  : 'Position occupied'}
          </span>
        </div>

        {#if sectionDrag.swapTargetId}
          <div
            class="section-grid-drop-preview swap-origin"
            style={[
              `grid-row:${sectionDrag.swapTargetRow} / span ${sectionDrag.swapTargetRowSpan}`,
              `grid-column:${sectionDrag.swapTargetColumn} / span ${sectionDrag.swapTargetColumnSpan}`
            ].join(';')}
            aria-hidden="true"
          >
            <span>
              Move {sectionDrag.swapTargetName} here
            </span>
          </div>
        {/if}
      {/if}
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

{#if sectionDrag}
  <div
    class="section-grid-drag-float"
    style={`transform:translate3d(${sectionDrag.x + 14}px, ${sectionDrag.y + 14}px, 0)`}
    aria-hidden="true"
  >
    <span></span>

    {config.sections.find(
      (section) => section.id === sectionDrag?.sectionId
    )?.title ?? 'Section'}
  </div>
{/if}
