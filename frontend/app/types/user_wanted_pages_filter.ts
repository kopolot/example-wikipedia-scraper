export type UserWantedPageCriteria = {
    name: string
    siteNames?: string[]
    keywords?: string[]
    titleContains?: string
}

export type UserWantedPagesFilter = {
    id: number
    userId: number
    createdAt: Date
    updatedAt: Date
    deletedAt?: Date
    filterData: UserWantedPageCriteria
}
