import Vue from 'vue';
import VueRouter from 'vue-router';
import Home from '../views/Home.vue';
import Complaints from "../views/Complaints.vue";
import Login from "../views/Login.vue";

Vue.use(VueRouter);

const routes = [
    {
        path: '/',
        name: 'home',
        component: Home,
    },
    {
        path: '/login',
        name: 'login',
        component: Login,
    },
    {
        path: '/complaints',
        name: 'complaints',
        component: Complaints,
    },
    {
        path: "*",
        redirect: "/",
    }
];

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes,
});
//
// router.beforeEach((to, from, next) => {
//     // redirect to login page if not logged in and trying to access a restricted page
//     const publicPages = ['/login', "/"];
//     const authRequired = !publicPages.includes(to.path);
//     const loggedIn = sessionStorage.getItem('jwt');
//
//     if (authRequired && !loggedIn) {
//         return next('/login');
//     }
//
//     next();
// });

export default router;
