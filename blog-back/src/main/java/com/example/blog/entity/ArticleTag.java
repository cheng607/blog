package com.example.blog.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.Data;

@Data
@TableName("article_tag")
public class ArticleTag {
    @TableId(type = IdType.AUTO)
    private Integer id; // 关联ID（自增）

    private Integer articleId; // 文章ID

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public Integer getArticleId() {
        return articleId;
    }

    public Integer getTagId() {
        return tagId;
    }

    private Integer tagId; // 标签ID
    public void setArticleId(Integer articleId) {
        this.articleId = articleId;
    }
    public void setTagId(Integer tagId) {
        this.tagId = tagId;
    }
}