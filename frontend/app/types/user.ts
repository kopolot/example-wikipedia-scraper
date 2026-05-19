export type UserType = {
    id?: number;
    email: string;
    username: string;
    role: UserRoles;
    lastLoginAt?: Date;
    emailVerified?: boolean;
    createdAt?: Date;
    updatedAt?: Date;
    subscriptionLevel?: number;
    subscriptionExpiration?: Date;
} | null;

export type UserRoles = 'user' | 'admin'
