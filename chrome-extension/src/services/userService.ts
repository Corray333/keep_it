import { useAccountStore } from '@/stores/account'
import axios, { type InternalAxiosRequestConfig } from 'axios'
import { isAxiosError } from 'axios'
import { useRouter } from 'vue-router'

interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
    useAuth?: boolean
}

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true,
})

api.interceptors.request.use((config: CustomAxiosRequestConfig) => {
    console.log(config.useAuth)
    if (config.useAuth) {
        const token = useAccountStore().authorization   
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`
        }
        console.log('Authorization header added:', config.headers['Authorization'])
    }
    return config
})

api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;
        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;
            
            try {
                const { data } = await axios.post(`${import.meta.env.VITE_API_URL}/auth/renew-tokens`, {}, {
                    withCredentials: true
                })
                useAccountStore().setAuthorization(data.authorization)
                originalRequest.headers['Authorization'] = `${data.authorization}`
                return api(originalRequest)
            } catch {
                const router = useRouter()
                router.push('/login')
            }
        }
        return Promise.reject(error);
    }
)

const logError = (error: unknown)=>{
    console.log(error)
}

export class UserTransport {

    static renewTokens = async ()=>{
        try {
            const { data } = await axios.post(`${import.meta.env.VITE_API_URL}/auth/renew-tokens`, {}, {
                withCredentials: true
            })
            useAccountStore().setAuthorization(data.authorization)
        } catch(error){
            console.log(error)
        }
    }
}

export {api, logError}