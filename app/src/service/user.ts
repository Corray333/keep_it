import type { User } from "@/entities/user"
import { UserTransport } from "@/transport/user"

const userTransport = new UserTransport()
export class UserService {
    checkUsernameOrEmail = async (checkStr: string) : Promise<User|null>=>{
        return userTransport.FindUser(checkStr)
    }

    checkCode = async (checkStr: string, code: string) : Promise<boolean> =>{
        return userTransport.CheckCode(checkStr, code)
    }
}