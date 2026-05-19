import { describe, it, expect } from 'vitest';
import { useAuth } from '../../../app/composables/useAuth';

// UWAGA: Ten test zakłada, że backend API jest uruchomiony i dostępny pod odpowiednim adresem!
// Test integracyjny sprawdza prawdziwą komunikację z API (nie mockuje useApi)

describe('useAuth integration', () => {
  it('should login and get user data from token', async () => {
    const auth = useAuth();
    const login = 'kopolot@gmail.com';
    const password = '------';

    const result = await auth.login(login, password);
    expect(result).toBeDefined();
    expect(typeof result).toBe('object');
    expect(result.success).toBe(true);
    const token = auth.getToken();
    expect(token).toBeDefined();
    expect(typeof token).toBe('string');
    const user = auth.getUser();
    expect(user).toBeDefined();
    expect(user?.email).toBe(login);
    expect(['user', 'admin']).toContain(user?.role);
    expect(typeof user?.id).toBe('number');
    expect(auth.isAuthenticated()).toBe(true);
  });

  it('should fail login with wrong credentials', async () => {
    const auth = useAuth();
    const result = await auth.login('wrong@test.com', 'wrongpassword');
    expect(result).toBeDefined();
    expect(result.success).toBe(false);
  });

  it('should verify email with valid token', async () => {
    // tu przerobic zeby było dobrze
    const auth = useAuth();
    const validToken = '28bf1867-a42d-45fe-ab2f-03668c1ed17a';
    const result = await auth.verifyEmail(validToken);
    expect(result).toBeDefined();
    expect(typeof result).toBe('object');
    expect(result.success).toBe(true);
  });

  it('should register a new user', async () => {
    const auth = useAuth();
    const username = 'testuser';
    const email = `testuser${Date.now()}@example.com`;
    const password = 'TestPassword123!';

    const result = await auth.register(username, email, password, password);
    expect(result).toBeDefined();
    expect(typeof result).toBe('object');
    expect(result.success).toBe(true);
  });

  it('should load user correctly from token', async () => {
    const auth = useAuth();
    const login = 'kopolot@gmail.com';
    const password = '------';

    const result = await auth.login(login, password);
    expect(result).toBeDefined();
    expect(typeof result).toBe('object');
    expect(result.success).toBe(true);
    await auth.loadUserData();
    const user = auth.getUser();
    expect(user).toBeDefined();
    expect(user?.email).toBe(login);
    expect(['user', 'admin']).toContain(user?.role);
    expect(typeof user?.id).toBe('number');
    expect(user.subscriptionLevel).toBeDefined();
  });
});