import type { Note, Tag } from "@/entities/note"
import { api, logError } from "./main"


export class NoteTransport {
    getNotes = async (offset: number, tags: string[], category: string) : Promise<Note[]>=>{
        let tagsFilter = ""
        for (let i = 0; i < tags.length; i++) {
            tagsFilter += `&tag=${tags[i]}`
        }
        if (category) tagsFilter += `&category=${category}`
        try {
            const response = await api.get(`/notes?offset=${offset}${tagsFilter}`)
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

    setNoteCategory = async (noteId: string, categoryiD: string) : Promise<boolean> =>{
        try {
            await api.put(`/notes/${noteId}/categories/${categoryiD}`)
            return true
        } catch (error) {
            logError(error)
            return false
        }
    }

    deleteNotes = async (noteIDs: string[]) : Promise<boolean> =>{
        let query = '/notes?'
        for (const noteID of noteIDs){
            query += `note_id=${noteID}&`
        }
        try {
            await api.delete(query)
            return true
        } catch (error) {
            logError(error)
            return false
        }
    }
}