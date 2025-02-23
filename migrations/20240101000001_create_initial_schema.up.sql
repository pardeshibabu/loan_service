-- 1. borrowers
CREATE TABLE borrowers (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    credit_score INT
);

-- 2. employees
CREATE TABLE employees (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    role ENUM('field_officer', 'field_validator', 'admin') NOT NULL,
    status ENUM('active', 'inactive') NOT NULL DEFAULT 'active'
);

-- 3. investors
CREATE TABLE investors (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    kyc_status ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'pending',
    kyc_documents TEXT
);

-- 4. loans
CREATE TABLE loans (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    borrower_id BIGINT NOT NULL,
    principal_amount DECIMAL(15,2) NOT NULL,
    rate DECIMAL(5,2) NOT NULL,
    roi DECIMAL(5,2) NOT NULL,
    status ENUM('proposed', 'approved', 'invested', 'disbursed') NOT NULL DEFAULT 'proposed',
    approval_date DATETIME,
    field_validator_id BIGINT,
    validation_proof_url TEXT,
    disbursement_date DATETIME,
    field_officer_id BIGINT,
    signed_agreement_url TEXT,
    agreement_letter_url TEXT,
    FOREIGN KEY (borrower_id) REFERENCES borrowers(id),
    FOREIGN KEY (field_validator_id) REFERENCES employees(id),
    FOREIGN KEY (field_officer_id) REFERENCES employees(id)
);

-- 5. investments
CREATE TABLE investments (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    loan_id BIGINT NOT NULL,
    investor_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    status ENUM('active', 'completed', 'cancelled') NOT NULL DEFAULT 'active',
    investment_date DATETIME,
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    FOREIGN KEY (investor_id) REFERENCES investors(id)
);

-- 6. loan_state_history
CREATE TABLE loan_state_histories (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    loan_id BIGINT NOT NULL,
    changed_by_id BIGINT NOT NULL,
    from_status ENUM('proposed', 'approved', 'invested', 'disbursed'),
    to_status ENUM('proposed', 'approved', 'invested', 'disbursed') NOT NULL,
    changed_at DATETIME NOT NULL,
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    FOREIGN KEY (changed_by_id) REFERENCES employees(id)
);

-- 7. documents
CREATE TABLE documents (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    loan_id BIGINT NOT NULL,
    document_type ENUM('kyc', 'validation_proof', 'agreement_letter', 'signed_agreement') NOT NULL,
    file_url TEXT NOT NULL,
    uploaded_by_id BIGINT NOT NULL,
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    FOREIGN KEY (uploaded_by_id) REFERENCES employees(id)
);

-- 8. notifications
CREATE TABLE notifications (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    loan_id BIGINT NOT NULL,
    investor_id BIGINT,
    type ENUM('agreement_letter', 'investment_confirmation', 'disbursement_notice') NOT NULL,
    status ENUM('pending', 'sent', 'failed') DEFAULT 'pending',
    email_content TEXT,
    sent_at DATETIME,
    created_at DATETIME NOT NULL,
    
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    FOREIGN KEY (investor_id) REFERENCES investors(id)
);

-- 9. KYC history
CREATE TABLE kyc_histories (
    id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    investor_id BIGINT NOT NULL,
    from_status ENUM('pending', 'approved', 'rejected') NOT NULL,
    to_status ENUM('pending', 'approved', 'rejected') NOT NULL,
    reviewer_id BIGINT NOT NULL,
    comments TEXT,
    reviewed_at DATETIME NOT NULL,
    FOREIGN KEY (investor_id) REFERENCES investors(id),
    FOREIGN KEY (reviewer_id) REFERENCES employees(id)
); 