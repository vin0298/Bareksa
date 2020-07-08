-- #1 PK
INSERT INTO news_articles(author, title, content, time_published, uuid, topic, status)
VALUES ('Vincent B', 'Article First', 'Loren', '2020-07-22 19:10:25-07',
        '5269c1b4-8665-4acc-a86a-855fa20122c8', 'Topic First', 'Draft');

-- #2 PK
INSERT INTO news_articles(author, title, content, time_published, uuid, topic, status)
VALUES ('Olivia', 'Article Second', 'Loren', '2020-07-21 19:10:25-07',
        '85fdc8f7-f564-487c-b41a-1487cd4e6847', 'Topic Second', 'Published');

-- #3 PK
INSERT INTO news_articles(author, title, content, time_published, uuid, topic, status)
VALUES ('Winson', 'Article Third', 'Loren', '2020-07-23 19:10:25-07',
        '98ce7669-341b-444e-bdef-9066b3421f4d', 'Topic First', 'Draft');

-- #1 PK
INSERT INTO tags(tag_name, uuid) 
VALUES ('First Loren Tag', 'd316539e-dcc8-4b4e-b8e1-ee53d5f3529a');

-- #2 PK
INSERT INTO tags(tag_name, uuid) 
VALUES ('Sec Loren Tag', '4a4ef0bb-7f0a-4ee3-adf4-da0ca7aa3850');

-- #3 PK
INSERT INTO tags(tag_name, uuid) 
VALUES ('Third Loren Tag', '21561967-4d7b-42ae-9705-2137543f3eff');

INSERT INTO news_tags(article_id, tag_id)
VALUES(1, 1);

INSERT INTO news_tags(article_id, tag_id)
VALUES(1, 2);

INSERT INTO news_tags(article_id, tag_id)
VALUES(2, 2);

INSERT INTO news_tags(article_id, tag_id)
VALUES(3, 2);

INSERT INTO news_tags(article_id, tag_id)
VALUES(3, 1);
