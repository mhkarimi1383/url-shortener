<script setup lang="ts">
import { inject } from 'vue';
import { loginStateCookie, loginInfoCookie, setToken } from '@/lib/api';
import type { VueCookies } from 'vue-cookies';
import router from '@/router';
import { message } from 'ant-design-vue';

const $cookies = inject<VueCookies>('$cookies');

const finish = message.loading('Logging out');
$cookies?.remove(loginInfoCookie);
$cookies?.remove(loginStateCookie);
setToken(null);
setTimeout(finish, 1000);
message.success('Logged out');
router.push('/user/login').finally(() => location.reload());
</script>
