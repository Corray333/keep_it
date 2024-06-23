import axios from 'axios'
import {useRouter} from 'vue-router'

function getCookie(name) {
    const value = `; ${document.cookie}`
    const parts = value.split(`; ${name}=`)
    if (parts.length === 2) return parts.pop().split(';').shift()
}

const refreshTokens = async () => {
    try {
        let { data } = await axios.get( `/api/users/refresh`, {
            headers: {
                'Refresh': localStorage.getItem('Refresh'),
            }
        })

        localStorage.setItem('Authorization', data.authorization)
        localStorage.setItem('Refresh', data.refresh)
    } catch (error) {
        alert('Error refreshing tokens')
        const router = useRouter()
        router.push('/login')
    }
}

export { getCookie, refreshTokens };