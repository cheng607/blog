package com.example.blog.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.UpdateWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.example.blog.entity.Article;
import com.example.blog.entity.ArticleTag;
import com.example.blog.entity.Tag;
import com.example.blog.mapper.ArticleMapper;
import com.example.blog.mapper.ArticleTagMapper;
import com.example.blog.mapper.TagMapper;
import com.example.blog.service.ArticleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

@Service
public class ArticleServiceImpl extends ServiceImpl<ArticleMapper, Article> implements ArticleService {

    @Autowired
    private ArticleTagMapper articleTagMapper;

    @Autowired
    private TagMapper tagMapper;

    @Override
    @Transactional
    public Integer saveArticle(Article article) {
        // 1. 设置时间字段（如果不使用自动填充，手动设置）
        LocalDateTime now = LocalDateTime.now();
        article.setPublishTime(now);
        article.setUpdateTime(now);

        // 2. 保存文章主体
        baseMapper.insert(article);
        Integer articleId = article.getId();

        // 3. 处理标签关联（使用tagIds）
        List<Integer> tagIds = article.getTagIds();
        if (tagIds != null && !tagIds.isEmpty()) {
            // 校验标签是否存在
            List<ArticleTag> articleTags = new ArrayList<>();
            for (Integer tagId : tagIds) {
                Tag existTag = tagMapper.selectById(tagId);
                if (existTag == null) {
                    throw new RuntimeException("标签ID不存在：" + tagId);
                }
                // 构建文章-标签关联对象
                ArticleTag at = new ArticleTag();
                at.setArticleId(articleId);
                at.setTagId(tagId);
                articleTags.add(at);
            }
            // 批量插入关联记录
            articleTagMapper.batchInsert(articleTags);
        }
        return articleId;
    }

    @Override
    public List<Article> getArticleList() {
        return baseMapper.selectAllWithTags();
    }

    @Override
    public List<Article> getArticleByUserId(Integer userId) {
        QueryWrapper<Article> queryWrapper = new QueryWrapper<>();
        queryWrapper.eq("user_id", userId);
        return baseMapper.selectList(queryWrapper);
    }

    @Override
    public Article getArticleById(Integer id) {
        return baseMapper.selectByIdWithTags(id);
    }

    @Override
    public List<Article> getArticleByTagId(Integer tagId) {
        return baseMapper.selectByTagId(tagId);
    }

    @Override
    @Transactional
    public Article updateArticle(Article article) {
        // 更新时间字段
        article.setUpdateTime(LocalDateTime.now());

        // 更新文章主体
        baseMapper.updateById(article);
        Integer articleId = article.getId();

        // 更新标签关联（先删除旧关联，再处理新的tagIds）
        articleTagMapper.deleteByArticleId(articleId);

        List<Integer> tagIds = article.getTagIds();
        if (tagIds != null && !tagIds.isEmpty()) {
            List<ArticleTag> articleTags = new ArrayList<>();
            for (Integer tagId : tagIds) {
                Tag existTag = tagMapper.selectById(tagId);
                if (existTag == null) {
                    throw new RuntimeException("标签ID不存在：" + tagId);
                }
                ArticleTag at = new ArticleTag();
                at.setArticleId(articleId);
                at.setTagId(tagId);
                articleTags.add(at);
            }
            articleTagMapper.batchInsert(articleTags);
        }
        // 返回更新后的完整文章对象，包含标签信息
        return this.getArticleById(articleId);
    }

    @Override
    @Transactional
    public void deleteArticleById(Integer id) {
        // 先删除标签关联数据
        articleTagMapper.deleteByArticleId(id);
        // 再删除文章主体
        baseMapper.deleteById(id);
    }
    @Override
    public Integer incrementViewCount(Integer id) {
        // 1. 先查询当前浏览量（避免更新时覆盖其他字段，且处理空值）
        Article article = baseMapper.selectById(id);
        if (article == null) {
            throw new RuntimeException("文章不存在");
        }

        // 2. 计算新浏览量（空值则初始化为1，否则+1）
        Integer currentViews = article.getViewsCount();
        int newViews = currentViews == null ? 1 : currentViews + 1;

        // 3. 仅更新浏览量字段（高效更新，只改需要的字段）
        UpdateWrapper<Article> updateWrapper = new UpdateWrapper<Article>()
                .eq("id", id)
                .set("views_count", newViews); // 数据库字段名（注意和实体类映射）
        baseMapper.update(null, updateWrapper);

        // 4. 返回最新浏览量
        return newViews;
    }
}