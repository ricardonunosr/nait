CREATE TABLE IF NOT EXISTS staff(
    username TEXT PRIMARY KEY,
    email TEXT,
    firstname TEXT,
    lastname TEXT,
    password TEXT
);

INSERT INTO
    staff (username, email, firstname, lastname, password)
VALUES
    (
        "mariaribeiro3",
        "maria@gmail.com",
        "Maria",
        "Ribeiro",
        "123"
    ),
    (
        "inessilva_",
        "ines@gmail.com",
        "Ines",
        "Silva",
        "123"
    ),
    (
        "joaocorreia99",
        "joao@gmail.com",
        "Joao",
        "Correia",
        "123"
    );

CREATE TABLE IF NOT EXISTS events_name(
    event_id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_name TEXT NOT NULL
);

INSERT INTO
    events_name (event_name)
VALUES
    ("Funk Brasileiro"),
    ("80s Classics"),
    ("Baile das Novinhas");

CREATE TABLE IF NOT EXISTS events(
    event_date DATE PRIMARY KEY,
    event_id TEXT,
    event_payment_url TEXT,
    CONSTRAINT fk_event_id FOREIGN KEY(event_id) REFERENCES events_name(event_id) ON DELETE CASCADE
);

INSERT INTO
    events (event_date, event_id, event_payment_url)
VALUES
    (
        "2023-10-18",
        1,
        "https://buy.stripe.com/test_9AQg0F6vV5Z8fGE4gh"
    ),
    (
        "2023-10-12",
        2,
        "https://buy.stripe.com/test_9AQg0F6vV5Z8fGE4gh"
    ),
    (
        "2023-10-02",
        3,
        "https://buy.stripe.com/test_9AQg0F6vV5Z8fGE4gh"
    );