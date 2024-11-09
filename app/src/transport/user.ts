import type { User } from "@/entities/user"
import { api, logError } from "./main"

interface SignUpResponse {
    authorization: string
    user: User
}

export class UserTransport {
    findUserLogin = async (checkStr: string) : Promise<User|null>=>{
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

    signUp = async (username: string, password: string, code: string) : Promise<SignUpResponse | null> =>{
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

    logIn = async (username: string, password: string, code: string) : Promise<SignUpResponse | null> =>{
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

    codeExists = async (username: string, syn: number) : Promise<boolean> =>{
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