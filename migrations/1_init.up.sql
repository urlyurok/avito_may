CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Определение ENUM типов
CREATE TYPE service_type AS ENUM (
    'CONSTRUCTION',
    'DELIVERY',
    'MANUFACTURING'  -- Исправлено опечатка
    );

CREATE TYPE tender_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CLOSED'
    );

CREATE TYPE bid_status AS ENUM (
    'CREATED',
    'PUBLISHED',
    'CANCELED',
    'APPROVED',
    'REJECTED'
    );

CREATE TYPE author_type AS ENUM (
    'ORGANIZATION',
    'USER'
    );

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
    );

-- Создание таблиц
CREATE TABLE IF NOT EXISTS organization (
                                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                            name VARCHAR(100) NOT NULL,
                                            description TEXT,
                                            type organization_type,
                                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS employee (
                                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                        username VARCHAR(50) UNIQUE NOT NULL,
                                        first_name VARCHAR(50),
                                        last_name VARCHAR(50),
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tender (
                                      id SERIAL PRIMARY KEY,
                                      name VARCHAR(100),
                                      description TEXT,
                                      service_type service_type,
                                      status tender_status DEFAULT 'CREATED',
                                      organization_id UUID REFERENCES organization(id) ON DELETE CASCADE, -- Исправлен тип данных на UUID
                                      version INT DEFAULT 1,
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bid (
                                   id SERIAL PRIMARY KEY,
                                   name VARCHAR(100),
                                   description TEXT,
                                   status bid_status DEFAULT 'CREATED',
                                   tender_id INT REFERENCES tender(id) ON DELETE CASCADE,
                                   author_type author_type,  -- Исправлен тип данных
                                   author_id UUID REFERENCES employee(id) ON DELETE CASCADE, -- Исправлен тип данных на UUID
                                   version INT DEFAULT 1,
                                   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bid_review (
                                          id SERIAL PRIMARY KEY,
                                          description TEXT,
                                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
                                                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                                        organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
                                                        user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);