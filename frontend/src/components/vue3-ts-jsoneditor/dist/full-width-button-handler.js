const getElement = (selector) => window == null ? void 0 : window.document.querySelector(selector);
const createElement = (tagName) => window == null ? void 0 : window.document.createElement(tagName);
export { createElement, getElement };
