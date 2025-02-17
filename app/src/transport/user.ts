import type { User } from "@/entities/user"
import { api, logError } from "./main"
import { useAccountStore } from "@/stores/account"
import axios from "axios"

interface SignUpResponse {
    authorization: string
    user: User
}

export class UserTransport {
    static findUserLogin = async (checkStr: string) : Promise<User|null>=>{
        try {
            const response = await api.post('/users/login-find', {
                checkStr: checkStr
            }
        )
            return response.data
        } catch (error) {
            logError(error)
            return null
        }
    }

    static signUp = async (username: string, password: string, code: string) : Promise<SignUpResponse | null> =>{
        try {
            const {data} = await api.post('/auth/signup', {
                username: username,
                password: password,
                code: code
            })
            return {
                authorization: data.authorization,
                user: data.user
            }
        } catch (error) {
            logError(error)            
            return null
        }
    }

    static logIn = async (username: string, password: string, code: string) : Promise<SignUpResponse | null> =>{
        try {
            const {data} = await api.post('/auth/login', {
                username: username,
                password: password,
                code: code
            })
            return {
                authorization: data.authorization,
                user: data.user
            }
        } catch (error) {
            logError(error)
            return null
        }
    }

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

    static codeExists = async (username: string, syn: number) : Promise<boolean> =>{
        try {
            const response = await api.post('/auth/code-exists', {
                username: username,
                syn: syn
            })
            return response.data.exists
        } catch (error) {
            logError(error)
            return false
        }
    }
}