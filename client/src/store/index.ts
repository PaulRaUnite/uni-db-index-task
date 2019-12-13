import Vue from 'vue';
import Vuex from 'vuex';
import VuexPersistence from 'vuex-persist'
import base64url from "base64url";

Vue.use(Vuex);

interface State {
    token: string | null
    user_id: number | null
}

const vuexLocal = new VuexPersistence<State>({
    storage: window.sessionStorage
});

export default new Vuex.Store<State>({
    state: {
        token: null,
        user_id: null
    },
    mutations: {
        login(state: State, new_token: string) {
            state.user_id = JSON.parse(base64url.decode(new_token.split(".")[1])).user_id;
            state.token = new_token
        }
    },
    actions: {},
    modules: {},
    plugins: [vuexLocal.plugin]
});
