import type { Note } from "@/entities/note"
import { api, logError } from "./main"


export class NoteTransport {
    getNotes = async (offset: number) : Promise<Note[]>=>{
        try {
            const response = await api.get(`/notes?offset=${offset}`)
            return response.data
        } catch (error) {
            logError(error)
            return []
        }
    }
}