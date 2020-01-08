import Vue from 'vue';
import Vuex from 'vuex';
import VuexPersistence from 'vuex-persist'
import base64url from "base64url";

Vue.use(Vuex);

interface State {
    token: string | null
    user_id: number | null
    username: string | null
    account_type: number | null
    logged: boolean
    layout: string
    cart: Object
    cart_total: number
}

const vuexLocal = new VuexPersistence<State>({
    storage: window.sessionStorage
});

const store = new Vuex.Store<State>({
    state: {
        token: null,
        user_id: null,
        username: null,
        logged: false,
        layout: 'app-layout',
        account_type: null,
        cart: {},
        cart_total: 0,
    },
    mutations: {
        signin(state: State, payload: { token: string, username: string }) {
            let claims = JSON.parse(base64url.decode(payload.token.split(".")[1]));
            state.user_id = claims.user_id;
            state.token = payload.token;
            state.logged = true;
            state.account_type = claims.account_type;
            state.username = payload.username;
        },
        signout(state: State) {
            state.user_id = null;
            state.token = null;
            state.logged = false;
            state.account_type = null;
            state.username = null;
            state.cart = {};
            state.cart_total = 0;
        },
        add_to_cart(state: State, payload: { id: number, amount: number, price: number, description: String }) {
            let m = state.cart;
            m[payload.id] = payload;
            state.cart = m;
            let total = 0;
            for (let [key, item] of Object.entries(m)) {
                total += item.amount * item.price;
            }
            state.cart_total = total;
        },
        clear_cart(state: State) {
            state.cart = {};
            state.cart_total = 0;
        }
    },
    actions: {},
    modules: {},
    plugins: [vuexLocal.plugin]
});

export default store;