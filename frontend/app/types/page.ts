export type PageType = {
    id: number
    createdAt: Date
    updatedAt: Date
    deletedAt?: Date
    siteName?: string
    title: string
    url: string
    content: string
    textField1?: string
    textField2?: string
    textField3?: string
    externalId?: string
}