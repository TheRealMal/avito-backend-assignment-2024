-- Добавляем 1000 записей в таблицу banners
INSERT INTO banners (is_active, feature, content, updated_at)
SELECT
   (random() > 0.5)::boolean AS is_active,
   floor(random()*((2147483647)::bigint-(1)::bigint+1))::bigint AS feature,
   convert_to('{"test":true}', 'LATIN1') AS content,
   CURRENT_DATE + (i || ' days')::interval AS updated_at
FROM generate_series(1, 1000) AS s(i);

-- Добавляем по 1-3 тега для каждой записи из banners
WITH banner_ids AS (
  SELECT id FROM banners
),
random_tags AS (
  SELECT 
    floor(random()*((2147483647)::bigint-(1)::bigint+1))::bigint AS tag_id,
    id as banner_id,
    generate_series(1, (random() * 3)::int)
  FROM 
    banner_ids
)
INSERT INTO tags (id, banner_id)
SELECT tag_id, banner_id
FROM random_tags;