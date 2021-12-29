DELETE FROM employees;

INSERT INTO employees (
    employee_id,
    auth0_id,
    email_address,
    first_name,
    middle_name,
    last_name,
    phone_number,
    birth_date,
    hire_date,
    picture,
    language,
    country,
    city,
    zipcode,
    salary,
    position
) VALUES
('bc4cd1a1-4f0e-4e39-9960-e6b1cfe388db', 'auth0|619517c51b1a2e00690365af','ivor@devpie.io','Ivor','Scott','Cummings','9083971266','1989-11-01','2018-09-01','','english','united states of america','Basking Ridge','07920','160,000','Senior Software Engineer'),
('35d27bb4-c39e-4c10-9a64-aabc2490ec4d', 'auth0|61cb8be92bb93500699c27c4','people@devpie.io','Pe','','Ople','15144558913','1984-2-01','2020-04-01','','german','germany','berlin','12049','67,000','Software Engineer');