export type JWTClaims = {
    email: string;
    role: 'user' | 'admin';
    username: string;
    user_id: number;
    issued_at: number;
    expires_at: number;
    created_at: number;
    updated_at: number;
};
