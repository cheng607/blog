package com.example.blog.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.example.blog.entity.ArticleTag;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

import java.util.List;

@Mapper
public interface ArticleTagMapper extends BaseMapper<ArticleTag> {
    // 根据文章ID删除关联的标签
    void deleteByArticleId(@Param("articleId") Integer articleId);

    // 批量插入文章与标签的关联
    void batchInsert(@Param("list") List<ArticleTag> articleTags);
}