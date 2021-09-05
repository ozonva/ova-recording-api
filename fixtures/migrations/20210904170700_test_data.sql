-- +goose Up
-- +goose StatementBegin
INSERT INTO appointments (user_id, name, description, start_time, end_time) VALUES
(1, 'Some appointment 1', 'Some description 1', to_timestamp('20210903-111111','YYYYMMDD-HH24MISS'), to_timestamp('20210904-111111','YYYYMMDD-HH24MISS')),
(2, 'Some appointment 2', 'Some description 2', to_timestamp('20210905-111111','YYYYMMDD-HH24MISS'), to_timestamp('20210906-111111','YYYYMMDD-HH24MISS')),
(3, 'Some appointment 3', 'Some description 3', to_timestamp('20210907-111111','YYYYMMDD-HH24MISS'), to_timestamp('20210908-111111','YYYYMMDD-HH24MISS'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
