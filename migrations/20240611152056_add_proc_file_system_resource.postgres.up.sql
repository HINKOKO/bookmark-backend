INSERT INTO public.bookmarks (url, title, description, user_id, project_id)
SELECT 'https://medium.com/geekculture/linux-proc-pid-directory-part-five-a10dacf49b4a', 'Proc filesystem', 'Excellent blog about proc', 67, p.id
FROM public.projects p
JOIN public.categories c ON p.category_id = c.id
WHERE p.name = 'proc_filesystem' AND c.category = 'system-linux';
