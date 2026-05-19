export const useUser = () => {
    const api = useApi();

    async function changeEmail(newEmail: string, password: string): Promise<{ success: boolean; error: string | null }> {
        try {
            const response = await api.post({
                url: '/user/change_email',
                body: { new_email: newEmail, password },
            });
            if (response.body?.success) {
                // useAuth().logout(); // Wyloguj użytkownika po zmianie emaila
                return { success: true, error: null };
            } else {
                console.error('Error changing email:', response.body?.errors);
                return { success: false, error: 'user.change_email_error' };
            }
        } catch (e) {
            console.error('Error changing email:', e);
            return { success: false, error: 'request_failed' };
        }
    }

    async function changePassword(oldPassword: string, newPassword: string): Promise<{ success: boolean; error: string | null }> {
        try {
            const response = await api.post({
                url: '/user/change_password',
                body: { old_password: oldPassword, new_password: newPassword },
            });
            if (response.body?.success) {
                return { success: true, error: null };
            } else {
                console.error('Error changing password:', response.body?.errors);
                return { success: false, error: 'user.change_password_error' };
            }
        } catch (e) {
            console.error('Error changing password:', e);
            return { success: false, error: 'request_failed' };
        }
    }

    async function logoutAll(): Promise<{ success: boolean; error: string | null }> {
        try {
            const response = await api.post({
                url: '/user/logout_all',
            });
            if (response.body?.success) {
                return { success: true, error: null };
            } else {
                console.error('Error in logout all:', response.body?.errors);
                return { success: false, error: 'user.logout_all_error' };
            }
        } catch (e) {
            console.error('Error in logout all:', e);
            return { success: false, error: 'request_failed' };
        }
    }

    return {
        changeEmail,
        changePassword,
        logoutAll,
    };
}