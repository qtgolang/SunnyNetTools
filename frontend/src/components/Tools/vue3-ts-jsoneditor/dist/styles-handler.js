const createStyle = (id, styleString) => {
  if (window == null ? void 0 : window.document.getElementById(id))
    return;
  const head = window == null ? void 0 : window.document.getElementsByTagName("head")[0];
  const style = window == null ? void 0 : window.document.createElement("style");
  style.setAttribute("id", id);
  style.textContent = styleString;
  head.appendChild(style);
};
const setFullWidthButtonStyle = async () => {
  const { fullWidthButton } = await import("./string-styles.js");
  createStyle("full-width-button", fullWidthButton);
};
const setDarkThemeStyle = async () => {
  const { darkTheme } = await import("./string-styles.js");
  createStyle("dark-theme", darkTheme);
};
export { setDarkThemeStyle, setFullWidthButtonStyle };
