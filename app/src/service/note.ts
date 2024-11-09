import type { Note } from "@/entities/note"
import { NoteTransport } from "@/transport/note"

export class NoteService {
    noteTransport = new NoteTransport()
    getNotes = async (offset: number) : Promise<Note[]>=>{
        return await this.noteTransport.getNotes(offset)
    }
}