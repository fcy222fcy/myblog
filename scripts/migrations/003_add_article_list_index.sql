-- 迁移脚本：为文章列表热查询补齐复合索引
-- 用途：为「已存在」的博客数据库补充 articles(status, is_top, created_at) 复合索引。
--
-- 背景：
--   前台文章列表 / 分类页 / 搜索的主查询形如
--     SELECT ... FROM articles WHERE status = 'published'
--     ORDER BY is_top DESC, created_at DESC LIMIT ? OFFSET ?
--   库中当前只有 status 单列索引（idx_articles_status），MySQL 需先按 status 过滤，
--   再对结果集做 filesort。建成 (status, is_top, created_at) 复合索引后，ORDER BY
--   的两列与索引列顺序一致，可直接反向扫描索引拿到有序结果，省掉排序步骤。
--
-- 说明：
--   - 本脚本只补这一条索引，幂等：先查 information_schema 判断是否已存在，不存在才 ALTER。
--   - article_tags 无需处理：tag_id 上已由外键 fk_article_tags_tag 隐式建了索引
--     （MySQL 会为外键列自动创建索引），再补一条属于重复索引，只会增加写放大。
--   - InnoDB 支持 Online DDL，加索引期间不阻塞读写，但会产生一次表的重建开销，
--     数据量大时建议低峰执行。
--
-- 执行方式（需先 USE 到博客库或指定库名）：
--   mysql -u<user> -p<password> blog < scripts/migrations/003_add_article_list_index.sql

SET @idx_exists := (
    SELECT COUNT(*) FROM information_schema.statistics
    WHERE table_schema = DATABASE()
      AND table_name = 'articles'
      AND index_name = 'idx_articles_status_top_created'
);
SET @ddl := IF(@idx_exists = 0,
    'ALTER TABLE articles ADD INDEX idx_articles_status_top_created (status, is_top, created_at)',
    'SELECT 1');
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 验证：
--   SHOW INDEX FROM articles WHERE Key_name = 'idx_articles_status_top_created';
--   期望看到三行，Seq_in_index 依次为 1=status、2=is_top、3=created_at。
