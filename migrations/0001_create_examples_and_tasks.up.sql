CREATE TABLE IF NOT EXISTS examples (
    id VARCHAR(64) PRIMARY KEY,
    expression TEXT NOT NULL,
    response VARCHAR(64) NOT NULL, -- variable
    calculated BOOLEAN NOT NULL DEFAULT FALSE,
    result DOUBLE PRECISION,
    user_id VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индекс по пользователю — для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_examples_user_id ON examples(user_id);

CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(64) PRIMARY KEY,
    example_id VARCHAR(64) NOT NULL REFERENCES examples(id) ON DELETE CASCADE,
    
    -- Числовые значения (уже разрешённые)
    value1 FLOAT8 NOT NULL,
    value2 FLOAT8 NOT NULL,
    result FLOAT8 NOT NULL,
    
    -- Операция
    sign VARCHAR(2) NOT NULL,
    
    -- Имя переменной, куда сохранён результат
    variable VARCHAR(64) NOT NULL,
    
    -- Порядок шага
    "order" INTEGER NOT NULL,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_example_id ON tasks(example_id);
CREATE INDEX IF NOT EXISTS idx_tasks_variable ON tasks(variable);
CREATE INDEX IF NOT EXISTS idx_tasks_order ON tasks("order");

