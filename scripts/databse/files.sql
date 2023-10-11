CREATE TABLE files (
    id SERIAL,
    folder_id INT,
    owner_id INT NOT NULL,
    name VARCHAR(700) NOT NULL,
    type VARCHAR(50) NOT NULL,
    path VARCHAR(250) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP NOT NULL,
    deleted BOOL DEFAULT false,
    PRIMARY KEY (id),
    CONSTRAINT fk_files_folder_id FOREIGN KEY (folder_id) REFERENCES folders(id),
    CONSTRAINT fk_files_owner_id FOREIGN KEY (owner_id) REFERENCES users(id)
)