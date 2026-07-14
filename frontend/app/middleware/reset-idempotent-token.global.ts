import { v4 as uuid } from 'uuid';

export default defineNuxtRouteMiddleware((to, from) => {
    const getIdempotentToken = () => useState('idempotentToken', () => uuid());
    getIdempotentToken().value = uuid();
});