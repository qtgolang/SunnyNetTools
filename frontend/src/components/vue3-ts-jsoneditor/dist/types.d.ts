import type { Validator, OnRenderValue, OnClassName, OnRenderMenu, Mode, JSONValue, JSONPatchDocument } from 'vanilla-jsoneditor';
export interface JSONEditorOptions {
    readOnly?: boolean;
    indentation?: number | string;
    tabSize?: number;
    mode?: Mode;
    mainMenuBar?: boolean;
    navigationBar?: boolean;
    statusBar?: boolean;
    askToFormat?: boolean;
    escapeControlCharacters?: boolean;
    escapeUnicodeCharacters?: boolean;
    flattenColumns?: boolean;
    validator?: Validator;
    queryLanguagesIds?: QueryLanguageId[];
    queryLanguageId?: QueryLanguageId;
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
export declare type QueryLanguageId = 'javascript' | 'lodash' | 'jmespath';
export interface OnTransformArguments {
    operations: JSONPatchDocument;
    json: JSONValue;
    transformedJson: JSONValue;
}
export interface TransformArguments {
    id?: string;
    rootPath?: [];
    onTransform: (args: OnTransformArguments) => void;
    onClose: () => void;
}
//# sourceMappingURL=types.d.ts.map