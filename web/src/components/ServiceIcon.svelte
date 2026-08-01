<script lang="ts">
  import type {
    Item
  } from '../lib/types';

  import {
    dashboardIconURL,
    iconCandidates,
    itemInitial,
    localIconURL,
    serviceFaviconURL
  } from '../lib/icons';

  let {
    item,
    decorative = true
  }: {
    item: Item;
    decorative?: boolean;
  } = $props();

  const sources = $derived.by(
    (): string[] => {
      if (
        item.icon.type === 'none'
      ) {
        return [];
      }

      const result: string[] = [];

      const local =
        localIconURL(item);

      if (local) {
        result.push(local);
      }

      if (
        item.icon.type === 'auto'
      ) {
        for (
          const candidate
          of iconCandidates(item)
        ) {
          result.push(
            dashboardIconURL(
              candidate
            )
          );
        }

        const favicon =
          serviceFaviconURL(
            item.url
          );

        if (favicon) {
          result.push(
            favicon
          );
        }
      }

      return [
        ...new Set(result)
      ];
    }
  );

  let sourceIndex =
    $state(0);

  $effect(() => {
    item.name;
    item.url;
    item.icon.type;
    item.icon.value;

    sourceIndex = 0;
  });

  const currentSource =
    $derived(
      sources[sourceIndex] ??
      ''
    );

  function handleError() {
    sourceIndex += 1;
  }
</script>

<span
  class:initial-only={!currentSource}
  class="service-icon-content"
>
  {#if currentSource}
    <img
      src={currentSource}
      alt={decorative
        ? ''
        : item.name}
      onerror={handleError}
    />
  {:else}
    <span>
      {itemInitial(item)}
    </span>
  {/if}
</span>
