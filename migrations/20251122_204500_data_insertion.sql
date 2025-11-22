INSERT INTO cities (name, region, location)
VALUES ('Kyiv', 'Kyiv', ST_SetSRID(ST_Point(30.5238, 50.4547), 4326)),
       ('Kamianets-Podilskyi', 'Khmelnytskyi', ST_SetSRID(ST_Point(26.58516, 48.67882), 4326)),
       ('Khmelnytskyi', 'Khmelnytskyi', ST_SetSRID(ST_Point(26.9965, 49.4216), 4326));


INSERT INTO reports (user_id, location, city_id, description, category_id, current_status_id)
VALUES
(1, ST_SetSRID(ST_Point(30.5245, 50.4552), 4326), '1f9588e6-a558-454b-92b5-48f5cf1e8ece',
 'Road damaged near the central square', 1, 0),

(1, ST_SetSRID(ST_Point(30.5209, 50.4531), 4326), '1f9588e6-a558-454b-92b5-48f5cf1e8ece',
 'Broken street light in the pedestrian area', 2, 0),

(1, ST_SetSRID(ST_Point(26.5861, 48.6794), 4326), 'a9a18285-9558-4d44-ae39-c2856c3d3e2b',
 'Overflowing trash container behind residential block', 3, 1),

(1, ST_SetSRID(ST_Point(26.5842, 48.6781), 4326), 'a9a18285-9558-4d44-ae39-c2856c3d3e2b',
 'Illegal parking blocking the sidewalk', 1, 0),

(1, ST_SetSRID(ST_Point(26.9974, 49.4209), 4326), 'bf318435-4223-43ad-9c1a-b52bbb0deeea',
 'Water leak near apartment complex', 2, 2),

(1, ST_SetSRID(ST_Point(26.9957, 49.4223), 4326), 'bf318435-4223-43ad-9c1a-b52bbb0deeea',
 'Fallen tree blocking a small road', 4, 1);

select *
from reports