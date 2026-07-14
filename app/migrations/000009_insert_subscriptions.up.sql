INSERT INTO products (id, name, description, price, "type")
VALUES
  (1, 'Basic Subscription For 30 Days', '30', 10, 'subscription'),
  (2, 'Premium Subscription For 30 Days', '30', 20, 'subscription'),
  (3, 'Basic Subscription For 90 Days', '90', 25, 'subscription'),
  (4, 'Premium Subscription For 90 Days', '90', 50, 'subscription');

INSERT INTO subscription_levels (id, "level", "limit", "name")
VALUES
    (0, 0 ,0 , 'None'),
    (1, 1, 3, 'Basic'),
    (2, 2, 10, 'Premium');

INSERT INTO subscription_level_products (subscription_level_id, product_id)
VALUES
    (1, 1),
    (2, 2),
    (1, 3),
    (2, 4);