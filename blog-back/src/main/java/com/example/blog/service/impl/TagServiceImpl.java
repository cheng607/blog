package com.example.blog.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.example.blog.entity.Tag;
import com.example.blog.mapper.TagMapper;
import com.example.blog.service.TagService;
import com.example.blog.utils.ResultCode;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
public class TagServiceImpl extends ServiceImpl<TagMapper, Tag> implements TagService {

    @Override
    @Transactional
    public Integer saveTag(Tag tag) {
        // 校验标签名是否已存在（数据库有uk_tag_name唯一索引）
        QueryWrapper<Tag> queryWrapper = new QueryWrapper<>();
        queryWrapper.eq("name", tag.getName());
        if (baseMapper.selectCount(queryWrapper) > 0) {
            throw new RuntimeException(ResultCode.DATA_ALREADY_EXISTED.message());
        }
        baseMapper.insert(tag);
        return tag.getId();
    }

    @Override
    public List<Tag> listTagsWithArticleCount() {
        return baseMapper.selectTagsWithArticleCount();
    }

    @Override
    public Tag getByIdWithArticleCount(Integer id) {
        return baseMapper.selectByIdWithArticleCount(id);
    }
}