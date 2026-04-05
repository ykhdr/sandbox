CREATE OR REPLACE VIEW order_details AS
SELECT
    o.id AS order_id,
    u.name AS user_name,
    p.name AS product_name,
    o.quantity,
    o.total_price,
    o.status,
    o.created_at
FROM orders o
JOIN users u ON o.user_id = u.id
JOIN products p ON o.product_id = p.id;

CREATE OR REPLACE VIEW monthly_sales_summary AS
SELECT
    date_trunc('month', created_at) AS month,
    COUNT(*) AS total_orders,
    SUM(total_price) AS total_revenue,
    AVG(total_price) AS avg_check
FROM orders
GROUP BY date_trunc('month', created_at)
ORDER BY month;
