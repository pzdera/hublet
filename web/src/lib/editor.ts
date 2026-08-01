export type DashboardSelection = {
  type: 'dashboard';
};

export type SectionSelection = {
  type: 'section';
  sectionId: string;
};

export type ItemSelection = {
  type: 'item';
  sectionId: string;
  itemId: string;
};

export type EditorSelection =
  | DashboardSelection
  | SectionSelection
  | ItemSelection;
