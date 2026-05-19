import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useAuth } from '../../../app/composables/useAuth';
import type { UserType } from '../../../app/types/user';

// Mock Nuxt composables
const { mockApiPost, mockCookie, mockUser } = vi.hoisted(() => {
    // Mock useCookie
    const mockCookie = { value: null as string | null };
    vi.stubGlobal('useCookie', vi.fn(() => mockCookie));

    // Mock useState
    const mockUser = { value: null as UserType };
    vi.stubGlobal('useState', vi.fn(() => mockUser));

    // Mock computed
    vi.stubGlobal('computed', vi.fn((fn) => ({ value: fn() })));

    // Mock readonly
    vi.stubGlobal('readonly', vi.fn((ref) => ref));

    // Mock navigateTo
    vi.stubGlobal('navigateTo', vi.fn());

    // Mock import.meta.client
    Object.defineProperty(import.meta, 'client', {
        value: true,
        writable: true
    });

    // Mock useApi - TYLKO to co zwraca useApi!
    const mockApiPost = vi.fn();
    vi.stubGlobal('useApi', vi.fn(() => ({
        get: vi.fn(),
        post: mockApiPost,
        put: vi.fn(),
        delete: vi.fn(),
        patch: vi.fn()
    })));

    return { mockApiPost, mockCookie, mockUser };
});

describe('useAuth', () => {
    const testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJyb2xlIjoidXNlciIsInN1YiI6MSwiZXhwIjo5OTk5OTk5OTk5LCJpYXQiOjE2MDAwMDAwMDB9.test-signature";

    beforeEach(() => {
        vi.clearAllMocks();
        // Reset wszystkich mock refs używając zmiennych z hoisted
        mockCookie.value = null;
        mockUser.value = null;
    });

    it('should login successfully with valid credentials', async () => {
        const mockResponse = {
            statusCode: 200,
            body: {
                success: true,
                data: { token: testToken }
            },
            headers: {}
        };

        mockApiPost.mockResolvedValue(mockResponse);

        const auth = useAuth();
        const { success: result } = await auth.login('test@test.com', 'password');

        expect(result).toBe(true);
        expect(mockApiPost).toHaveBeenCalledWith({
            url: '/user/login',
            body: JSON.stringify({ login: 'test@test.com', password: 'password' }),
        });
        expect(auth.getToken()).toBe(testToken);
    });

    it('should return false on login failure', async () => {
        const mockResponse = {
            statusCode: 401,
            body: {
                success: false,
                errors: [{ message: 'Invalid credentials' }]
            },
            headers: {}
        };

        mockApiPost.mockResolvedValue(mockResponse);

        const auth = useAuth();
        const { success: result } = await auth.login('test@test.com', 'wrong-password');

        expect(result).toBe(false);
    });

    it('should decode token correctly', () => {
        const auth = useAuth();
        const decoded = auth.decodeToken(testToken);

        expect(decoded).toEqual({
            email: 'test@test.com',
            role: 'user',
            user_id: 1,
            issued_at: 1600000000,
            expires_at: 9999999999
        });
    });

    it('should return null for invalid token', () => {
        const auth = useAuth();
        const decoded = auth.decodeToken('invalid.token');

        expect(decoded).toBeNull();
    });

    it('should check if user has specific role', () => {
        const auth = useAuth();

        // Mock user z admin role używając zmiennej z hoisted
        mockUser.value = {
            email: 'admin@test.com',
            role: 'admin',
            id: 1,
            username: 'admin'
        };

        expect(auth.hasRole('admin')).toBe(true);
        expect(auth.hasRole('user')).toBe(false);
    });

    it('should check authentication status correctly', () => {
        const auth = useAuth();

        // Mock valid token w cookie używając zmiennej z hoisted
        mockCookie.value = testToken;

        expect(auth.isAuthenticated()).toBe(true);

        // Test expired token
        const expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJyb2xlIjoidXNlciIsInN1YiI6MSwiZXhwIjoxNjAwMDAwMDAwLCJpYXQiOjE2MDAwMDAwMDB9.test-signature";
        mockCookie.value = expiredToken;

        expect(auth.isAuthenticated()).toBe(false);
    });

    it('should register successfully with valid data', async () => {
        const mockResponse = {
            statusCode: 201,
            body: {
                success: true
            },
            headers: {}
        };

        mockApiPost.mockResolvedValue(mockResponse);

        const auth = useAuth();
        const { success: result } = await auth.register('newuser', 'newuser@test.com', 'password', 'password');

        expect(result).toBe(true);
        expect(mockApiPost).toHaveBeenCalledWith({
            url: '/user/',
            body: { username: 'newuser', email: 'newuser@test.com', password: 'password', repeat_password: 'password' },
        });
    });
});