--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Debian 14.5-1.pgdg110+1)
-- Dumped by pg_dump version 14.5 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

-- users table
CREATE TABLE public.users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at timestamp without time zone default CURRENT_TIMESTAMP,
    updated_at timestamp without time zone default CURRENT_TIMESTAMP
);

-- categories table
CREATE TABLE public.categories (
    id SERIAL PRIMARY KEY,
    category VARCHAR(255) NOT NULL UNIQUE,
    created_at timestamp without time zone default CURRENT_TIMESTAMP,
    updated_at timestamp without time zone default CURRENT_TIMESTAMP
);

-- Insert the projects categories
INSERT INTO public.categories (category) VALUES ('system-linux'), ('system-algorithms'), ('blockchain'), ('malloc'), ('simple-shell');

CREATE TABLE public.projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    category_id INTEGER NOT NULL,
    created_at timestamp without time zone default CURRENT_TIMESTAMP,
    updated_at timestamp without time zone default CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES public.categories (id) ON DELETE SET NULL
);

CREATE TABLE public.bookmarks (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    title VARCHAR(255),
    description TEXT,
    user_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,
    created_at timestamp without time zone default CURRENT_TIMESTAMP,
    updated_at timestamp without time zone default CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES public.users (id) ON DELETE SET NULL,
    FOREIGN KEY (project_id) REFERENCES public.projects (id) ON DELETE SET NULL
);

-- Create ratings table
CREATE TABLE public.ratings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    bookmark_id INTEGER NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    created_at timestamp without time zone default CURRENT_TIMESTAMP,
    updated_at timestamp without time zone default CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES public.users (id) ON DELETE CASCADE,
    FOREIGN KEY (bookmark_id) REFERENCES public.bookmarks (id) ON DELETE CASCADE,
    UNIQUE (user_id, bookmark_id) -- Ensures a user can only rate a bookmark once
);

-- index to improve requests performance
CREATE INDEX idx_user_id ON public.bookmarks (user_id);
CREATE INDEX idx_project_id ON public.bookmarks (project_id);
CREATE INDEX idx_user_id_bookmark_id ON public.ratings (user_id, bookmark_id);

-- categories are  ('system-linux'), ('system-algorithms'), ('blockchain'), ('malloc'), ('simple-shell');

-- Insert the Holbies projects
INSERT INTO public.projects (name, category_id)
VALUES
('libasm', (SELECT id FROM public.categories WHERE category = 'system-linux')),
('ls', (SELECT id FROM public.categories WHERE category = 'system-linux')),
('multithreading', (SELECT id FROM public.categories WHERE category = 'system-linux')),
('nm-objdump', (SELECT id FROM public.categories WHERE category = 'system-linux')),
('proc_filesystem', (SELECT id FROM public.categories WHERE category = 'system-linux')),
('readelf', (SELECT id FROM public.categories WHERE category = 'system-linux')),
('signals', (SELECT id FROM public.categories WHERE category = 'system-linux')),
('sockets', (SELECT id FROM public.categories WHERE category = 'system-linux')),
('strace', (SELECT id FROM public.categories WHERE category = 'system-linux')),
('graphs', (SELECT id FROM public.categories WHERE category = 'system-algorithms')),
('huffman_coding', (SELECT id FROM public.categories WHERE category = 'system-algorithms')),
('nary_trees and red-black trees', (SELECT id FROM public.categories WHERE category = 'system-algorithms')),
('pathfinding', (SELECT id FROM public.categories WHERE category = 'system-algorithms')),
('blockchain', (SELECT id FROM public.categories WHERE category = 'blockchain')),
('crypto', (SELECT id FROM public.categories WHERE category = 'blockchain')),
('malloc project', (SELECT id FROM public.categories WHERE category = 'malloc')),
('The shell project', (SELECT id FROM public.categories WHERE category = 'simple-shell'));






