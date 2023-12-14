import {createApp} from 'vue';
import App from './App.vue';
import JsonEditorPlugin from '@/JsonEditorPlugin';

const app = createApp(App);

app.use(JsonEditorPlugin, {
  componentName: 'JsonEditor',
  options: {},
});
app.mount('#app');
