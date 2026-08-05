<script lang="ts">
  import {
    searchWeatherLocations,
    type WeatherLocation
  } from '../lib/api';

  import type {
    HubletConfig
  } from '../lib/types';

  let {
    config
  }: {
    config: HubletConfig;
  } = $props();

  let query = $state('');
  let results = $state<WeatherLocation[]>([]);
  let searching = $state(false);
  let searchError = $state('');

  const hasLocation = $derived(
    config.modules.weather.latitude !== null &&
    config.modules.weather.longitude !== null &&
    config.modules.weather.location.trim() !== ''
  );

  $effect(() => {
    const term = query.trim();

    if (term.length < 2) {
      results = [];
      searching = false;
      searchError = '';
      return;
    }

    const controller = new AbortController();

    const timer = window.setTimeout(
      async () => {
        searching = true;
        searchError = '';

        try {
          results = await searchWeatherLocations(
            term,
            controller.signal
          );
        } catch (reason) {
          if (controller.signal.aborted) {
            return;
          }

          results = [];
          searchError =
            reason instanceof Error
              ? reason.message
              : String(reason);
        } finally {
          if (!controller.signal.aborted) {
            searching = false;
          }
        }
      },
      350
    );

    return () => {
      window.clearTimeout(timer);
      controller.abort();
    };
  });

  function chooseLocation(
    location: WeatherLocation
  ): void {
    const weather = config.modules.weather;

    weather.provider = 'open-meteo';
    weather.mode = 'current';
    weather.location = location.name;
    weather.country = location.country ?? '';
    weather.admin1 = location.admin1 ?? '';
    weather.latitude = location.latitude;
    weather.longitude = location.longitude;
    weather.enabled = true;

    query = '';
    results = [];
    searchError = '';
  }

  function locationDetails(
    location: WeatherLocation
  ): string {
    return [
      location.admin1,
      location.country
    ].filter(Boolean).join(', ');
  }
</script>

<section class="widgets-panel">
  <header class="settings-panel-header">
    <div class="settings-panel-mark">☀</div>

    <div>
      <p>WIDGETS</p>
      <h2>Weather</h2>
    </div>
  </header>

  <div class="widgets-panel-content">
    <label class="toggle-control weather-toggle">
      <span>
        <strong>Show weather</strong>

        <small>
          Display current conditions in the dashboard header.
        </small>
      </span>

      <input
        bind:checked={config.modules.weather.enabled}
        disabled={!hasLocation}
        type="checkbox"
      />
    </label>

    <div class="widget-settings-group">
      <div class="widget-settings-heading">
        <strong>Weather service</strong>
        <small>Current conditions, refreshed every 10 minutes.</small>
      </div>

      <select
        bind:value={config.modules.weather.provider}
        aria-label="Weather service"
      >
        <option value="open-meteo">Open-Meteo</option>
      </select>

      <small class="weather-units-note">
        Metric units: °C, km/h and mm.
      </small>
    </div>

    <div class="widget-settings-group">
      <label class="widget-settings-heading" for="weather-location-search">
        <strong>Location</strong>
        <small>Search and choose the correct place from the list.</small>
      </label>

      {#if hasLocation}
        <div class="selected-weather-location">
          <span aria-hidden="true">⌖</span>

          <div>
            <strong>{config.modules.weather.location}</strong>
            <small>
              {[config.modules.weather.admin1,
                config.modules.weather.country]
                .filter(Boolean).join(', ')}
            </small>
          </div>

          <span>Selected</span>
        </div>
      {/if}

      <div class="weather-location-search">
        <span aria-hidden="true">⌕</span>

        <input
          id="weather-location-search"
          bind:value={query}
          type="search"
          autocomplete="off"
          placeholder="City or place"
        />

        {#if searching}
          <span class="weather-search-state">…</span>
        {/if}
      </div>

      {#if searchError}
        <p class="weather-search-message error">
          {searchError}
        </p>
      {:else if query.trim().length >= 2 &&
        !searching && results.length === 0}
        <p class="weather-search-message">
          No matching locations found.
        </p>
      {/if}

      {#if results.length > 0}
        <div class="weather-location-results">
          {#each results as location (location.id)}
            <button
              type="button"
              onclick={() => chooseLocation(location)}
            >
              <span aria-hidden="true">⌖</span>

              <span>
                <strong>{location.name}</strong>
                <small>{locationDetails(location)}</small>
              </span>

              <span>Choose</span>
            </button>
          {/each}
        </div>
      {/if}
    </div>

    <a
      class="weather-attribution"
      href="https://open-meteo.com/"
      target="_blank"
      rel="noreferrer"
    >
      Weather data by Open-Meteo
    </a>
  </div>
</section>
