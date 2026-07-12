-- +goose Up
-- +goose StateMentBegin



-- btree_gist нужен чтобы в GIST индексе исключения использовать
-- равенство по скалярному court_id вместе с оператором && по диапозону
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- courts - доступные корты
CREATE TABLE IF NOT EXISTS courts (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    open_time time NOT NULL DEFAULT '07:00',
    close_time time NOT NULL DEFAULT '23:00',
    address text NOT NULL,
    is_active boolean NOT NULL DEFAULT true
);

-- users - пользователи Telegram. Здесь будут использоваться ПДн: имя и телефон.
CREATE TABLE IF NOT EXISTS users (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    telegram_id bigint NOT NULL UNIQUE, -- постоянный id из телеграмма
    is_admin boolean NOT NULL default false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz
);


-- reservations - таблица с бронями пользователей и админа
-- kind='booking' - обычная бронь пользователя user_id обязателен
-- kind='block' - слот закрыт админом (запись по звонку, тех работы)
--                 user_id=NULL created_by admin
CREATE TABLE IF NOT EXISTS reservations (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    court_id bigint NOT NULL REFERENCES courts(id),
    user_id bigint REFERENCES users(id) ON DELETE SET NULL,
    kind text NOT NULL DEFAULT 'booking',
    during tstzrange NOT NULL,         -- [начало, конец), 14:00-16:00
    status text NOT NULL DEFAULT 'pending',
    created_at timestamptz NOT NULL DEFAULT now(),
    reviwed_at timestamptz,
    cancelled_at timestamptz,

    CONSTRAINT reservations_kind_chk
        CHECK (kind IN ('booking', 'block')),

    CONSTRAINT reservations_status_chk
        CHECK (status IN ('pending', 'confirmed', 'rejected', 'cancelled')),

    -- проверка диапозона бронированая
    CONSTRAINT reservations_during_chk
        CHECK (NOT isempty(during)
                AND lower(during) IS NOT NULL
                AND upper(during) IS NOT NULL),

    -- у обычной брони должен быть пользователь
    CONSTRAINT reservations_booking_user_chk
        CHECK ((kind = 'booking' AND user_id IS NOT NULL)
                OR kind = 'block'),

    -- ОСНОВНАЯ логика: два занимающих слот диапозона (pending/confirmed) на
    -- одном корте не пересекаются - это дает авто-сериализацию заявок
    CONSTRAINT reservation_no_overlap
        EXCLUDE USING gist (
            court_id WITH =,
            during WITH &&
        ) WHERE (status IN ('pending', 'confirmed'))
);


-- Активные (pending + confirmed брони пользователя); для /my и лимита броней
CREATE INDEX IF NOT EXISTS idx_reservations_user_active
    ON reservations (user_id)
    WHERE status IN ('pending', 'confirmed') AND kind = 'booking';

-- Очередь модерации; ожидающие решение админа
CREATE INDEX IF NOT EXISTS idx_reservations_pending
    ON reservations (created_at)
    WHERE status = 'pending' AND kind = 'booking';

-- Корт по умолчанию
INSERT INTO courts(name, open_time, close_time, address)
VALUES ('Корт Неймарк', '07:00', '23:00', 'Большие Овраги 12к17');

-- +goose StateMentEnd

-- +goose Down
-- +goose StateMentBegin
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS courts;
-- +goose StateMentEnd
