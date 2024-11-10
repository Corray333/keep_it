import type { Note, Tag } from "@/entities/note"
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

    createTag = async (text: string, color: string) : Promise<boolean> =>{
        try {
            await api.post('/tags', {
                text: text,
                color: color
            })
            return true
        } catch (error) {
            logError(error)
            return false
        }
    }

    deleteTag = async (text: string) : Promise<boolean> =>{
        try {
            await api.delete(`/tags/${text}`)
            return true
        } catch (error) {
            logError(error)
            return false
        }
    }

    getTags = async () : Promise<Tag[]> =>{
        try {
            const response = await api.get('/tags')
            return response.data
        } catch (error) {
            logError(error)
            return []
        }
    }

    removeTagFromNote = async (noteId: string, tagText: string) : Promise<boolean> =>{
        try {
            await api.delete(`/notes/${noteId}/tags/${tagText}`)
            return true
        } catch (error) {
            logError(error)
            return false
        }
    }

    addTagToNote = async (noteId: string, tagText: string, color: string, isNew: boolean) : Promise<boolean> =>{
        try {
            await api.post(`/notes/${noteId}/tags`,{
                text: tagText,
                color: color,
                isNew: isNew
            })
            return true
        } catch (error) {
            logError(error)
            return false
        }
    }
}