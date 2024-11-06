export interface MessageI {
    text: string
    type: MessageType
}

export class Message implements MessageI {
    text: string
    type: MessageType

    constructor(text: string, type: MessageType){
        this.text = text
        this.type = type
    }
}

export enum MessageType{
    ERROR,
    INFO,
    SUCCESS,
    WARNING
}