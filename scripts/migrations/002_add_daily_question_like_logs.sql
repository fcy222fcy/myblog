-- 迁移脚本：新增 daily_question_like_logs 表
-- 用途：为「已存在」的博客数据库补充「每日一问点赞记录」表，
--       支撑 DailyQuestionLikeLog 实体的按 IP 防重复点赞功能（事务 + 唯一索引）。
--
-- 说明：
--   - 使用 IF NOT EXISTS，可重复安全执行（幂等），不会对已存在的表造成破坏。
--   - 唯一索引 (question_id, visitor_ip) 保证同一访客对同一问题只能成功点赞一次；
--     并发请求下由数据库唯一约束兜底，多余插入被 OnConflict DoNothing 忽略。
--   - 外键关联 daily_questions(id)；删除/批量删除问题时由业务层事务先清理点赞记录
--     （见 DailyQuestionRepository.Delete / BatchDelete），再删问题本身。
--
-- 执行方式（任选其一，需先 USE 到博客库或指定库名）：
--   mysql -u<user> -p<password> blog < scripts/migrations/002_add_daily_question_like_logs.sql
--   或在数据库客户端中直接运行本文件内容。

CREATE TABLE IF NOT EXISTS daily_question_like_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    question_id BIGINT UNSIGNED NOT NULL,
    visitor_ip VARCHAR(50) NOT NULL,
    created_at DATETIME(3) NULL,
    UNIQUE KEY uk_question_ip (question_id, visitor_ip),
    INDEX idx_daily_question_like_question_id (question_id),
    INDEX idx_daily_question_like_visitor_ip (visitor_ip),
    CONSTRAINT fk_daily_question_like_question FOREIGN KEY (question_id) REFERENCES daily_questions(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
