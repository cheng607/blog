// TagService接口
package com.example.blog.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.example.blog.entity.Tag;

import java.util.List;

public interface TagService extends IService<Tag> {
    // 新增标签（校验名称唯一）
    Integer saveTag(Tag tag);
    Tag getByIdWithArticleCount(Integer id);
    List<Tag> listTagsWithArticleCount();
}