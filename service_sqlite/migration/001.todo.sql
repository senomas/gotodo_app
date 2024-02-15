PRAGMA foreign_keys = ON;
PRAGMA integrity_check;

CREATE TABLE IF NOT EXISTS todo_category (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS todo (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  category_id INTEGER NOT NULL,
  done BOOLEAN NOT NULL DEFAULT FALSE,
  FOREIGN KEY (category_id) REFERENCES todo_category (id)
);
	
