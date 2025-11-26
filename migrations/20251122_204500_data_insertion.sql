INSERT INTO cities (name, region, location)
VALUES ('Kyiv', 'Kyiv', ST_SetSRID(ST_Point(30.5238, 50.4547), 4326)),
       ('Kamianets-Podilskyi', 'Khmelnytskyi', ST_SetSRID(ST_Point(26.58516, 48.67882), 4326)),
       ('Khmelnytskyi', 'Khmelnytskyi', ST_SetSRID(ST_Point(26.9965, 49.4216), 4326));


INSERT INTO users (email, password_hash, name, surname, phone_number, reg_city_id)
VALUES ('admin@gmail.com', '$2a$10$8kRRzb7VDH0ph/T.D.v4G.l/5EBm9S6AZWhGh9VDjgP3.oXfydmJy', 'Admin', 'Admin', 'Tphonenumber', (SELECT city_id FROM cities WHERE name = 'Kyiv'));


INSERT INTO reports (user_id, location, city_id, description, category_id, current_status_id)
VALUES ((SELECT user_id FROM users WHERE email = 'admin@gmail.com'),
        ST_SetSRID(ST_Point(30.5245, 50.4552), 4326),
        (SELECT city_id FROM cities WHERE name = 'Kyiv'),
        'Road damaged near the central square',
        1,
        0),

       ((SELECT user_id FROM users WHERE email = 'admin@gmail.com'),
        ST_SetSRID(ST_Point(30.5209, 50.4531), 4326),
        (SELECT city_id FROM cities WHERE name = 'Kyiv'),
        'Broken street light in the pedestrian area',
        2,
        0),

       ((SELECT user_id FROM users WHERE email = 'admin@gmail.com'),
        ST_SetSRID(ST_Point(26.5861, 48.6794), 4326),
        (SELECT city_id FROM cities WHERE name = 'Kamianets-Podilskyi'),
        'Overflowing trash container behind residential block',
        3,
        1),

       ((SELECT user_id FROM users WHERE email = 'admin@gmail.com'),
        ST_SetSRID(ST_Point(26.5842, 48.6781), 4326),
        (SELECT city_id FROM cities WHERE name = 'Kamianets-Podilskyi'),
        'Illegal parking blocking the sidewalk',
        1,
        0),

       ((SELECT user_id FROM users WHERE email = 'admin@gmail.com'),
        ST_SetSRID(ST_Point(26.9974, 49.4209), 4326),
        (SELECT city_id FROM cities WHERE name = 'Khmelnytskyi'),
        'Water leak near apartment complex',
        2,
        2),

       ((SELECT user_id FROM users WHERE email = 'admin@gmail.com'),
        ST_SetSRID(ST_Point(26.9957, 49.4223), 4326),
        (SELECT city_id FROM cities WHERE name = 'Khmelnytskyi'),
        'Fallen tree blocking a small road',
        4,
        1);

SELECT *
from reports
