package com.example.blog.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.Data;

import java.time.LocalDateTime;

@Data
@TableName("tag") // 对应数据库的 tag 表
public class Tag {
    public LocalDateTime getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(LocalDateTime createdAt) {
        this.createdAt = createdAt;
    }

    public Integer getCount() {
        return count;
    }

    public void setCount(Integer count) {
        this.count = count;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    @TableId(type = IdType.AUTO) // 主键自增
    private Integer id;

    private String name; // 标签名（数据库有该列）

    private String description; // 标签描述（数据库有该列）

    // 关键修改：添加 @TableField(exist = false)，标记为非数据库字段
    @TableField(exist = false)
    private Integer count; // 仅用于前端展示“标签关联文章数量”，数据库无该列

    @TableField("created_at") // 对应数据库的 created_at 列
    private LocalDateTime createdAt;
}