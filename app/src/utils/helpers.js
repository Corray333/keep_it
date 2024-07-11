import axios from 'axios'

function getCookie(name) {
    const value = `; ${document.cookie}`
    const parts = value.split(`; ${name}=`)
    if (parts.length === 2) return parts.pop().split(';').shift()
}

const refreshTokens = async (store) => {
    try {
        let { data } = await axios.get( `${import.meta.env.VITE_API_URL}/users/refresh`, {
            withCredentials: true
        })
        store.commit("setAccess", data.authorization)
    } catch (error) {
        console.log(error)
        alert(error)
    }
}

export { getCookie, refreshTokens };