import type { User } from "@/entities/user"
import { UserTransport } from "@/transport/user"

const userTransport = new UserTransport()
export class UserService {
    findUserLogin = async (checkStr: string) : Promise<User|null>=>{
        return userTransport.FindUserLogin(checkStr)
    }

    checkCode = async (checkStr: string, code: string) : Promise<boolean> =>{
        return userTransport.CheckCode(checkStr, code)
    }

    signUp = async (username: string, password: string, code: string) : Promise<{authorization: string, user: User} | null> =>{
        return userTransport.SignUp(username, password, code)
    }

    logIn = async (username: string, password: string, code: string) : Promise<{authorization: string, user: User} | null> =>{
        return userTransport.LogIn(username, password, code)
    }
}