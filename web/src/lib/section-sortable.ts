import Sortable from 'sortablejs';

export type SectionSortableParameters = {
  enabled: boolean;
  onReorder: (orderedSectionIDs: string[]) => void;
};

export function sectionSortable(
  node: HTMLElement,
  initialParameters: SectionSortableParameters
) {
  let parameters = initialParameters;

  const sortable = Sortable.create(node, {
    animation: 190,

    handle: '.canvas-section-drag-handle',
    draggable: '.editor-canvas-section',
    dataIdAttr: 'data-section-id',

    disabled: !parameters.enabled,

    ghostClass: 'section-sortable-ghost',
    chosenClass: 'section-sortable-chosen',
    dragClass: 'section-sortable-drag',

    /*
     * Fallback mode je pouzdaniji za višeredni CSS grid,
     * naročito u Firefox-u.
     */
    forceFallback: true,
    fallbackOnBody: true,
    fallbackTolerance: 5,

    /*
     * Položaj kursora određuje da li section ide pre
     * ili posle section-a iznad kog se nalazi.
     */
    swapThreshold: 0.55,
    invertSwap: true,
    invertedSwapThreshold: 0.55,

    delay: 0,
    delayOnTouchOnly: true,
    touchStartThreshold: 4,

    scroll: true,
    scrollSensitivity: 90,
    scrollSpeed: 14,

    onStart() {
      document.body.classList.add(
        'section-sort-in-progress'
      );
    },

    onEnd() {
      document.body.classList.remove(
        'section-sort-in-progress'
      );

      parameters.onReorder(
        sortable.toArray()
      );
    }
  });

  return {
    update(
      nextParameters: SectionSortableParameters
    ) {
      parameters = nextParameters;

      sortable.option(
        'disabled',
        !parameters.enabled
      );
    },

    destroy() {
      document.body.classList.remove(
        'section-sort-in-progress'
      );

      sortable.destroy();
    }
  };
}
