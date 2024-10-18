import './assets/main.css';
import App from './App.vue';
import router from './router';
import { createApp } from 'vue';
import Antd from 'ant-design-vue';
import VueCookies from 'vue-cookies';
import 'ant-design-vue/dist/reset.css';

const app = createApp(App);

app.use(router);
app.use(Antd);
app.use(VueCookies, {
  secure: true,
  expires: '20d',
  sameSite: 'Strict',
  path: '/BASE_URI/ui/',
  domain: window.location.hostname,
});

app.mount('#app');
