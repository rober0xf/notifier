CREATE TYPE token_purpose AS ENUM('email_verification', 'password_reset');

CREATE TYPE user_role AS ENUM('user', 'admin');

CREATE TYPE transaction_type AS ENUM('expense', 'income', 'subscription');

CREATE TYPE category_type AS ENUM(
  'electronics',
  'entertainment',
  'education',
  'clothing',
  'work',
  'sports'
);

CREATE TYPE frequency_type AS ENUM('daily', 'weekly', 'monthly', 'yearly');
