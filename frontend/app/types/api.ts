export interface ApiInterface {
    post(request: PostRequest): Promise<ApiResponse>;
    get(request: GetRequest): Promise<ApiResponse>;
    put(request: PutRequest): Promise<ApiResponse>;
    delete(request: DeleteRequest): Promise<ApiResponse>;
    patch(request: PatchRequest): Promise<ApiResponse>;
}

export type ErrorDetail = {
    message: string
    field?: string
    tag?: string
    param?: string
}

export type ApiResponseBody = {
    success: boolean
    data?: object
    errors?: ErrorDetail[]
}

export type ApiResponse = {
    statusCode: number
    body?: ApiResponseBody
    headers?: Record<string, string>
}

export type ApiRequest = {
    url: string
    method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
    headers?: Record<string, string>
    body?: string | object
    query_params?: Record<string, string[]>
    // to mozna automagionczie zrobic ze daje tylko tu np i sie doda do url
    path_params?: Record<string, string>
}

export type GetRequest = Omit<ApiRequest, 'method' | 'body'>;
export type PostRequest = Omit<ApiRequest, 'method'>;
export type PutRequest = Omit<ApiRequest, 'method'>;
export type DeleteRequest = Omit<ApiRequest, 'method' | 'body'>
export type PatchRequest = Omit<ApiRequest, 'method'>;

// osobne pole na token dla wygody uzycia?

