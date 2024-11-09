export interface Note {
    id: string
    creator: number
    tags: Tag[]
    title: string
    source: string
    original: string
    content: unknown[]
    cover: string
    createdAt: number
    copiedAt: number
    type: number
    checked: boolean
    categoryId?: string
    icon: Icon
}

export interface Tag {
    id: number
    text: string
    color: string
    owner: number
}

export interface Icon {
    type: string
    data: string
}
