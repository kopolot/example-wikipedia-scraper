import { describe, it, expect } from "vitest";
import { useApi } from '../../../app/composables/useApi';

describe('useApi integration', () => {
    it('should make a GET request integaration', async () => {
        const api = useApi();
        const response = await api.get({ url: '/health', method: 'GET' });
        expect(response).toBeDefined();
        expect(response.statusCode).toBe(200);
        expect(response.body).toBeDefined();
        expect(typeof response.body.success).toBe('boolean');
        expect(response.body.success).toBe(true);
        expect(response.body.data).toBeDefined();
        expect(response.body.data.status).toBe('ok');
    });
});
