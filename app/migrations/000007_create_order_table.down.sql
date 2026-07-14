DROP TABLE order_items;
DROP INDEX IF EXISTS idx_order_item_order;
DROP TYPE order_item_type;

DROP TABLE "orders";

DROP INDEX IF EXISTS idx_order_user;
DROP INDEX IF EXISTS idx_order_status;
DROP INDEX IF EXISTS idx_order_payment_id;

DROP TYPE order_status;