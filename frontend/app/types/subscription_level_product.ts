import type { Product } from "./product";
import type { SubscriptionLevel } from "./subscription_level";

export type SubscriptionLevelProduct = {
    subscriptionLevelId: number;
    productId: number;
    subscriptionLevel: SubscriptionLevel;
    product: Product;
}