CREATE TYPE course_role AS ENUM ('STUDENT', 'TEACHER');

CREATE TABLE enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
    role_in_course course_role NOT NULL,
    enrolled_at TIMESTAMP DEFAULT now(),
    UNIQUE(user_id, course_id)
);