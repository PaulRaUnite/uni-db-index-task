import Vue from 'vue';
import App from './App.vue';
import './registerServiceWorker';
import router from './router';
import store from './store';
import Complaints from "@/views/Complaints.vue";
import Login from "@/views/AuthPage.vue";
// @ts-ignore
import Va from 'vue-atlas';
import 'vue-atlas/dist/vue-atlas.css';
import AppLayout from "@/layouts/AppLayout.vue";
import BlankLayout from "@/layouts/BlankLayout.vue";

Vue.config.productionTip = false;

Vue.use(Va, 'en');

Vue.component("app-layout", AppLayout);
Vue.component("blank-layout", BlankLayout);

new Vue({
    router,
    store,
    components: {App, Login, Complaints},
    render: (h) => h(App),
}).$mount('#app');
