<script lang="ts">
  import type {
    HubletConfig
  } from '../lib/types';

  import AppearancePanel from './AppearancePanel.svelte';
  import ShortcutsPanel from './ShortcutsPanel.svelte';

  let {
    config
  }: {
    config: HubletConfig;
  } = $props();

  let activeTab = $state<
    'appearance' | 'shortcuts'
  >('appearance');
</script>

<aside class="dashboard-settings-panel">
  <nav
    class="dashboard-settings-tabs"
    aria-label="Dashboard settings"
  >
    <button
      class:active={activeTab === 'appearance'}
      type="button"
      onclick={() => {
        activeTab = 'appearance';
      }}
    >
      Appearance
    </button>

    <button
      class:active={activeTab === 'shortcuts'}
      type="button"
      onclick={() => {
        activeTab = 'shortcuts';
      }}
    >
      Shortcuts
      <span>{config.shortcuts.length}</span>
    </button>
  </nav>

  <div class="dashboard-settings-content">
    {#if activeTab === 'appearance'}
      <AppearancePanel {config} />
    {:else}
      <ShortcutsPanel {config} />
    {/if}
  </div>
</aside>
