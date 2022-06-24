CREATE TABLE users (
                       id    SERIAL PRIMARY KEY,
                       userId    int,
                       fullName  VARCHAR(100),
                       username  VARCHAR(100),
                       email VARCHAR(100),
                       phoneNumber VARCHAR(100),
                       status int
);
CREATE TABLE massageType (
                             id    SERIAL PRIMARY KEY,
                             mType VARCHAR(100)
);
INSERT INTO massageType(mType) VALUES('💆 Шейно воротниковый массаж');
INSERT INTO massageType(mType) VALUES('🧖 Лечебный массаж');

CREATE TABLE massageSchedule (
                                 id    SERIAL PRIMARY KEY,
                                 mid int references massageType(id),
                                 mDate DATE DEFAULT CURRENT_DATE,
                                 mTime TIME,
                                 uId int references users(id),
                                 status int,
                                 isCanceled bool default false
);