/* Менеджеры. ФИО + компания, которую представляет */
CREATE TABLE IF NOT EXISTS managers
(
    id         bigserial
        CONSTRAINT managers_pk PRIMARY KEY,
    name       varchar(255) not null,
    surname    varchar(255) not null,
    patronymic varchar(255),
    company_id bigint,
    created_at timestamp default now(),
    updated_at timestamp default now()
);

/* Вакансии. Менеджер может указать не свою компанию в вакансии (например, если работает без компании или на несколько) */
CREATE TABLE IF NOT EXISTS vacancies
(
    id            bigserial
        CONSTRAINT vacancies_pk PRIMARY KEY,
    manager_id    bigint       not null,
    company_id    bigint,
    position      varchar(255) not null,       -- должность
    description   text,                        -- описание вакансии
    salary        int4range,
    skills        varchar(255)[],              -- необходимые навыки (возможно, стоило бы повесить индекс)
    city          varchar(255),                -- город расположение
    work_schedule varchar(255),
    status        varchar(255) default 'Open', -- актуальность вакансии
    created_at    timestamp    default now(),
    updated_at    timestamp    default now()
);

/* Соискатели */
CREATE TABLE IF NOT EXISTS applicants
(
    id                       bigserial
        CONSTRAINT applicants_pk PRIMARY KEY,
    name                     varchar(255)       not null,
    surname                  varchar(255)       not null,
    patronymic               varchar(255),
    speciality               varchar(255)       not null,       -- специальность
    years_of_work_experience smallint,
    phone_number             varchar(50) unique not null,
    birthday_date            date,                              -- для определения возраста
    gender                   varchar(255),
    city                     varchar(255),
    status                   varchar(255) default 'Unemployed', -- статус занятости
    created_at               timestamp    default now(),
    updated_at               timestamp    default now()
);

/* Отклики */
CREATE TABLE IF NOT EXISTS vacancy_application
(
    id           bigserial
        CONSTRAINT vacancy_applications_pk PRIMARY KEY,
    applicant_id bigint not null,
    vacancy_id   bigint not null,
    status       varchar(255) default 'Not accepted', -- решение менеджера
    created_at   timestamp    default now(),
    updated_at   timestamp    default now()
);

/* Доп. таблица для получения сведений о компаниях */
CREATE TABLE IF NOT EXISTS companies
(
    id               bigserial
        CONSTRAINT company_pk PRIMARY KEY,
    title            varchar(255) not null,
    foundation_date  date,
    specialization   varchar(255),
    resident_country varchar(255),
    created_at       timestamp default now(),
    updated_at       timestamp default now()
);

/* Для поиска всех менеджеров в компании. Например, соискатель ищет всех HR Ozon */
CREATE INDEX managers_company_id_index ON managers (company_id);

/* Для поиска всех вакансий в компании/у определенного менеджера.
   Например, соискатель ищет все вакансии в Ozon, HR ищет все свои вакансии */
CREATE INDEX vacancies_company_id_manager_id_index ON vacancies (company_id, manager_id);

/* Для поиска всех откликов по вакансии, по определенному соискателю.
   Использует и HR, и соискатель */
CREATE INDEX vacancy_applications_vacancy_id_applicant_id_index ON vacancy_application (vacancy_id, applicant_id);

/* Для поиска всех соискателей по специальности.
   Например, HR ищет всех соискателей по специальности Frontend */
CREATE INDEX applicants_speciality_index ON applicants (speciality);

