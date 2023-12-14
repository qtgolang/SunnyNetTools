import type {
  Validator,
  QueryLanguage,
  OnRenderValue,
  OnClassName,
  OnRenderMenu,
  Mode,
  JSONValue,
} from 'vanilla-jsoneditor';
export interface JSONEditorOptions {
  readOnly?: boolean;
  indentation?: number | string;
  tabSize?: number;
  mode?: Mode;
  mainMenuBar?: boolean;
  navigationBar?: boolean;
  statusBar?: boolean;
  escapeControlCharacters?: boolean;
  escapeUnicodeCharacters?: boolean;
  validator?: Validator;
  queryLanguages?: QueryLanguage[];
  queryLanguageId?: string;
  onRenderValue?: OnRenderValue;
  onClassName?: OnClassName;
  onRenderMenu?: OnRenderMenu;
  height?: string | number;
  fullWidthButton?: boolean;
  darkTheme?: boolean;
}
export declare type TextContent = {
  json?: undefined;
  text: string;
};
export declare type JSONContent = {
  json: JSONValue;
  text?: undefined;
};
export declare type Content = JSONContent | TextContent;
export declare type Path = Array<string | number | symbol>;
//# sourceMappingURL=types.d.ts.map
