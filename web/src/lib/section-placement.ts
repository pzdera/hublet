import type {
  Section,
  SectionWidth
} from './types';

export function sectionWidthSpan(
  width: SectionWidth
): number {
  switch (width) {
    case 'narrow':
      return 3;

    case 'medium':
      return 4;

    case 'wide':
      return 6;

    case 'extra-wide':
      return 8;

    case 'full':
      return 12;
  }
}

export function maximumStartColumn(
  width: SectionWidth
): number {
  return 13 - sectionWidthSpan(width);
}

export function normalizeStartColumn(
  section: Section
): void {
  if (section.startColumn <= 0) {
    section.startColumn = 0;
    return;
  }

  section.startColumn = Math.min(
    maximumStartColumn(section.width),
    Math.max(1, section.startColumn)
  );
}

export function sectionPlacementStyle(
  section: Section
): string {
  if (section.startColumn <= 0) {
    return '';
  }

  return [
    `grid-column-start:${section.startColumn}`,
    `grid-column-end:span ${sectionWidthSpan(section.width)}`
  ].join(';');
}
