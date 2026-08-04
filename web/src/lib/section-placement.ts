import type {
  Section,
  SectionWidth
} from './types';

export const dashboardGridColumns = 24;
export const dashboardGridRowHeight = 8;
export const dashboardGridGap = 8;
export const dashboardGridRowStep =
  dashboardGridRowHeight + dashboardGridGap;

export function legacyWidthSpan(
  width: SectionWidth
): number {
  switch (width) {
    case 'narrow':
      return 6;
    case 'medium':
      return 8;
    case 'wide':
      return 12;
    case 'extra-wide':
      return 16;
    case 'full':
      return 24;
  }
}

export function sectionColumnSpan(
  section: Section
): number {
  const span = section.gridColumnSpan ||
    legacyWidthSpan(section.width);

  return Math.max(
    4,
    Math.min(dashboardGridColumns, span)
  );
}

export function sectionRowSpan(
  section: Section
): number {
  return Math.max(1, section.gridRowSpan || 1);
}

export function measuredRowSpan(
  height: number
): number {
  return Math.max(
    1,
    Math.ceil(
      (height + dashboardGridGap) /
      dashboardGridRowStep
    )
  );
}

export function placementFits(
  sections: Section[],
  movingSectionId: string,
  row: number,
  column: number
): boolean {
  const moving = sections.find(
    (section) => section.id === movingSectionId
  );

  if (!moving || row < 1 || column < 1) {
    return false;
  }

  const columnSpan = sectionColumnSpan(moving);
  const rowSpan = sectionRowSpan(moving);
  const endColumn = column + columnSpan - 1;
  const endRow = row + rowSpan - 1;

  if (endColumn > dashboardGridColumns) {
    return false;
  }

  return !sections.some((section) => {
    if (section.id === movingSectionId) {
      return false;
    }

    const sectionStartColumn = section.gridColumn;
    const sectionEndColumn =
      sectionStartColumn + sectionColumnSpan(section) - 1;
    const sectionStartRow = section.gridRow;
    const sectionEndRow =
      sectionStartRow + sectionRowSpan(section) - 1;

    const columnsOverlap =
      column <= sectionEndColumn &&
      endColumn >= sectionStartColumn;
    const rowsOverlap =
      row <= sectionEndRow &&
      endRow >= sectionStartRow;

    return columnsOverlap && rowsOverlap;
  });
}

export function placementCollisions(
  sections: Section[],
  movingSectionId: string,
  row: number,
  column: number
): Section[] {
  const moving = sections.find(
    (section) => section.id === movingSectionId
  );

  if (!moving) {
    return [];
  }

  const movingEndColumn =
    column + sectionColumnSpan(moving) - 1;
  const movingEndRow =
    row + sectionRowSpan(moving) - 1;

  return sections.filter((section) => {
    if (section.id === movingSectionId) {
      return false;
    }

    const sectionEndColumn =
      section.gridColumn + sectionColumnSpan(section) - 1;
    const sectionEndRow =
      section.gridRow + sectionRowSpan(section) - 1;

    return (
      column <= sectionEndColumn &&
      movingEndColumn >= section.gridColumn &&
      row <= sectionEndRow &&
      movingEndRow >= section.gridRow
    );
  });
}

export type SectionSwapPlan = {
  movingRow: number;
  movingColumn: number;
  targetRow: number;
  targetColumn: number;
};

type GridPosition = {
  row: number;
  column: number;
  distance: number;
};

function swapDestinationPositions(
  sections: Section[],
  incoming: Section,
  destination: Section,
  preferredRow = destination.gridRow,
  preferredColumn = destination.gridColumn
): GridPosition[] {
  const incomingColumnSpan = sectionColumnSpan(incoming);
  const incomingRowSpan = sectionRowSpan(incoming);
  const destinationColumnEnd =
    destination.gridColumn +
    sectionColumnSpan(destination) -
    1;
  const destinationRowEnd =
    destination.gridRow + sectionRowSpan(destination) - 1;
  const firstColumn = Math.max(
    1,
    destination.gridColumn - incomingColumnSpan + 1
  );
  const lastColumn = Math.min(
    dashboardGridColumns - incomingColumnSpan + 1,
    destinationColumnEnd
  );
  const firstRow = Math.max(
    1,
    destination.gridRow - incomingRowSpan + 1
  );
  const lastRow = destinationRowEnd;
  const positions: GridPosition[] = [];

  for (let row = firstRow; row <= lastRow; row += 1) {
    for (
      let column = firstColumn;
      column <= lastColumn;
      column += 1
    ) {
      if (
        placementFits(
          sections,
          incoming.id,
          row,
          column
        )
      ) {
        positions.push({
          row,
          column,
          distance:
            Math.abs(row - preferredRow) +
            Math.abs(column - preferredColumn)
        });
      }
    }
  }

  return positions.sort(
    (left, right) => left.distance - right.distance
  );
}

function positionsOverlap(
  first: GridPosition,
  firstSection: Section,
  second: GridPosition,
  secondSection: Section
): boolean {
  const firstColumnEnd =
    first.column + sectionColumnSpan(firstSection) - 1;
  const firstRowEnd =
    first.row + sectionRowSpan(firstSection) - 1;
  const secondColumnEnd =
    second.column + sectionColumnSpan(secondSection) - 1;
  const secondRowEnd =
    second.row + sectionRowSpan(secondSection) - 1;

  return (
    first.column <= secondColumnEnd &&
    firstColumnEnd >= second.column &&
    first.row <= secondRowEnd &&
    firstRowEnd >= second.row
  );
}

export function sectionSwapPlan(
  sections: Section[],
  movingSectionId: string,
  targetSectionId: string,
  preferredMovingRow?: number,
  preferredMovingColumn?: number
): SectionSwapPlan | null {
  const moving = sections.find(
    (section) => section.id === movingSectionId
  );
  const target = sections.find(
    (section) => section.id === targetSectionId
  );

  if (!moving || !target || moving.id === target.id) {
    return null;
  }

  const sectionsWithoutSwapPair = sections.filter(
    (section) =>
      section.id !== moving.id && section.id !== target.id
  );
  const movingPositions = swapDestinationPositions(
    [...sectionsWithoutSwapPair, moving],
    moving,
    target,
    preferredMovingRow,
    preferredMovingColumn
  );
  const targetPositions = swapDestinationPositions(
    [...sectionsWithoutSwapPair, target],
    target,
    moving
  );
  let bestPlan: SectionSwapPlan | null = null;
  let bestDistance = Number.POSITIVE_INFINITY;

  for (const movingPosition of movingPositions) {
    for (const targetPosition of targetPositions) {
      const distance =
        movingPosition.distance + targetPosition.distance;

      if (distance >= bestDistance) {
        continue;
      }

      if (
        positionsOverlap(
          movingPosition,
          moving,
          targetPosition,
          target
        )
      ) {
        continue;
      }

      bestDistance = distance;
      bestPlan = {
        movingRow: movingPosition.row,
        movingColumn: movingPosition.column,
        targetRow: targetPosition.row,
        targetColumn: targetPosition.column
      };
    }
  }

  return bestPlan;
}

export function maximumGridBottom(
  sections: Section[]
): number {
  return Math.max(
    1,
    ...sections.map(
      (section) =>
        section.gridRow + sectionRowSpan(section) - 1
    )
  );
}

export function firstAvailablePlacement(
  sections: Section[],
  movingSectionId: string
): { row: number; column: number } {
  const moving = sections.find(
    (section) => section.id === movingSectionId
  );

  if (!moving) {
    return { row: 1, column: 1 };
  }

  const columnSpan = sectionColumnSpan(moving);
  const maximumColumn =
    dashboardGridColumns - columnSpan + 1;
  const maximumRow = maximumGridBottom(sections) + 1;

  for (let row = 1; row <= maximumRow; row += 1) {
    for (
      let column = 1;
      column <= maximumColumn;
      column += 1
    ) {
      if (
        placementFits(
          sections,
          movingSectionId,
          row,
          column
        )
      ) {
        return { row, column };
      }
    }
  }

  return {
    row: maximumRow,
    column: 1
  };
}

export function normalizeGridLayout(
  sections: Section[]
): void {
  const placed: Section[] = [];

  for (const section of sections) {
    section.gridColumnSpan = sectionColumnSpan(section);
    section.gridRowSpan = sectionRowSpan(section);

    const candidateSections = [
      ...placed,
      section
    ];

    const validCurrentPlacement =
      Number.isInteger(section.gridRow) &&
      Number.isInteger(section.gridColumn) &&
      placementFits(
        candidateSections,
        section.id,
        section.gridRow,
        section.gridColumn
      );

    if (!validCurrentPlacement) {
      const placement = firstAvailablePlacement(
        candidateSections,
        section.id
      );

      section.gridRow = placement.row;
      section.gridColumn = placement.column;
    }

    placed.push(section);
  }
}

export function sectionPlacementStyle(
  section: Section
): string {
  const columnSpan = sectionColumnSpan(section);
  const rowSpan = sectionRowSpan(section);
  const tabletSpan = Math.min(
    12,
    Math.max(1, Math.ceil(columnSpan / 2))
  );

  return [
    `grid-row:${section.gridRow} / span ${rowSpan}`,
    `grid-column:${section.gridColumn} / span ${columnSpan}`,
    `--tablet-column-span:${tabletSpan}`
  ].join(';');
}
