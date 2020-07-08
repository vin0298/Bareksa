DROP TABLE IF EXISTS news_articles, news_tags, tags;

CREATE TABLE news_articles (
    article_id SERIAL PRIMARY KEY,
    author VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    time_published TIMESTAMPTZ  NOT NULL,
    uuid UUID NOT NULL,
    topic VARCHAR(255) NOT NULL,
    status VARCHAR(25) NOT NULL
);

CREATE TABLE tags_tw (
    tag_id SERIAL PRIMARY KEY,
    tag_name VARCHAR(255) NOT NULL unique,
    uuid UUID NOT NULL
);

CREATE TABLE news_tags_tw (
    article_id INT NOT NULL REFERENCES news_articles(article_id) ON UPDATE CASCADE ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(tag_id) ON UPDATE CASCADE ON DELETE CASCADE
);
