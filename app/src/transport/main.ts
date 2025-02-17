import { useAccountStore } from '@/stores/account'
import { useComponentsStore } from '@/stores/components'
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
    if (config.useAuth) {
        const token = useAccountStore().authorization   
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`
        }
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
    const componentsStore = useComponentsStore()
    console.log(error)
    if (isAxiosError(error)){
        componentsStore.newError(error.response?.data)
    } else {
        componentsStore.newError()
    }
}

export {api, logError}