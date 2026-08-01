<script lang="ts">
  import type {
    Section
  } from '../lib/types';

  import {
    maximumStartColumn,
    normalizeStartColumn,
    sectionWidthSpan
  } from '../lib/section-placement';

  let {
    section
  }: {
    section: Section;
  } = $props();

  let dragging = $state(false);

  const span = $derived(
    sectionWidthSpan(section.width)
  );

  const maximumStart = $derived(
    maximumStartColumn(section.width)
  );

  const selectedStart = $derived(
    section.startColumn > 0
      ? section.startColumn
      : 0
  );

  $effect(() => {
    section.width;
    normalizeStartColumn(section);
  });

  function chooseStart(
    column: number
  ) {
    if (
      column < 1 ||
      column > maximumStart
    ) {
      return;
    }

    section.startColumn = column;
  }

  function beginsSelectedRange(
    column: number
  ): boolean {
    return selectedStart === column;
  }

  function insideSelectedRange(
    column: number
  ): boolean {
    if (selectedStart === 0) {
      return false;
    }

    return (
      column >= selectedStart &&
      column < selectedStart + span
    );
  }

  function unavailable(
    column: number
  ): boolean {
    return column > maximumStart;
  }

  function finishDragging() {
    dragging = false;
  }
</script>

<svelte:window
  onpointerup={finishDragging}
  onpointercancel={finishDragging}
/>

<div class="section-placement-control">
  <div class="section-placement-toolbar">
    <button
      class:active={selectedStart === 0}
      type="button"
      onclick={() => {
        section.startColumn = 0;
      }}
    >
      Auto
    </button>

    <span>
      {#if selectedStart === 0}
        Automatic placement
      {:else}
        Columns {selectedStart}–{selectedStart + span - 1}
      {/if}
    </span>
  </div>

  <div
    class="section-placement-grid"
    role="group"
    aria-label="Section start column"
    onpointerleave={() => {
      dragging = false;
    }}
  >
    {#each Array.from({ length: 12 }, (_, index) => index + 1) as column}
      <button
        class:range-start={beginsSelectedRange(column)}
        class:selected={insideSelectedRange(column)}
        class:unavailable={unavailable(column)}
        type="button"
        disabled={unavailable(column)}
        aria-label={`Start section at column ${column}`}
        title={
          unavailable(column)
            ? `Section width does not fit from column ${column}`
            : `Start at column ${column}`
        }
        onpointerdown={(event) => {
          if (unavailable(column)) {
            return;
          }

          event.preventDefault();
          dragging = true;
          chooseStart(column);
        }}
        onpointerenter={() => {
          if (dragging) {
            chooseStart(column);
          }
        }}
        onclick={() => {
          chooseStart(column);
        }}
      >
        <span>{column}</span>
      </button>
    {/each}
  </div>

  <div class="section-placement-scale">
    <span>1</span>
    <span>Desktop grid · 12 columns</span>
    <span>12</span>
  </div>

  <p>
    Drag across the grid to position this section.
    Mobile layouts always use a single column.
  </p>
</div>
