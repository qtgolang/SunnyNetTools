export const getElement = (selector: string): HTMLElement | null => window?.document.querySelector(selector);

export const createElement = (tagName: string): HTMLElement => window?.document.createElement(tagName);
