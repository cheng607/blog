package com.example.blog.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.example.blog.entity.Article;

import java.util.List;

public interface ArticleService extends IService<Article> {
    // 保存文章（含标签关联）
    Integer saveArticle(Article article);
    List<Article> getArticleList();
    // 根据用户ID查询文章
    List<Article> getArticleByUserId(Integer userId);

    // 根据ID查询文章（含标签）
    Article getArticleById(Integer id);

    // 根据标签ID查询文章
    List<Article> getArticleByTagId(Integer tagId);

    // 更新文章（含标签关联）
    Article updateArticle(Article article);

    // 删除文章（数据库会级联删除article_tag关联记录）
    void deleteArticleById(Integer id);

    Integer incrementViewCount(Integer id);
}