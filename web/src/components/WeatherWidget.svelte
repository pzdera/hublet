<script lang="ts">
  import {
    onMount
  } from 'svelte';

  import {
    loadCurrentWeather,
    type CurrentWeather
  } from '../lib/api';

  let current = $state<CurrentWeather | null>(null);
  let failed = $state(false);

  function weatherSymbol(code: number): string {
    if (code === 0) return '☀';
    if (code <= 3) return '☁';
    if (code <= 48) return '≋';
    if (code <= 67) return '☂';
    if (code <= 77) return '❄';
    if (code <= 82) return '☂';
    if (code <= 86) return '❄';
    return 'ϟ';
  }

  function windDirection(degrees: number): string {
    const directions = [
      'N', 'NE', 'E', 'SE',
      'S', 'SW', 'W', 'NW'
    ];

    return directions[
      Math.round(degrees / 45) % 8
    ];
  }

  async function refresh(): Promise<void> {
    try {
      current = await loadCurrentWeather();
      failed = false;
    } catch {
      failed = true;
    }
  }

  onMount(() => {
    void refresh();

    const interval = window.setInterval(
      () => void refresh(),
      15 * 60 * 1000
    );

    return () => window.clearInterval(interval);
  });
</script>

{#if current}
  <div
    class="weather-widget"
    title={`Feels like ${Math.round(current.apparentTemperature)} °C · Gusts ${Math.round(current.windGusts)} km/h · Data by Open-Meteo`}
  >
    <span class="weather-widget-symbol" aria-hidden="true">
      {weatherSymbol(current.weatherCode)}
    </span>

    <div class="weather-widget-primary">
      <strong>{Math.round(current.temperature)}°</strong>
      <span>{current.location}</span>
    </div>

    <div class="weather-widget-meta">
      <span title="Wind">
        <b>Wind</b>

        <span>
          {windDirection(current.windDirection)}
          {Math.round(current.windSpeed)} km/h
        </span>
      </span>

      <span title="Current precipitation">
        <b>Rain</b>

        <span>
          {current.precipitation.toFixed(1)} mm
        </span>
      </span>
    </div>
  </div>
{:else if failed}
  <button
    class="weather-widget weather-widget-error"
    type="button"
    title="Weather is temporarily unavailable"
    onclick={() => void refresh()}
  >
    <span aria-hidden="true">↻</span>
    <span>Weather</span>
  </button>
{:else}
  <div
    class="weather-widget weather-widget-loading"
    aria-label="Loading weather"
  >
    <span></span>
    <span></span>
  </div>
{/if}
