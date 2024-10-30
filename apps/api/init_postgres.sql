-- Connect to the newly created or existing database
\c news_and_topic_management;

-- Drop the existing tables if they exist
DROP TABLE IF EXISTS news_topic CASCADE;
DROP TABLE IF EXISTS topic CASCADE;
DROP TABLE IF EXISTS news CASCADE;
DROP TABLE IF EXISTS author CASCADE;

-- Table structure for table `news`
CREATE TABLE news
(
    id         SERIAL PRIMARY KEY,
    title      VARCHAR(45) NOT NULL,
    content    TEXT        NOT NULL,
    author_id  INTEGER   DEFAULT 0,
    status     VARCHAR(20) NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Inserting realistic data for table `news`
INSERT INTO news (id, title, content, author_id, status, updated_at, created_at)
VALUES (1, 'Health Benefits of a Mediterranean Diet',
        '<p>The Mediterranean diet has been associated with various health benefits...</p>',
        1, 'published', '2024-10-28 09:15:00', '2024-10-28 09:00:00'),
       (2, 'AI and the Future of Work',
        '<p>Artificial Intelligence is transforming the workplace by automating tasks...</p>',
        2, 'published', '2024-10-27 14:30:00', '2024-10-27 14:00:00'),
       (3, '10 Best Travel Destinations for 2024',
        '<p>Looking to plan your 2024 vacation? Here are ten must-visit destinations...</p>',
        3, 'draft', '2024-10-26 11:20:00', '2024-10-26 11:00:00'),
       (4, 'Climate Change: What You Can Do to Help',
        '<p>As global temperatures continue to rise, individuals have the power to make a difference...</p>',
        4, 'published', '2024-10-25 08:10:00', '2024-10-25 08:00:00'),
       (5, '5 Tips for Boosting Your Mental Health',
        '<p>Prioritizing mental health is essential. Here are five tips to improve your well-being...</p>',
        1, 'archived', '2024-10-24 13:45:00', '2024-10-24 13:00:00'),
       (6, 'Advancements in Renewable Energy Technologies',
        '<p>Renewable energy sources such as solar and wind are becoming more efficient and accessible...</p>',
        2, 'published', '2024-10-23 10:30:00', '2024-10-23 10:00:00');

-- Table structure for table `topic`
CREATE TABLE topic
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Inserting data for table `topic`
INSERT INTO topic (id, name)
VALUES (1, 'Health'),
       (2, 'Technology'),
       (3, 'Travel'),
       (4, 'Environment'),
       (5, 'Mental Health');

-- Table structure for table `news_topic` (associates `news` with `topic`)
CREATE TABLE news_topic
(
    id       SERIAL PRIMARY KEY,
    news_id  INTEGER NOT NULL REFERENCES news (id) ON DELETE CASCADE,
    topic_id INTEGER NOT NULL REFERENCES topic (id) ON DELETE CASCADE,
    UNIQUE (news_id, topic_id)
);

-- Inserting data for table `news_topic`
INSERT INTO news_topic (news_id, topic_id)
VALUES (1, 1),
       (1, 4),
       (2, 2),
       (3, 3),
       (4, 4),
       (5, 5),
       (6, 2);

-- Table structure for table `author`
CREATE TABLE author
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(200) DEFAULT '',
    created_at TIMESTAMP    DEFAULT NOW(),
    updated_at TIMESTAMP    DEFAULT NOW()
);

-- Inserting data for table `author`
INSERT INTO author (id, name, created_at, updated_at)
VALUES (1, 'Doni', '2017-05-18 13:50:19', '2017-05-18 13:50:19'),
       (2, 'Deni', '2017-05-19 14:00:00', '2017-05-19 14:00:00'),
       (3, 'Dani', '2017-05-20 15:00:00', '2017-05-20 15:00:00'),
       (4, 'Dini', '2017-05-21 16:00:00', '2017-05-21 16:00:00');
