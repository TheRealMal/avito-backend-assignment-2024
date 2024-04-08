SELECT * FROM banners
SELECT * FROM tags

INSERT INTO banners (is_active, feature, content)
VALUES (true, 1, convert_to('{"title": "some_title", "text": "some_text", "url": "some_url"}', 'LATIN1'))

INSERT INTO banners (is_active, feature, content)
VALUES (true, 1, '{"title": "some_title", "text": "some_text", "url": "some_url"}')


UPDATE banners
SET feature = 3
WHERE id = 1;


INSERT INTO tags (banner_id)
VALUES (1), (1)


SELECT b.id, b.feature, b.content, b.created_at, b.updated_at, b.is_active, array_agg(t.id) FROM banners b JOIN tags t ON t.banner_id = b.id GROUP BY b.id LIMIT 100 OFFSET 0;