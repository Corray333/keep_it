import type { User } from "@/entities/user"
import { UserTransport } from "@/transport/user"

export class UserService {
    userTransport = new UserTransport()
    findUserLogin = async (checkStr: string) : Promise<User|null>=>{
        return UserTransport.findUserLogin(checkStr)
    }

    signUp = async (username: string, password: string, code: string) : Promise<{authorization: string, user: User} | null> =>{
        return UserTransport.signUp(username, password, code)
    }

    logIn = async (username: string, password: string, code: string) : Promise<{authorization: string, user: User} | null> =>{
        return UserTransport.logIn(username, password, code)
    }

    codeExists = async (username: string, syn: number) : Promise<boolean> =>{
        return UserTransport.codeExists(username, syn)
    }
}