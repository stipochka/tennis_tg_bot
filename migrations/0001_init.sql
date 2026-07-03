BEGIN;
    -- btree_gist нужен чтобы в GIST индексе исключения использовать
    -- равенство по скалярному court_id вместе с оператором && по диапозону
    CREATE EXTESION IF NOT EXISTS btree_gist;

    -- courts - доступные корты
    CREATE TABLE IF NOT EXISTS courts (
        id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        name text NOT NULL,
        open_time time NOT NULL DEFAULT '07:00',
        close_time time NOT NULL DEFAULT '23:00',
        address text NOT NULL,
        is_active boolean NOT NULL DEFAULT true,
    );

    -- users - пользователи Telegram. Здесь будут использоваться ПДн: имя и телефон.
    CREATE TABLE users (
        id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        telegram_id bigint NOT NULL UNIQUE, -- постоянный id из телеграмма
        is_admin boolean NOT NULL default false,
        created_at timestampz NOT NULL DEFAULT now(),
        updated_at timestampz NOT NULL DEFAULT now(),
    );


    -- reservations - таблица с бронями пользователей и админа
    -- kind='booking' - обычная бронь пользователя user_id обязателен
    -- kind='block' - слот закрыт админом (запись по звонку, тех работы)
    --                 user_id=NULL created_by admin
    CREATE TABLE reservations (
        id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        court_id bigint NOT NULL REFERENCES courts(id),
        user_id bigint REFERENCES users(user_id) ON DELETE SET NULL,
        kind text NOT NULL DEFAULT 'booking',
        during tstzrange NOT NULL,         -- [начало, конец), 14:00-16:00

    )

END
