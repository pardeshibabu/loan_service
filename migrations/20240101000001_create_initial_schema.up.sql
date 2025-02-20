-- 1. borrowers
CREATE TABLE borrowers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    id_number VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- 2. employees
CREATE TABLE employees (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role ENUM('field_validator', 'field_officer', 'admin') NOT NULL,
    status ENUM('active', 'inactive') DEFAULT 'active',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- 3. investors
CREATE TABLE investors (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    kyc_status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
    kyc_documents TEXT,
    total_investment_amount DECIMAL(15,2) DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- 4. loans
CREATE TABLE loans (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    borrower_id BIGINT NOT NULL,
    principal_amount DECIMAL(15,2) NOT NULL,
    rate DECIMAL(5,2) NOT NULL,
    roi DECIMAL(5,2) NOT NULL,
    status ENUM('proposed', 'approved', 'invested', 'disbursed') DEFAULT 'proposed',
    approval_date DATETIME,
    field_validator_id BIGINT,
    validation_proof_url TEXT,
    disbursement_date DATETIME,
    field_officer_id BIGINT,
    signed_agreement_url TEXT,
    agreement_letter_url TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    
    FOREIGN KEY (borrower_id) REFERENCES borrowers(id),
    FOREIGN KEY (field_validator_id) REFERENCES employees(id),
    FOREIGN KEY (field_officer_id) REFERENCES employees(id)
);

-- 5. investments
CREATE TABLE investments (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    loan_id BIGINT NOT NULL,
    investor_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    agreement_letter_url TEXT,
    status ENUM('active', 'completed', 'cancelled') DEFAULT 'active',
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    FOREIGN KEY (investor_id) REFERENCES investors(id)
);

-- 6. loan_state_history
CREATE TABLE loan_state_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    loan_id BIGINT NOT NULL,
    changed_by_id BIGINT NOT NULL,
    from_status VARCHAR(20),
    to_status VARCHAR(20) NOT NULL,
    changed_at DATETIME NOT NULL,
    
    FOREIGN KEY (loan_id) REFERENCES loans(id),
    FOREIGN KEY (changed_by_id) REFERENCES employees(id)
);

-- 7. documents
CREATE TABLE documents (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    loan_id BIGINT NOT NULL,
    document_type ENUM('validation_proof', 'agreement_letter', 'signed_agreement', 'kyc') NOT NULL,
    file_url TEXT NOT NULL,
    uploaded_by_id BIGINT NOT NULL,
    created_at DATETIME NOT NULL,
    
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