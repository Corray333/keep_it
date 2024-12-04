import type { Note, Tag } from "@/entities/note"
import { NoteTransport } from "@/transport/note"

export class NoteService {
    noteTransport = new NoteTransport()
    getNotes = async (offset: number, tags: string[], category: string) : Promise<Note[]>=>{
        return await this.noteTransport.getNotes(offset, tags, category)
    }

    deleteNotes = async (noteIDs: string[]) : Promise<boolean>=>{
        return await this.noteTransport.deleteNotes(noteIDs)
    }

    createTag = async (text: string, color: string) : Promise<boolean> =>{
        return await this.noteTransport.createTag(text, color)
    }

    deleteTag = async (text: string) : Promise<boolean> =>{
        return await this.noteTransport.deleteTag(text)
    }

    getTags = async () : Promise<Tag[]> =>{
        return await this.noteTransport.getTags()
    }

    removeTagFromNote = async (noteId: string, tagText: string) : Promise<boolean> =>{
        return await this.noteTransport.removeTagFromNote(noteId, tagText)
    }

    addTagToNote = async (noteId: string, tagText: string, color: string, isNew: boolean) : Promise<boolean> =>{
        return await this.noteTransport.addTagToNote(noteId, tagText, color, isNew)
    }

    setNoteCategory = async (noteId: string, categoryID: string) : Promise<boolean> =>{
        return await this.noteTransport.setNoteCategory(noteId, categoryID)
    }
}