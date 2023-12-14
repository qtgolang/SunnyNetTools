const createStyle = (id: string, styleString: string): void => {
  if (window?.document.getElementById(id)) return;

  const head = window?.document.getElementsByTagName('head')[0];
  const style = window?.document.createElement('style');
  style.setAttribute('id', id);
  style.textContent = styleString;

  head.appendChild(style);
};

export const setFullWidthButtonStyle = async (): Promise<void> => {
  const {fullWidthButton} = await import('./string-styles');

  createStyle('full-width-button', fullWidthButton);
};

export const setDarkThemeStyle = async (): Promise<void> => {
  const {darkTheme} = await import('./string-styles');

  createStyle('dark-theme', darkTheme);
};
