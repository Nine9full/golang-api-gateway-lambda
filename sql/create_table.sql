CREATE TABLE todos (
    id INT IDENTITY(1,1) PRIMARY KEY, -- Identity column with auto-increment
    name NVARCHAR(MAX),               -- Name column
    active BIT                        -- Represents the 'completed' field in the struct
);