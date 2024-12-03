import { CategoryTransport } from "@/transport/category"

export class CategoryService {
    static getCategories = async () =>{
        return CategoryTransport.getCategories()
    }

    static createCategory = async (name: string, parentID: string) =>{
        return CategoryTransport.createCategory(name, parentID)
    }
}