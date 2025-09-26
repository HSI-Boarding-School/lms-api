CREATE TYPE status_type AS ENUM ('DRAFT', 'SUBMITTED', 'LOCKED');

CREATE TABLE log_books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    student_id_id UUID REFERENCES users(id) ON DELETE CASCADE,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    status status_type DEFAULT 'DRAFT',
    submited_at TIMESTAMP,
    locked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE log_book_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    log_book_id UUID REFERENCES log_books(id) ON DELETE CASCADE,
    entry_date DATE NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);