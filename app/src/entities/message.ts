export interface Message {
    text: string
    type: MessageType
}

export enum MessageType{
    ERROR,
    INFO
}