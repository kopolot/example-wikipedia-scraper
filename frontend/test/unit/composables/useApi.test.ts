
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useApi } from '../../../app/composables/useApi';

const testToken = 'test-token-123';

// Mock Nuxt composables używając vi.hoisted
const mockFetch = vi.hoisted(() => {
    const mockFetch = vi.fn();

    // Mock globalnych funkcji
    vi.stubGlobal('$fetch', mockFetch);
    vi.stubGlobal('useRuntimeConfig', vi.fn(() => ({
        public: {
            apiBase: 'http://test-api.com'
        }
    })));
    vi.stubGlobal('useToken', vi.fn(() => ({
        getToken: vi.fn(() => testToken)
    })));

    return mockFetch;
});

describe('useApi', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should make a GET request', async () => {
        const mockResponse = { success: true, data: { message: 'OK' } };
        mockFetch.mockResolvedValue(mockResponse);

        const api = useApi();
        const response = await api.get({ url: '/health', method: 'GET' });

        expect(response).toBeDefined();
        expect(response.statusCode).toBe(200);
        expect(response.body).toEqual(mockResponse);
        expect(mockFetch).toHaveBeenCalledWith('http://test-api.com/health', {
            method: 'GET',
            headers: expect.any(Object),
            onResponseError: expect.any(Function)
        });
    });

    it('login failure POST request', async () => {
        const mockErrorResponse = {
            success: false,
            errors: [{ message: 'Invalid credentials' }]
        };

        // Mock error response
        const mockError = new Error('HTTP Error') as Error & {
            response: { status: number };
            data: unknown;
        };
        mockError.response = { status: 401 };
        mockError.data = mockErrorResponse;

        mockFetch.mockRejectedValue(mockError);

        const api = useApi();
        const response = await api.post({
            url: '/user/login',
            method: 'POST',
            body: JSON.stringify({ login: 'wrong', password: 'credentials' }),
        });

        expect(response.statusCode).toBe(401);
        expect(response.body).toEqual(mockErrorResponse);
        expect(mockFetch).toHaveBeenCalledWith('http://test-api.com/user/login', {
            method: 'POST',
            headers: expect.any(Object),
            body: { login: 'wrong', password: 'credentials' },
            onResponseError: expect.any(Function)
        });
    });

    it('should have token in headers when provided', async () => {
        const mockResponse = { success: true, data: { message: 'OK' } };
        mockFetch.mockResolvedValue(mockResponse);

        const api = useApi();

        const response = await api.get({
            url: '/protected/resource',
            method: 'GET',
        });

        expect(response).toBeDefined();
        expect(response.statusCode).toBe(200);
        expect(response.body).toEqual(mockResponse);
        expect(mockFetch).toHaveBeenCalledWith('http://test-api.com/protected/resource', {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${testToken}` },
            onResponseError: expect.any(Function)
        });
    });

    it('test handling of query parameters in GET request', async () => {
        const mockResponse = { success: true, data: { items: [1, 2, 3] } };
        mockFetch.mockResolvedValue(mockResponse);

        const api = useApi();

        const response = await api.get({
            url: '/items',
            method: 'GET',
            query_params: {
                category: ['books', 'electronics'],
                sort: ['asc']
            }
        });
        expect(response).toBeDefined();
        expect(response.statusCode).toBe(200);
        expect(response.body).toEqual(mockResponse);
        expect(mockFetch).toHaveBeenCalledWith(
            'http://test-api.com/items?category=books&category=electronics&sort=asc',
            {
                method: 'GET',
                headers: { 'Authorization': `Bearer ${testToken}` },
                onResponseError: expect.any(Function)
            }
        );
    }, 10000);
});
