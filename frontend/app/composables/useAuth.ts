import type { UserRoles, UserType } from '@/types/user';
import type { JWTClaims } from '~/types/auth';

export const useToken = () => {
    const tokenCookie = useCookie<string | null>('auth-token', {
        default: () => null,
        httpOnly: false,
        secure: typeof window !== 'undefined' ? window.location.protocol === 'https:' : false,
        sameSite: 'strict',
        maxAge: 60 * 60 * 24 * 7 // 7 dni
    });

    function getToken(): string | null {
        return tokenCookie.value;
    }

    function setToken(newToken: string) {
        tokenCookie.value = newToken;
    }

    function clearToken() {
        tokenCookie.value = null;
    }

    return {
        getToken,
        setToken,
        clearToken
    };
};

export const useAuth = () => {
    const api = useApi();
    const user = useState<UserType>('auth.user', () => null);
    const token = useToken();

    async function login(login: string, password: string): Promise<{ success: boolean, errors?: string[] }> {
        try {
            const response = await api.post({
                url: '/user/login',
                body: { login, password },
            });
            if (response.body?.success && response.body.data && 'token' in response.body.data) {
                setToken(response.body.data.token as string);
                return { success: true };
            } else /*if (responseBody && 'error' in responseBody) */ {
                console.error('Login failed:', response.body?.errors);
                return { success: false, errors: response.body?.errors?.map(err => err.message) };
            }
        } catch (e: unknown) {
            console.error('Login error:', e);
            return { success: false, errors: ['login.error'] };
        }
    }

    async function register(username: string, email: string, password: string, repeatPassword: string): Promise<{ success: boolean; errors?: string[] }> {
        try {
            const response = await api.post({
                url: '/user/',
                body: { username, email, password, repeat_password: repeatPassword },
            });
            if (response.body?.success) {
                return { success: true };
            } else {
                console.error('Registration failed:', response.body?.errors);
                return { success: false, errors: response.body?.errors?.map(err => err.message) };
            }
        } catch (e) {
            console.error('Registration error:', e);
            return { success: false, errors: ['registration.error'] };
        }
    }

    async function forgotPassword(email: string): Promise<{ success: boolean; errors?: string[] }> {
        try {
            const response = await api.post({
                url: '/user/forgot_password',
                body: { email },
            });
            if (response.body?.success) {
                return { success: true };
            } else {
                console.error('Error in forgot password:', response.body?.errors);
                return { success: false, errors: response.body?.errors?.map(err => err.message) };
            }
        } catch (e) {
            console.error('Error in forgot password:', e);
            return { success: false, errors: ['request_failed'] };
        }
    }

    async function resetPassword(token: string, newPassword: string): Promise<{ success: boolean; errors?: string[] }> {
        try {
            const response = await api.post({
                url: '/user/reset_password',
                body: { token, new_password: newPassword },
            });
            if (response.body?.success) {
                return { success: true };
            } else {
                console.error('Error in reset password:', response.body?.errors);
                return { success: false, errors: response.body?.errors?.map(err => err.message) };
            }
        } catch (e) {
            console.error('Error in reset password:', e);
            return { success: false, errors: ['request_failed'] };
        }
    }

    async function verifyEmail(token: string): Promise<{ success: boolean; errors?: string[] }> {
        try {
            const response = await api.post({
                url: '/user/verify_email',
                body: { token },
            });
            if (response.body?.success) {
                return { success: true };
            } else {
                console.error('Error in verify email:', response.body?.errors);
                return { success: false, errors: response.body?.errors?.map(err => err.message) };
            }
        } catch (e) {
            console.error('Error in verify email:', e);
            return { success: false, errors: ['request_failed'] };
        }
    }

    async function logout() {
        token.clearToken();
        user.value = null;
        // Przekieruj na stronę logowania
        await navigateTo('/');
    }

    function decodeToken(token: string): JWTClaims | null {
        try {
            const parts = token.split('.');
            if (parts.length !== 3 || !parts[1]) {
                throw new Error('Invalid token format');
            }
            const payload = parts[1];
            const decodedPayload = atob(payload);
            const userData = JSON.parse(decodedPayload);
            return {
                email: userData.email,
                role: userData.role,
                user_id: Number(userData.sub),
                issued_at: userData.iat,
                expires_at: userData.exp,
                username: userData.username,
                created_at: userData.created_at,
                updated_at: userData.updated_at,
            };
        } catch (e) {
            console.error('Error decoding token:', e);
            return null;
        }
    }

    function setToken(newToken: string) {
        const decodedToken = decodeToken(newToken);
        if (decodedToken) {
            setUserFromDecodedToken(decodedToken);
            token.setToken(newToken);
        } else {
            console.error('Invalid token format');
        }
    }

    function getToken(): string | null {
        return token.getToken();
    }

    function getUser(): UserType | null {
        // Jeśli user nie istnieje ale token tak, spróbuj go zdekodować
        const tokenValue = token.getToken();
        if (!user.value && tokenValue) {
            const decoded = decodeToken(tokenValue);
            if (decoded) {
                setUserFromDecodedToken(decoded);
            }
        }
        return user.value;
    }

    function isAuthenticated(): boolean {
        const currentToken = getToken();
        if (!currentToken) return false;

        const decoded = decodeToken(currentToken);
        if (!decoded) return false;

        // Sprawdź czy token nie wygasł
        const now = Math.floor(Date.now() / 1000);
        const isValid = decoded.expires_at > now;

        if (!isValid) {
            // Token wygasł, wyczyść
            logout();
            return false;
        }

        return true;
    }

    function hasRole(role: 'user' | 'admin'): boolean {
        const currentUser = getUser();
        return currentUser?.role === role;
    }

    function setUserFromDecodedToken(decoded: JWTClaims) {
        const decodedUser: UserType = {
            email: decoded.email,
            role: decoded.role,
            id: Number(decoded.user_id),
            username: decoded.username,
            createdAt: unixTimestampToDate(decoded.created_at),
            updatedAt: unixTimestampToDate(decoded.updated_at),
        };
        user.value = decodedUser;
    }

    // czy to musi byc w auth czy powiino byc w user
    async function loadUserData() {
        try {
            if (!user.value || !user.value.id) {
                throw new Error('User ID is missing');
            }
            const response = await api.get({
                url: '/user/' + user.value.id,
            });
            if (response.body?.success && response.body.data) {
                const userData = response.body.data as {
                    email: string;
                    username: string;
                    role: UserRoles;
                    lastLoginAt: string;
                    emailVerified?: boolean;
                    createdAt: string;
                    updatedAt: string;
                    subscriptionLevel: number;
                    subscriptionExpiration: string;
                };
                if (userData) {
                    user.value = {
                        ...user.value,
                        email: userData.email,
                        username: userData.username,
                        role: userData.role,
                        lastLoginAt: new Date(userData.lastLoginAt),
                        emailVerified: userData.emailVerified ?? false,
                        createdAt: new Date(userData.createdAt),
                        updatedAt: new Date(userData.updatedAt),
                        subscriptionLevel: userData.subscriptionLevel,
                        subscriptionExpiration: userData.subscriptionExpiration ? new Date(userData.subscriptionExpiration) : undefined,
                    };
                }
            } else {
                console.error('Failed to load user data:', response.body?.errors);
            }
        } catch (e) {
            console.error('Error loading user data:', e);
        }
    }

    const initializeAuth = () => {
        const tokenValue = token.getToken();
        if (tokenValue && !user.value) {
            getUser();
        }
    };

    if (import.meta.client) {
        initializeAuth();
    }

    const isLoggedIn = computed(() => {
        return isAuthenticated() && !!user.value;
    });

    return {
        // Reactive state
        user: readonly(user),
        isLoggedIn: readonly(isLoggedIn),

        // Methods
        login,
        register,
        logout,
        forgotPassword,
        resetPassword,
        getUser,
        getToken,
        isAuthenticated,
        hasRole,
        verifyEmail,
        loadUserData,

        // Utils
        decodeToken
    };
}