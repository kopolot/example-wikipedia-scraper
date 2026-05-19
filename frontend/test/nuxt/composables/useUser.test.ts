import { describe, it, beforeAll, expect } from 'vitest';
import { useUser } from '../../../app/composables/useUser';
import { useAuth } from '../../../app/composables/useAuth';

const login = 'kopolot@gmail.com';
const password = '------';

beforeAll(async () => {
    const auth = useAuth();
    const result = await auth.login(login, password);
    if (!result.success) {
        throw new Error('Nie udało się zalogować testowego użytkownika');
    }
}, 1000000);

describe('useUser integration', () => {
    it('should change email via API', async () => {
        const userModule = useUser();
        const newEmail = 'kopolot@gmail.com';

        const result = await userModule.changeEmail(newEmail, password);
        expect(result).toBeDefined();
        expect(result.success).toBe(true);
        expect(result.error).toBeNull();
        userModule.changeEmail(login, password); // Przywrócenie oryginalnego emaila
    });
});