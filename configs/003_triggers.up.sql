-- Создание триггера для установки значения created_at при вставке новой записи
CREATE OR REPLACE FUNCTION set_created_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_created_at_trigger
BEFORE INSERT ON banners
FOR EACH ROW
EXECUTE FUNCTION set_created_at();

-- Создание триггера для обновления значения updated_at при вставке или обновлении записи
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_trigger
BEFORE INSERT OR UPDATE ON banners
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
