<style lang="scss">
$titlebar-height: 28px;

.titlebar {
  position: relative;
  flex-grow: 0;
  flex-shrink: 0;
  width: 100%;
  display: flex;
  flex-direction: row;
  top: 0;
  left: 0;
  right: 0;
  z-index: 99999;
  height: $titlebar-height;

  -webkit-app-region: drag;

  &.titlebar-style-dark {
    color: #fff;
    background: #2d3135;
  }

  &.titlebar-style-light {
    color: #2c2c2c;
    background: #efefef;
  }

  .titlebar-resize-handle {
    position: absolute;
    top: 0;
    left: 0;
    -webkit-app-region: no-drag;

    &.top {
      width: 100%;
      height: 3px;
    }

    &.right {
      left: auto;
      right: 0;
      width: 3px;
      height: $titlebar-height;
    }

    &.left {
      width: 3px;
      height: $titlebar-height;
    }
  }

  .titlebar-header {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .titlebar-icon, .titlebar-name {
    display: flex;
    align-content: center;
    align-self: center;
    flex-grow: 0;
    flex-shrink: 0;
    font-size: 14px;
    line-height: $titlebar-height;
    padding: 0 12px;
    height: $titlebar-height;

    > svg, > img {
      display: block;
      align-content: center;
      align-self: center;
      width: auto;
      height: 16px;
    }
  }

  .titlebar-icon ~ .titlebar-name {
    padding-left: 0;
  }

  &.titlebar-platform-darwin {
    .titlebar-header {
      width: 100%;
      text-align: center;
      position: absolute;
    }
  }

  .titlebar-menu {
    display: flex;
    -webkit-app-region: no-drag;

    .titlebar-menu-item {
      min-width: 0;
      position: relative;
      cursor: pointer;

      button {
        border: none;
        box-shadow: none;
        background: transparent;
        height: 100%;
        width: 100%;

        position: relative;
        display: flex;
        flex-direction: row;
        align-items: center;
        margin: 0;
        color: currentColor;
        font-size: 13px;
        padding: 0 10px;
        outline: none;

        &:hover {
          background-color: rgba(0, 0, 0, 0.2);
        }
      }
    }
  }

  .titlebar-buttons {
    display: flex;
    flex-direction: row;
    flex-grow: 0;
    flex-shrink: 0;
    margin-left: auto;

    button {
      -webkit-app-region: no-drag;
      display: inline-block;
      position: relative;
      width: 45px;
      height: 100%;
      padding: 0;
      margin: 0;
      overflow: hidden;
      border: none;
      box-shadow: none;
      border-radius: 0;
      color: currentColor;
      background-color: transparent;
      line-height: 10px;
      outline: none;

      svg {
        fill: currentColor;
      }

      &:hover {
        background-color: rgba(0, 0, 0, 0.2);
        color: currentColor;
      }

      &.close:hover {
        background-color: #e81123;
        color: #fff;
      }
    }
  }

  .titlebar-buttons-osx {
    display: flex;
    box-sizing: border-box;
    padding: 10px;
    width: 70px;
    flex-direction: row;
    -webkit-box-pack: justify;
    justify-content: space-between;
    align-items: center;

    .macButton {
      -webkit-app-region: no-drag;
      -webkit-box-sizing: border-box;
      border-radius: 50%;
      box-sizing: border-box;
      height: 12px;
      width: 12px;

      background-color: #dcdcdc;
      border-color: #d1d1d1;

      &.macButtonClose {
        background-color: #fc615d;
      }

      &.macButtonMinimize {
        background-color: #fdbc40;
      }

      &.macButtonMaximize {
        background-color: #34c749;
      }

      svg {
        display: block;
        visibility: hidden;
      }
    }

    &:hover {
      .macButton {
        svg {
          visibility: visible;
        }
      }
    }
  }
}
</style>
<template>
  <div class="titlebar" :class="[styleClass, stylePlatform]">
    <div>
      <div class="titlebar-resize-handle top"></div>
      <div class="titlebar-resize-handle right"></div>
      <div class="titlebar-resize-handle left"></div>

    </div>

    <div v-if="platform === 'darwin'" class="titlebar-buttons-osx">
      <div class="macButton macButtonClose" @click="onClose()" v-if="isClosable">
        <svg name="TitleBarCloseMac" width="12" height="12" viewBox="0 0 12 12">
          <path stroke="#4c0000" fill="none"
                d="M8.5,3.5 L6,6 L3.5,3.5 L6,6 L3.5,8.5 L6,6 L8.5,8.5 L6,6 L8.5,3.5 Z"></path>
        </svg>
      </div>
      <div class="macButton macButtonMinimize" @click="onMinimize()" v-if="isMinimizable">
        <svg name="TitleBarMinimizeMac" width="12" height="12" viewBox="0 0 12 12">
          <rect fill="#975500" width="8" height="2" x="2" y="5" fill-rule="evenodd"></rect>
        </svg>
      </div>
      <div class="macButton macButtonMaximize" @click="onMaximize()" v-if="isMaximizable">
        <svg name="TitleBarMaximizeMac" width="12" height="12" viewBox="0 0 12 12">
          <g fill="#006500" fill-rule="evenodd">
            <path
                d="M5,3 C5,3 5,6.1325704 5,6.48601043 C5,6.83945045 5.18485201,7 5.49021559,7 L9,7 L9,6 L8,6 L8,5 L7,5 L7,4 L6,4 L6,3 L5,3 Z"
                transform="rotate(180 7 5)"></path>
            <path
                d="M3,5 C3,5 3,8.1325704 3,8.48601043 C3,8.83945045 3.18485201,9 3.49021559,9 L7,9 L7,8 L6,8 L6,7 L5,7 L5,6 L4,6 L4,5 L3,5 Z"></path>
          </g>
        </svg>
      </div>
    </div>

    <div class="titlebar-header" :style="getHeaderStyle">
      <div class="titlebar-icon" v-if="showIcon" style="width: 30px">
        <slot name="icon"></slot>
      </div>
      <div class="titlebar-name" v-if="showTitle" :style="getTitleStyle">
        <slot name="title"></slot>
      </div>
    </div>

    <div class="titlebar-menu" v-if="platform !== 'darwin'">
      <div class="titlebar-menu-item" v-for="(item, index) in menu" :key="index">
        <button @click="item.click()">
          {{ item.label }}
        </button>
      </div>
    </div>

    <div class="titlebar-buttons" v-if="platform !== 'darwin'">
      <button aria-label="minimize" title="Minimize" tabindex="-1" @click="onMinimize()" v-if="isMinimizable">
        <svg aria-hidden="true" version="1.1" width="10" height="10">
          <path d="M 0,5 10,5 10,6 0,6 Z"></path>
        </svg>
      </button>

      <button aria-label="maximize" title="Maximize" tabindex="-1" @click="onMaximize()" v-if="isMaximizable">
        <svg aria-hidden="true" version="1.1" width="10" height="10">
          <path d="M 0,0 0,10 10,10 10,0 Z M 1,1 9,1 9,9 1,9 Z"></path>
        </svg>
      </button>

      <button aria-label="close" title="Close" tabindex="-1" class="close" @click="onClose()" v-if="isClosable">
        <svg aria-hidden="true" version="1.1" width="10" height="10">
          <path
              d="M 0,0 0,0.7 4.3,5 0,9.3 0,10 0.7,10 5,5.7 9.3,10 10,10 10,9.3 5.7,5 10,0.7 10,0 9.3,0 5,4.3 0.7,0 Z"></path>
        </svg>
      </button>
    </div>
  </div>
</template>
<script>
export default {
  name: "v-titlebar",
  props: {
    theme: {
      type: String,
      default: "light"
    },
    platform: {
      type: String,
      required: true
    },
    menu: {
      type: Array,
      default: function () {
        return [];
      }
    },
    isMinimizable: {
      type: Boolean,
      default: true
    },
    isMaximizable: {
      type: Boolean,
      default: true
    },
    isClosable: {
      type: Boolean,
      default: true
    },
    onMinimize: {
      type: Function
    },
    onMaximize: {
      type: Function
    },
    onClose: {
      type: Function
    },
    showIcon: {
      type: Boolean,
      default: true
    },
    showTitle: {
      type: Boolean,
      default: true
    }
  },
  computed: {
    styleClass() {
      return `titlebar-style-${this.theme}`
    },
    stylePlatform() {
      return `titlebar-platform-${this.platform}`;
    },
    getTitleStyle() {
      return this.showIcon ? 'width: calc(100% - 30px)' : "width: 100% ";
    },
    getHeaderStyle() {
      return this.platform==="darwin" ? 'width: calc(100% - 70px);left: 70px;' : "width: 100% ";
    }
  }
}
</script>