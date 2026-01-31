-- 创建数据库
CREATE DATABASE IF NOT EXISTS file_manager DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE file_manager;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '加密密码',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 文档表
CREATE TABLE IF NOT EXISTS documents (
    id VARCHAR(36) PRIMARY KEY COMMENT '文档UUID',
    file_name VARCHAR(255) NOT NULL COMMENT '原始文件名',
    storage_path VARCHAR(500) NOT NULL COMMENT '存储路径',
    file_size BIGINT NOT NULL COMMENT '文件大小(字节)',
    file_type VARCHAR(100) NOT NULL COMMENT 'MIME类型',
    file_extension VARCHAR(20) NOT NULL COMMENT '文件扩展名',
    md5_hash VARCHAR(32) COMMENT '文件MD5',
    upload_id VARCHAR(36) COMMENT '分片上传ID',
    status TINYINT DEFAULT 1 COMMENT '状态:1-正常,0-已删除',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_created_at (created_at),
    INDEX idx_status (status),
    INDEX idx_upload_id (upload_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档表';

-- 分片上传表
CREATE TABLE IF NOT EXISTS upload_chunks (
    upload_id VARCHAR(36) PRIMARY KEY COMMENT '上传ID',
    file_name VARCHAR(255) NOT NULL COMMENT '文件名',
    file_size BIGINT NOT NULL COMMENT '文件总大小',
    chunk_size INT NOT NULL COMMENT '分片大小',
    total_chunks INT NOT NULL COMMENT '总分片数',
    uploaded_chunks INT DEFAULT 0 COMMENT '已上传分片数',
    storage_path VARCHAR(500) COMMENT '存储路径',
    status TINYINT DEFAULT 1 COMMENT '状态:1-上传中,2-已完成,0-已取消',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    expired_at TIMESTAMP COMMENT '过期时间',
    INDEX idx_status (status),
    INDEX idx_expired_at (expired_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分片上传表';

-- 插入测试用户（密码: admin123, 使用 bcrypt 加密）
-- 注意：这个密码是测试用的，生产环境请修改
INSERT INTO users (username, password) VALUES
('admin', '$2a$10$7a.9QIsXq8IpU8PEUw6sp.ynvZ.eH2NMh.5Lu/6LtnUZ9RMm4DuYG')
ON DUPLICATE KEY UPDATE username=username;
