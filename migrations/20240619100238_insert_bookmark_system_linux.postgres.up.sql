INSERT INTO public.bookmarks (url, type, description, user_id, project_id)
VALUES 
('https://www.youtube.com/watch?v=wLXIWKUWpSs&ab_channel=DavyWybiral','video', 'Assembly explained easy', 75, (SELECT id FROM public.projects WHERE name = 'libasm')),
('https://www.youtube.com/watch?v=s3o5tixMFho&ab_channel=CodingOverflow', 'video', 'video tutorial on programming sockets in C ', 75, (SELECT id FROM public.projects WHERE name = 'sockets'));

