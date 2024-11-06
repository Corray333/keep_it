import type { User } from "@/entities/user"
import { api, logError } from "./main"

interface SignUpResponse {
    authorization: string
    user: User
}

export class UserTransport {
    FindUserLogin = async (checkStr: string) : Promise<User|null>=>{
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

    CheckCode = async (checkStr: string, code: string) : Promise<boolean> =>{
        try {
            const response = await api.post('/auth/check-code', {
                checkStr: checkStr,
                code: code
            })
            return response.data
        } catch (error) {
            logError(error)
            return false
        }
    }

    SignUp = async (username: string, password: string, code: string) : Promise<SignUpResponse | null> =>{
        try {
            const response = await api.post('/auth/signup', {
                username: username,
                password: password,
                code: code
            })
            return response.data.authorization, response.data.user
        } catch (error) {
            logError(error)            
            return null
        }
    }

    LogIn = async (username: string, password: string, code: string) : Promise<SignUpResponse | null> =>{
        try {
            const response = await api.post('/auth/login', {
                username: username,
                password: password,
                code: code
            })
            return response.data.authorization, response.data.user
        } catch (error) {
            logError(error)
            return null
        }
    }
}