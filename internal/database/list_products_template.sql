-- ListProductsBase template: the service will replace the placeholders.
SELECT
  p.product_id,
  p.name,
  p.slug,
  p.brand,
  pv.price,
  pv.primary_image_url,
  p.in_stock
FROM product p
JOIN product_variant pv
  ON pv.variant_id = p.default_variant_id

/*{{DYNAMIC_JOINS}}*/

WHERE p.store_id = $1
  AND p.deleted_at IS NULL

/*{{DYNAMIC_WHERE}}*/

ORDER BY p.product_id DESC
LIMIT $2 OFFSET $3;
