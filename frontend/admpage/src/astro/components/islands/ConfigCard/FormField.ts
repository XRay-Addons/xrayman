export interface FormFieldElement extends HTMLElement {
  updateData(data: Record<string, any>): void;
  refreshData(data: Record<string, any>): void;
}
