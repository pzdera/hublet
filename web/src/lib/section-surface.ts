import type {
  Section
} from './types';

function clamp(
  value: number,
  minimum: number,
  maximum: number
): number {
  return Math.min(
    maximum,
    Math.max(
      minimum,
      value
    )
  );
}

export function sectionSurfaceStyle(
  section: Section
): string {
  const configuredOpacity =
    clamp(
      section.surfaceOpacity,
      0,
      100
    ) / 100;

  const opacity =
    section.surface === 'solid'
      ? 1
      : section.surface === 'none'
        ? 0
        : configuredOpacity;

  const blur =
    clamp(
      section.surfaceBlur,
      0,
      40
    );

  const effectiveBlur =
    section.surface === 'glass'
      ? blur
      : 0;

  return [
    `--section-accent:${section.accent}`,
    `--section-surface-opacity:${opacity}`,
    `--section-opacity-percent:${opacity * 100}%`,
    `--section-surface-blur:${blur}px`,
    `--section-effective-blur:${effectiveBlur}px`,
    `--section-border-width:${
      section.showBorder
        ? '1px'
        : '0px'
    }`
  ].join(';');
}
