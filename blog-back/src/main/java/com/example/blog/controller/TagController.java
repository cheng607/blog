package com.example.blog.controller;

import com.example.blog.entity.Tag;
import com.example.blog.service.TagService;
import com.example.blog.utils.Result;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/tag")
public class TagController {

    @Autowired
    private TagService tagService;

    // 创建标签
    @PostMapping
    public Result saveTag(@RequestBody Tag tag) {
        Integer tagId = tagService.saveTag(tag);
        Result result = Result.success();
        result.simple().put("tagId", tagId);
        result.setData(tag);
        return result;
    }

    // 查询所有标签
    @GetMapping
    public Result getAllTags() {
        List<Tag> tags = tagService.listTagsWithArticleCount();
        return Result.success(tags);
    }

    // 根据ID查询标签
    @GetMapping("/{id}")
    public Result getTagById(@PathVariable Integer id) {
        Tag tag = tagService.getByIdWithArticleCount(id);
        return Result.success(tag);
    }

    // 更新标签
    @PutMapping
    public Result updateTag(@RequestBody Tag tag) {
        tagService.updateById(tag);
        return Result.success(tag);
    }

    // 删除标签
    @DeleteMapping("/{id}")
    public Result deleteTagById(@PathVariable Integer id) {
        tagService.removeById(id);
        return Result.success();
    }
}