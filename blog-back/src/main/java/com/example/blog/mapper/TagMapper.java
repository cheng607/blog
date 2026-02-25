package com.example.blog.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.example.blog.entity.Tag;
import org.apache.ibatis.annotations.Mapper;

import java.util.List;

@Mapper
public interface TagMapper extends BaseMapper<Tag> {
    // 继承BaseMapper，已包含CRUD基础方法
    List<Tag> selectTagsWithArticleCount();
    Tag selectByIdWithArticleCount(Integer id);
}