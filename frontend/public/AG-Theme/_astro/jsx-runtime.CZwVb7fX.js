import{r as i}from"./index.DHWipiO2.js";var p={exports:{}},o={};/**
 * @license React
 * react-jsx-runtime.production.min.js
 *
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */var l=i,d=Symbol.for("react.element"),m=Symbol.for("react.fragment"),x=Object.prototype.hasOwnProperty,y=l.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED.ReactCurrentOwner,a={key:!0,ref:!0,__self:!0,__source:!0};function _(t,e,u){var r,n={},f=null,s=null;u!==void 0&&(f=""+u),e.key!==void 0&&(f=""+e.key),e.ref!==void 0&&(s=e.ref);for(r in e)x.call(e,r)&&!a.hasOwnProperty(r)&&(n[r]=e[r]);if(t&&t.defaultProps)for(r in e=t.defaultProps,e)n[r]===void 0&&(n[r]=e[r]);return{$$typeof:d,type:t,key:f,ref:s,props:n,_owner:y.current}}o.Fragment=m;o.jsx=_;o.jsxs=_;p.exports=o;var E=p.exports;export{E as j};
