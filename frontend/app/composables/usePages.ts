import type { PageType } from "~/types/page";

export function usePages() {
    const api = useApi();

    async function fetchPageRecords(
        limit = 20,
        page = 1,
        sort = 'created_at',
        order = 'desc',
    ): Promise<{ pages: PageType[]; maxPage: number; error: string | null }> {
        try {
            const params: Record<string, string[]> = {
                page: [page.toString()],
                limit: [limit.toString()],
                sort: [sort],
                order: [order],
            };
            const response = await api.get({
                url: '/page-records/list',
                query_params: params,
            })
            if (response.body?.success) {
                const typedResponseData = response.body.data as { pageRecords: PageType[]; maxPage: number };
                const pageRecords = (typedResponseData.pageRecords as PageType[]).map(page => ({
                    ...page,
                    createdAt: new Date(page.createdAt),
                    updatedAt: new Date(page.updatedAt),
                    deletedAt: page.deletedAt ? new Date(page.deletedAt) : undefined,
                }));
                return { pages: pageRecords, maxPage: typedResponseData?.maxPage ?? 0, error: null };
            }
            return { pages: [], maxPage: 0, error: 'page_records.fetch_error' };
        } catch (e) {
            console.error('Error fetching pages:', e);
            return { pages: [], maxPage: 0, error: 'request_failed' };
        }
    }


    async function fetchPageById(id: string | number): Promise<{ pageRecord: PageType | null; error: string | null }> {
        try {
            const response = await api.get({ url: `/page-records/${id}` });
            if (response.body?.success && response.body.data) {
                return { pageRecord: response.body.data as PageType, error: null };
            } else {
                return { pageRecord: null, error: 'page_records.fetch_error' };
            }
        } catch (e) {
            console.error('Error fetching page record:', e);
            return { pageRecord: null, error: 'request_failed' };
        }
    }

    return {
        fetchPageRecords,
        fetchPageById,
    }
}