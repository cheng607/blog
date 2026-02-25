package com.example.blog.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.example.blog.entity.Article;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Select;

import java.util.List;

@Mapper
public interface ArticleMapper extends BaseMapper<Article> {
    // 查询所有文章（包含标签）
    @Select("SELECT a.* FROM article a") // 仅保留这一行
    List<Article> selectAll();

    // 删除原@Select注解，仅保留方法定义，通过XML映射
    Article selectByIdWithTags(Integer id);

    List<Article> selectAllWithTags();

    // 根据标签ID查询文章
    List<Article> selectByTagId(Integer tagId);
}