<script lang="ts">
  import type {
    HubletConfig
  } from '../lib/types';

  let {
    dashboard,
    preview = false
  }: {
    dashboard: HubletConfig['dashboard'];
    preview?: boolean;
  } = $props();

  const iconURL = $derived(
    dashboard.icon.type === 'local' &&
    dashboard.icon.value.trim()
      ? `/icons/${
          encodeURIComponent(
            dashboard.icon.value.trim()
          )
        }`
      : ''
  );

  const initial = $derived(
    dashboard.icon.value
      .trim()
      .slice(0, 2)
      .toUpperCase() ||
    dashboard.title
      .trim()
      .slice(0, 1)
      .toUpperCase() ||
    'H'
  );
</script>

<div
  class:preview
  class={`dashboard-brand dashboard-brand-icon-${dashboard.iconSize}`}
>
  {#if dashboard.icon.type !== 'none'}
    <span class="dashboard-brand-icon">
      {#if iconURL}
        <img
          src={iconURL}
          alt=""
        />
      {:else}
        <span>{initial}</span>
      {/if}
    </span>
  {/if}

  <div class="dashboard-brand-copy">
    {#if (
      dashboard.descriptionVisible &&
      dashboard.description.trim()
    )}
      <p>
        {dashboard.description}
      </p>
    {/if}

    <h1>
      {dashboard.title}
    </h1>
  </div>
</div>
