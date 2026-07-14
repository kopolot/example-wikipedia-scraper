import type { SubscriptionLevel } from "@/types/subscription_level";

export const useUser = () => {
    const api = useApi();
    const subscriptionData = useState<SubscriptionLevel | null>('user.subscriptionData', () => null);

    async function changeEmail(newEmail: string, password: string): Promise<{ success: boolean; error: string | null }> {
        try {
            const response = await api.post({
                url: '/user/change_email',
                body: { new_email: newEmail, password },
            });
            if (response.body?.success) {
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

    async function loadSubscriptionData(): Promise<SubscriptionLevel | null> {
        try {
            const response = await api.get({
                url: '/user/subscription_levels/' + useAuth().getUser()?.subscriptionLevel,
            });
            if (response.body?.success) {
                subscriptionData.value = response.body.data as SubscriptionLevel;
                return subscriptionData.value;
            } else {
                console.error('Error fetching subscription data:', response.body?.errors);
                return null;
            }
        } catch (e) {
            console.error('Error fetching subscription data:', e);
            return null;
        }
    }

    function getSubscriptionData(): SubscriptionLevel | null {
        if (!subscriptionData.value) {
            loadSubscriptionData();
        }
        return subscriptionData.value;
    }

    return {
        changeEmail,
        changePassword,
        logoutAll,
        loadSubscriptionData,
        getSubscriptionData,
        subscriptionData,
    };
}
