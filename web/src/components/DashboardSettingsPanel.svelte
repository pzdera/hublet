<script lang="ts">
  import type {
    HubletConfig
  } from '../lib/types';

  import AppearancePanel from './AppearancePanel.svelte';
  import WidgetsPanel from './WidgetsPanel.svelte';
  import ShortcutsPanel from './ShortcutsPanel.svelte';

  let {
    config
  }: {
    config: HubletConfig;
  } = $props();

  let activeTab = $state<
    'appearance' | 'widgets' | 'shortcuts'
  >('appearance');
</script>

<aside class="dashboard-settings-panel">
  <nav
    class="dashboard-settings-tabs"
    aria-label="Dashboard settings"
  >
    <button
      class:active={activeTab === 'widgets'}
      type="button"
      onclick={() => {
        activeTab = 'widgets';
      }}
    >
      Widgets
      {#if config.modules.weather.enabled}
        <span>1</span>
      {/if}
    </button>

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
    {:else if activeTab === 'widgets'}
      <WidgetsPanel {config} />
    {:else}
      <ShortcutsPanel {config} />
    {/if}
  </div>
</aside>
