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
  const opacity =
    clamp(
      section.surfaceOpacity,
      0,
      100
    ) / 100;

  const blur =
    clamp(
      section.surfaceBlur,
      0,
      40
    );

  return [
    `--section-accent:${section.accent}`,
    `--section-surface-opacity:${opacity}`,
    `--section-surface-blur:${blur}px`,
    `--section-border-width:${
      section.showBorder
        ? '1px'
        : '0px'
    }`
  ].join(';');
}
