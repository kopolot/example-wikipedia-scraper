import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useUser } from '../../../app/composables/useUser';

const { mockApiPost } = vi.hoisted(() => {
    // Mock useApi - TYLKO to co zwraca useApi!
    const mockApiPost = vi.fn();
    vi.stubGlobal('useApi', vi.fn(() => ({
        get: vi.fn(),
        post: mockApiPost,
        put: vi.fn(),
        delete: vi.fn(),
        patch: vi.fn()
    })));
    return { mockApiPost };
});

describe('useUser', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('should change email successfully', async () => {
        const userModule = useUser();

        mockApiPost.mockResolvedValueOnce({
            body: { success: true }
        });

        const result = await userModule.changeEmail('newemail@example.com', 'password123');
        expect(result).toEqual({ success: true, error: null });
        expect(mockApiPost).toHaveBeenCalledWith({
            url: '/user/change_email',
            body: { new_email: 'newemail@example.com', password: 'password123' },
        });
    });

    it('should handle error when changing email fails', async () => {
        const userModule = useUser();

        mockApiPost.mockResolvedValueOnce({
            body: { success: false, errors: ['Some error'] }
        });

        const result = await userModule.changeEmail('invalidemail', 'wrongpassword');

        expect(result).toEqual({ success: false, error: 'user.change_email_error' });
        expect(mockApiPost).toHaveBeenCalledWith({
            url: '/user/change_email',
            body: { new_email: 'invalidemail', password: 'wrongpassword' },
        });
    });

    it('should handle exception when API call fails', async () => {
        const userModule = useUser();

        mockApiPost.mockRejectedValueOnce(new Error('Network error'));

        const result = await userModule.changeEmail('newemail@example.com', 'password123');

        expect(result).toEqual({ success: false, error: 'request_failed' });
        expect(mockApiPost).toHaveBeenCalledWith({
            url: '/user/change_email',
            body: { new_email: 'newemail@example.com', password: 'password123' },
        });
    });
});