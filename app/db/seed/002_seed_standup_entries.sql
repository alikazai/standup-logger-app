-- db/seed/seed_standup_entries.sql

INSERT INTO standup_entries (user_id, date, yesterday, today, blockers) VALUES
(
  (SELECT id FROM users WHERE email = 'alice@example.com'),
  NOW() - INTERVAL '1 day',
  'Worked on user interface and fixed styling bugs.',
  'Start backend integration for login.',
  'None'
),
(
  (SELECT id FROM users WHERE email = 'bob@example.com'),
  NOW() - INTERVAL '1 day',
  'Reviewed pull requests and documented endpoints.',
  'Write unit tests and update API docs.',
  'Waiting on data from the frontend team.'
);
