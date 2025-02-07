use ozge_contract;

CREATE TABLE TOO (
                     id INT AUTO_INCREMENT PRIMARY KEY,
                     name VARCHAR(255) NOT NULL,
                     bin VARCHAR(20) NOT NULL,
                     registration_file TEXT NOT NULL,
                     ceo_name VARCHAR(255) NOT NULL,
                     ceo_order_file TEXT NOT NULL,
                     ceo_id_file TEXT NOT NULL,
                     representative_poa TEXT,
                     representative_id TEXT NOT NULL,
                     bank_details TEXT NOT NULL,
                     legal_address TEXT NOT NULL,
                     actual_address TEXT NOT NULL,
                     contact_details TEXT NOT NULL,
                     email VARCHAR(255) NOT NULL,
                     egov_file TEXT NOT NULL,
                     company_card TEXT NOT NULL,
                     company_code VARCHAR(255) NOT NULL
);

CREATE TABLE IP (
                    id INT AUTO_INCREMENT PRIMARY KEY,
                    name VARCHAR(255) NOT NULL,
                    iin VARCHAR(20) NOT NULL,
                    registration_file TEXT NOT NULL,
                    representative_poa TEXT,
                    representative_id TEXT NOT NULL,
                    bank_details TEXT NOT NULL,
                    legal_address TEXT NOT NULL,
                    actual_address TEXT NOT NULL,
                    contact_details TEXT NOT NULL,
                    email VARCHAR(255) NOT NULL,
                    company_card TEXT NOT NULL,
                    company_code VARCHAR(255) NOT NULL
);

CREATE TABLE Individual (
                            id INT AUTO_INCREMENT PRIMARY KEY,
                            full_name VARCHAR(255) NOT NULL,
                            iin VARCHAR(20) NOT NULL,
                            id_file TEXT NOT NULL,
                            bank_details TEXT NOT NULL,
                            legal_address TEXT NOT NULL,
                            actual_address TEXT NOT NULL,
                            contact_details TEXT NOT NULL,
                            email VARCHAR(255) NOT NULL,
                            company_code VARCHAR(255) NOT NULL
);
