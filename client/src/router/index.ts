import Vue from 'vue';
import VueRouter from 'vue-router';
import Home from '../views/Home.vue';
import Complaints from "../views/Complaints.vue";
import AuthPage from "@/views/AuthPage.vue";
import Inventory from "@/views/Inventory.vue";
import Users from "@/views/Users.vue";
import UserPage from "@/views/UserPage.vue";
import GoodPage from "@/views/GoodPage.vue";
import Test from "@/views/Test.vue";

Vue.use(VueRouter);

const routes = [
    {
        path: '/',
        name: 'home',
        component: Home,
    },
    {
        path: '/auth',
        name: 'auth',
        meta: {layout: "blank"},
        component: AuthPage,
    },
    {
        path: '/complaints',
        name: 'complaints',
        component: Complaints,
    },
    {
        path: '/inventory',
        name: 'inventory',
        component: Inventory,
    },
    {
        path: '/users',
        name: 'users',
        component: Users,
    },
    {
        path: '/user/:username',
        name: 'user',
        component: UserPage,
    },
    {
        path: '/test',
        name: 'test',
        meta: {layout: "test"},
        component: Test,
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
