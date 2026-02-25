package com.example.blog.entity;

import com.baomidou.mybatisplus.annotation.*;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

import java.time.LocalDateTime;
import java.util.List;

@Data
@TableName("article")
public class Article {
    @TableId(type = IdType.AUTO)
    private Integer id;
    private String title;
    private String content;
    private Integer viewsCount;

    public Integer getViewsCount() {
        return viewsCount;
    }

    public void setViewsCount(Integer viewsCount) {
        this.viewsCount = viewsCount;
    }

    @TableField("cover_image")
    private String coverImage;

    // 发布时间：插入时自动填充
    @TableField(value = "publish_time", fill = FieldFill.INSERT)
    private LocalDateTime publishTime;

    // 更新时间：插入和更新时自动填充
    @TableField(value = "update_time", fill = FieldFill.INSERT_UPDATE)
    private LocalDateTime updateTime;

    @TableField("user_id")
    private Integer userId;

    // 用于接收前端传递的标签ID数组（修正JSON映射为tagIds）
    @TableField(exist = false) // 不映射数据库表字段
    @JsonProperty("tagIds")  // 与前端参数名保持一致
    private List<Integer> tagIds;

    // 新增：接收前端发送的targets字段，并同步到tagIds
    @TableField(exist = false)
    @JsonProperty("targets")
    private List<Integer> targets;

    // 用于返回给前端的标签对象列表
    @TableField(exist = false)
    private List<Tag> tags;

    // 自动将targets的值同步到tagIds
    public void setTargets(List<Integer> targets) {
        this.targets = targets;
        this.tagIds = targets; // 关键：将targets的值赋给tagIds，让业务逻辑能正确处理
    }

    // 以下为自动生成的getter/setter（保持不变）
    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public String getCoverImage() {
        return coverImage;
    }

    public void setCoverImage(String coverImage) {
        this.coverImage = coverImage;
    }

    public LocalDateTime getPublishTime() {
        return publishTime;
    }

    public void setPublishTime(LocalDateTime publishTime) {
        this.publishTime = publishTime;
    }

    public LocalDateTime getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(LocalDateTime updateTime) {
        this.updateTime = updateTime;
    }

    public Integer getUserId() {
        return userId;
    }

    public void setUserId(Integer userId) {
        this.userId = userId;
    }

    public List<Integer> getTagIds() {
        return tagIds;
    }

    public void setTagIds(List<Integer> tagIds) {
        this.tagIds = tagIds;
    }

    public List<Tag> getTags() {
        return tags;
    }

    public void setTags(List<Tag> tags) {
        this.tags = tags;
    }

    public List<Integer> getTargets() {
        return targets;
    }
}