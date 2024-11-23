export interface Note {
    id: string
    creator: number
    tags: Tag[]
    title: string
    source: string
    original: string
    content: ContentElement[]
    cover: string
    createdAt: number
    copiedAt: number
    type: number
    checked: boolean
    categoryId?: string
    icon: Icon
}

export class NoteTemplate implements Note {
    id: string = ''
    creator: number = 0
    tags: Tag[] = []
    title: string = ''
    source: string = ''
    original: string = ''
    content: ContentElement[] = []
    cover: string = ''
    createdAt: number = 0
    copiedAt: number = 0
    type: number = 0
    checked: boolean = false
    categoryId: string = ''
    icon: Icon = { type: '', data: '' }

    constructor() {
        
    }
}

export interface ContentElement {
    type: string
}

export interface Tag {
    text: string
    color: string
}

export interface Icon {
    type: string
    data: string
}

export interface RichText{
    plain_text: string
    meta: TextMeta[] 
}

export interface Text{
    type: string
    rich_text: RichText
}

export interface Checkbox{
    type: string
    checked: boolean
    rich_text: RichText
}

export interface Image{
    type: string
    src: string
    caption: string
    width: number
    align: string
}
export interface TextMeta{
    offset: number
    length: number
    link: string
    color: string
    weight: string
    italic: boolean
    underline: boolean
    cross_out: boolean
}