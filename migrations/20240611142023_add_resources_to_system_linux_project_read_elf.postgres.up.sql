INSERT INTO public.bookmarks (url, title, description, user_id, project_id)
SELECT 'https://example.com/readelf-resource', 'Readelf Resource', 'Description of Readelf Resource', 67, p.id
FROM public.projects p
JOIN public.categories c ON p.category_id = c.id
WHERE p.name = 'readelf' AND c.category = 'system-linux';


