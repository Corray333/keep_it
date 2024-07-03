import { createStore } from 'vuex'

const store = createStore({
    state: {
      AccessToken: "",
    },
    getters: {},
    mutations: {
      setAccess:(state, token)=>{
        state.AccessToken = token
      }
    }
})

export default store