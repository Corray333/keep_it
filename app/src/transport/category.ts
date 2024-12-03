import type { Category } from "@/entities/category"
import { api, logError } from "./main"

export class CategoryTransport {
    static getCategories = async () : Promise<Category[]> =>{
        try {
            const response = await api.get('/categories')
            return response.data
        } catch (error) {
            logError(error)
            return []
        }
    }

    static createCategory = async (name: string, parentID: string) : Promise<Category | null> =>{
        try {
            const response = await api.post('/categories', {
                label: name,
                parentCategoryID: parentID ? parentID : null
            })
            return response.data
        } catch (error) {
            logError(error)
            return null
        }
    }
}