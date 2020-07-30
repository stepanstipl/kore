INSERT INTO identities (user_id, provider_email, provider) SELECT id, email, 'sso' FROM users u WHERE u.id != 1;
