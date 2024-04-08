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
