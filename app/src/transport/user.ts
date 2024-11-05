import type { User } from "@/entities/user"
import { useComponentsStore } from "@/stores/components"
import axios from "axios"
import { isAxiosError } from "axios"

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
})


interface SignUpResponse {
    authorization: string
    user: User
}

export class UserTransport {
    FindUser = async (checkStr: string) : Promise<User|null>=>{
        try {
            const response = await api.post('/users/login-find', {
                checkStr: checkStr
            }
        )
            return response.data
        } catch (error) {
            const componentsStore = useComponentsStore()
            console.log(error)
            componentsStore.newError()

            return null
        }
    }

    CheckCode = async (checkStr: string, code: string) : Promise<boolean> =>{
        try {
            const response = await api.post('/user/check-code', {
                checkStr: checkStr,
                code: code
            })
            return response.data
        } catch (error) {
            const componentsStore = useComponentsStore()
            console.log(error)
            componentsStore.newError()

            return false
        }
    }

    SignUp = async (username: string, password: string, code: string) : Promise<SignUpResponse | null> =>{
        try {
            const response = await api.post('/user/signup', {
                username: username,
                password: password,
                code: code
            })
            return response.data.authorization, response.data.user
        } catch (error) {
            const componentsStore = useComponentsStore()
            if (isAxiosError(error)) {
                componentsStore.newError(error.message)
            } else{
                componentsStore.newError()
            }
            
            return null
        }
    }

    LogIn = async (username: string, password: string, code: string) : Promise<SignUpResponse | null> =>{
        try {
            const response = await api.post('/user/login', {
                username: username,
                password: password,
                code: code
            })
            return response.data.authorization, response.data.user
        } catch (error) {
            const componentsStore = useComponentsStore()
            if (isAxiosError(error)) {
                componentsStore.newError(error.message)
            } else{
                componentsStore.newError()
            }
            
            return null
        }
    }
}