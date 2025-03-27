-- Veritabanı schema oluşturma

-- Extension oluşturma
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Kullanıcılar tablosu
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    avatar VARCHAR(255),
    bio TEXT,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Kategoriler tablosu
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(150) NOT NULL UNIQUE,
    description TEXT,
    parent_id INT REFERENCES categories(id) ON DELETE SET NULL,
    icon VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    order_number INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Etiketler tablosu
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Makale tablosu
CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(300) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    summary TEXT,
    featured_image VARCHAR(255),
    author_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    view_count INT DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    is_breaking BOOLEAN DEFAULT FALSE,
    is_editors_pick BOOLEAN DEFAULT FALSE,
    meta_title VARCHAR(150),
    meta_description VARCHAR(300),
    meta_keywords VARCHAR(300),
    allow_comments BOOLEAN DEFAULT TRUE,
    publish_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Makale-Etiket ilişki tablosu
CREATE TABLE IF NOT EXISTS article_tags (
    article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, tag_id)
);

-- Yorumlar tablosu
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    parent_id INT REFERENCES comments(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    name VARCHAR(100),
    email VARCHAR(150),
    ip_address VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Medya tablosu
CREATE TABLE IF NOT EXISTS media (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size INT NOT NULL,
    alt_text VARCHAR(255),
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Reklam alanları tablosu
CREATE TABLE IF NOT EXISTS ad_spaces (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(100) NOT NULL,
    content TEXT,
    image VARCHAR(255),
    url VARCHAR(500),
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'inactive',
    click_count INT DEFAULT 0,
    impression_count INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Sayfa tablosu
CREATE TABLE IF NOT EXISTS pages (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(300) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    author_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    meta_title VARCHAR(150),
    meta_description VARCHAR(300),
    meta_keywords VARCHAR(300),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- İlk admin kullanıcısı oluşturma (şifre: admin123)
INSERT INTO users (username, email, password, name, role)
VALUES ('admin', 'admin@habersite.com', '$2a$10$VgDVxrF8D2/VSuKshzCWEu9n3yB5SegvFiNP/.bIz/cWOqFcNX8H2', 'Admin Kullanıcı', 'admin');

-- İlk kategorileri oluşturma
INSERT INTO categories (name, slug, description)
VALUES 
('Gündem', 'gundem', 'Güncel haberler ve gelişmeler'),
('Ekonomi', 'ekonomi', 'Ekonomi haberleri ve piyasa gelişmeleri'),
('Spor', 'spor', 'Spor haberleri ve gelişmeleri'),
('Teknoloji', 'teknoloji', 'Teknoloji haberleri ve gelişmeler'),
('Sağlık', 'saglik', 'Sağlık haberleri ve bilgileri'),
('Kültür-Sanat', 'kultur-sanat', 'Kültür ve sanat haberleri'); 