import type { ApiRequest, ApiResponse, ApiInterface, ApiResponseBody, PostRequest, PutRequest, DeleteRequest, PatchRequest, GetRequest } from "@/types/api";
import type { FetchError, FetchResponse } from 'ofetch';
import type { AvailableRouterMethod, NitroFetchOptions, NitroFetchRequest } from 'nitropack/types'

export const useApi = (): ApiInterface => {
    const config = useRuntimeConfig()
    const baseURL = config.public.apiBase

    const buildUrl = (request: ApiRequest): string => {
        const baseUrlObject = new URL(baseURL);
        const requestUrl = request.url.startsWith('/') ? request.url : `/${request.url}`;
        // requestUrl = requestUrl.endsWith('/') ? requestUrl : `${requestUrl}/`;
        const url = new URL(baseUrlObject.pathname + requestUrl, baseUrlObject.origin);
        if (request.query_params) {
            Object.entries(request.query_params).forEach(([key, values]) => {
                values.forEach(value => url.searchParams.append(key, value));
            });
        }
        return url.toString();
    }

    const processResponse = (response: FetchResponse<unknown>): ApiResponse => {
        const body = response._data as ApiResponseBody;
        if (!body?.success) {
            console.warn('API Error:', body?.errors);
        }
        return {
            statusCode: response.status,
            body,
            headers: Object.fromEntries(response.headers)
        };
    }

    async function get(request: GetRequest): Promise<ApiResponse> {
        return fetchRequest({ ...request, method: 'GET' });
    }

    async function fetchRequest(
        request: ApiRequest
    ): Promise<ApiResponse> {
        try {
            const url = buildUrl(request);
            const fetchOptions = makeFetchOptions(request)
            const response: FetchResponse<unknown> = await $fetch.raw(url, fetchOptions);
            return processResponse(response);
        } catch (error) {
            const fetchError = error as FetchError;
            return {
                statusCode: fetchError.response?.status || 500,
                body: fetchError.data || {
                    success: false,
                    errors: [{ message: fetchError.message || 'request_failed' }]
                },
                headers: Object.fromEntries(fetchError.response?.headers || [])
            };
        }
    }

    function makeFetchOptions(request: ApiRequest) {
        const headers = request.headers || {} as Record<string, string>;
        const tokenValue = useToken().getToken();
        if (tokenValue)
            headers['Authorization'] = `Bearer ${tokenValue}`;
        const fetchOptions: NitroFetchOptions<NitroFetchRequest, AvailableRouterMethod<NitroFetchRequest>> = {
            method: request.method,
            headers: headers,
            onResponseError({ response }) {
                console.error(`API ${request.method} Error:`, response.status, response.statusText);
            },
            responseType: 'json',
        };
        if (request.body) {
            fetchOptions.body = typeof request.body === 'string'
                ? JSON.parse(request.body)
                : request.body;
        }
        return fetchOptions;
    }

    return {
        get,
        post: (request: PostRequest) => fetchRequest({ ...request, method: 'POST' }),
        put: (request: PutRequest) => fetchRequest({ ...request, method: 'PUT' }),
        delete: (request: DeleteRequest) => fetchRequest({ ...request, method: 'DELETE' }),
        patch: (request: PatchRequest) => fetchRequest({ ...request, method: 'PATCH' })
    };
}