import type { Section } from './types';

import {
  measuredRowSpan,
  normalizeGridLayout
} from './section-placement';

export type SectionMeasureParameters = {
  section: Section;
  sections: Section[];
  desktopOnly?: boolean;
};

export function measureSection(
  node: HTMLElement,
  initialParameters: SectionMeasureParameters
) {
  let parameters = initialParameters;
  let frame = 0;

  function measure() {
    cancelAnimationFrame(frame);

    frame = requestAnimationFrame(() => {
      if (
        parameters.desktopOnly &&
        window.innerWidth < 1200
      ) {
        return;
      }

      const nextSpan = measuredRowSpan(
        node.getBoundingClientRect().height
      );

      if (parameters.section.gridRowSpan === nextSpan) {
        return;
      }

      parameters.section.gridRowSpan = nextSpan;
      normalizeGridLayout(parameters.sections);
    });
  }

  const observer = new ResizeObserver(measure);
  observer.observe(node);
  measure();

  return {
    update(nextParameters: SectionMeasureParameters) {
      parameters = nextParameters;
      measure();
    },

    destroy() {
      cancelAnimationFrame(frame);
      observer.disconnect();
    }
  };
}
