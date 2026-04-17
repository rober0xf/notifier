CREATE TABLE payments (
  id serial PRIMARY KEY,
  user_id integer NOT NULL,
  name text NOT NULL,
  amount numeric(10, 2) NOT NULL,
  type transaction_type NOT NULL,
  category category_type NOT NULL,
  date text NOT NULL,
  due_date text,
  paid bool NOT NULL,
  paid_at text,
  recurrent bool NOT NULL,
  frequency frequency_type,
  receipt_url text
)
