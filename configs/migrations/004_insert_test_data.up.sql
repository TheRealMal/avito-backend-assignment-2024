INSERT INTO banners (is_active, feature, content) VALUES (true, 1, convert_to('{"title": "some_title", "text": "some_text", "url": "some_url"}', 'LATIN1'));

INSERT INTO tags (id, banner_id) VALUES (1, 1), (2, 1);