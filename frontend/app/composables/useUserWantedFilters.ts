import type { UserWantedPageCriteria, UserWantedPagesFilter } from '~/types/user_wanted_pages_filter'
import type { PageType } from '~/types/page'

function mapFilterDates(filter: UserWantedPagesFilter): UserWantedPagesFilter {
    return {
        ...filter,
        createdAt: new Date(filter.createdAt),
        updatedAt: new Date(filter.updatedAt),
        deletedAt: filter.deletedAt ? new Date(filter.deletedAt) : undefined,
    }
}

function mapPageDates(pages: PageType[]): PageType[] {
    return pages.map(page => ({
        ...page,
        createdAt: new Date(page.createdAt),
        updatedAt: new Date(page.updatedAt),
        deletedAt: page.deletedAt ? new Date(page.deletedAt) : undefined,
    }))
}

export function useUserWantedFilters() {
    const api = useApi()

    async function fetchFilters(): Promise<{ filters: UserWantedPagesFilter[]; error: string | null }> {
        try {
            const response = await api.get({ url: '/user_wanted_filters/' })
            if (response.body?.success && response.body.data) {
                const filters = (response.body.data as UserWantedPagesFilter[]).map(mapFilterDates)
                return { filters, error: null }
            }
            return { filters: [], error: 'user_wanted_filters.fetch_error' }
        } catch (e) {
            console.error('Error fetching user wanted filters:', e)
            return { filters: [], error: 'request_failed' }
        }
    }

    async function fetchFilteredPages(): Promise<{ pages: PageType[]; error: string | null }> {
        try {
            const response = await api.get({ url: '/user_wanted_filters/filtered_pages' })
            if (response.body?.success && response.body.data) {
                return { pages: mapPageDates(response.body.data as PageType[]), error: null }
            }
            return {
                pages: [],
                error: response?.body?.errors?.map(e => e.message).join(', ') ?? 'user_wanted_filters.pages_fetch_error',
            }
        } catch (e) {
            console.error('Error fetching filtered pages:', e)
            return { pages: [], error: 'request_failed' }
        }
    }

    async function createFilter(criteria: UserWantedPageCriteria): Promise<{ filter: UserWantedPagesFilter | null; error: string | null }> {
        try {
            const response = await api.post({ url: '/user_wanted_filters/', body: criteria })
            if (response.body?.success && response.body.data) {
                return { filter: mapFilterDates(response.body.data as UserWantedPagesFilter), error: null }
            }
            return {
                filter: null,
                error: response?.body?.errors?.map(e => e.message).join(', ') ?? 'user_wanted_filters.create_error',
            }
        } catch (e) {
            console.error('Error creating user wanted filter:', e)
            return { filter: null, error: 'request_failed' }
        }
    }

    async function fetchFilter(id: number): Promise<{ filter: UserWantedPagesFilter | null; error: string | null }> {
        try {
            const response = await api.get({ url: `/user_wanted_filters/${id}` })
            if (response.body?.success && response.body.data) {
                return { filter: mapFilterDates(response.body.data as UserWantedPagesFilter), error: null }
            }
            return { filter: null, error: 'user_wanted_filters.fetch_error' }
        } catch (e) {
            console.error('Error fetching user wanted filter:', e)
            return { filter: null, error: 'request_failed' }
        }
    }

    async function fetchPagesByFilters(filterIds: number[]): Promise<{ pages: PageType[]; error: string | null }> {
        try {
            const response = await api.get({
                url: '/user_wanted_filters/pages_by_filters',
                query_params: { id: filterIds.map(String) },
            })
            if (response.body?.success) {
                return { pages: mapPageDates((response.body.data ?? []) as PageType[]), error: null }
            }
            return {
                pages: [],
                error: response?.body?.errors?.map(e => e.message).join(', ') ?? 'user_wanted_filters.pages_fetch_error',
            }
        } catch (e) {
            console.error('Error fetching pages by filter:', e)
            return { pages: [], error: 'request_failed' }
        }
    }

    async function updateFilter(id: number, criteria: UserWantedPageCriteria): Promise<{ filter: UserWantedPagesFilter | null; error: string | null }> {
        try {
            const response = await api.put({ url: `/user_wanted_filters/${id}`, body: criteria })
            if (response.body?.success && response.body.data) {
                return { filter: mapFilterDates(response.body.data as UserWantedPagesFilter), error: null }
            }
            return {
                filter: null,
                error: response?.body?.errors?.map(e => e.message).join(', ') ?? 'user_wanted_filters.update_error',
            }
        } catch (e) {
            console.error('Error updating user wanted filter:', e)
            return { filter: null, error: 'request_failed' }
        }
    }

    async function deleteFilter(id: number): Promise<{ success: boolean; error: string | null }> {
        try {
            const response = await api.delete({ url: `/user_wanted_filters/${id}` })
            if (response.body?.success) {
                return { success: true, error: null }
            }
            return { success: false, error: 'user_wanted_filters.delete_error' }
        } catch (e) {
            console.error('Error deleting user wanted filter:', e)
            return { success: false, error: 'request_failed' }
        }
    }

    return {
        fetchFilters,
        fetchFilteredPages,
        createFilter,
        fetchFilter,
        fetchPagesByFilters,
        updateFilter,
        deleteFilter,
    }
}
