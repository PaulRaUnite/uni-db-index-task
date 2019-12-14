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
}

const vuexLocal = new VuexPersistence<State>({
    storage: window.sessionStorage
});

export default new Vuex.Store<State>({
    state: {
        token: null,
        user_id: null,
        username: null,
        logged: false,
        layout: 'app-layout',
        account_type: null,
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
        },
    },
    actions: {},
    modules: {},
    plugins: [vuexLocal.plugin]
});
