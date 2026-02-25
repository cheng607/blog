package com.example.blog.controller;

import com.example.blog.entity.Article;
import com.example.blog.service.ArticleService;
import com.example.blog.utils.Result;
import com.example.blog.utils.ResultCode;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/article")
public class ArticleController {

    @Autowired
    private ArticleService articleService;

    // 创建文章（含标签）
    @PostMapping
    public Result saveArticle(@RequestBody Article article) {
        Integer articleId = articleService.saveArticle(article);
        // 保存成功后，查询完整的文章信息（包含标签）返回给前端
        Article savedArticle = articleService.getArticleById(articleId);
        return Result.success(savedArticle);
    }

    // 查询所有文章
    @GetMapping("/list")
    public Result getArticleList() {
        List<Article> articles = articleService.getArticleList();
        return Result.success(articles);
    }

    // 根据ID查询文章（含标签）
    @GetMapping("/{id}")
    public Result getArticleById(@PathVariable Integer id) {
        // 参数校验
        if (id == null || id <= 0) {
            return Result.error(ResultCode.PARAM_IS_INVALID);
        }
        Article article = articleService.getArticleById(id);
        if (article == null) {
            return Result.error(ResultCode.ARTICLE_NOT_EXIST);
        }
        return Result.success(article);
    }

    // 根据用户ID查询文章
    @GetMapping("/user/{userId}")
    public Result getArticleByUserId(@PathVariable Integer userId) {
        // 参数校验
        if (userId == null || userId <= 0) {
            return Result.error(ResultCode.PARAM_IS_INVALID);
        }
        List<Article> articles = articleService.getArticleByUserId(userId);
        return Result.success(articles);
    }

    // 根据标签ID查询文章
    @GetMapping("/tag/{tagId}")
    public Result getArticleByTagId(@PathVariable Integer tagId) {
        // 参数校验
        if (tagId == null || tagId <= 0) {
            return Result.error(ResultCode.PARAM_IS_INVALID);
        }
        List<Article> articles = articleService.getArticleByTagId(tagId);
        return Result.success(articles);
    }

    // 更新文章（含标签）
    @PutMapping("/update")
    public Result updateArticle(@RequestBody Article article) {
        // 参数校验
        if (article == null || article.getId() == null) {
            return Result.error(ResultCode.PARAM_IS_INVALID);
        }
        // 调用修改后的服务方法，返回完整的文章对象
        Article updatedArticle = articleService.updateArticle(article);
        return Result.success(updatedArticle);
    }

    @DeleteMapping("/delete/{id}")
    public Result deleteArticle(@PathVariable("id") Integer id) {
        // 参数校验
        if (id == null || id <= 0) {
            return Result.error(ResultCode.PARAM_IS_INVALID);
        }
        // 检查文章是否存在
        Article article = articleService.getArticleById(id);
        if (article == null) {
            return Result.error(ResultCode.ARTICLE_NOT_EXIST);
        }
        articleService.deleteArticleById(id);
        return Result.success();
    }

    @PostMapping("/save")
    public Result saveArticle2(@RequestBody Article article) {  // 方法名修改避免冲突
        // 参数校验
        if (article == null || article.getUserId() == null) {
            return Result.error(ResultCode.PARAM_IS_INVALID);
        }
        Integer articleId = articleService.saveArticle(article);
        // 保存成功后，查询完整的文章信息（包含标签）返回给前端
        Article savedArticle = articleService.getArticleById(articleId);
        return Result.success(savedArticle);
    }
    @PostMapping("/view/{id}")
    public Result viewArticle(@PathVariable Integer id) {
        if (id == null || id <= 0) {
            return Result.error(ResultCode.PARAM_IS_INVALID);
        }
        // 调用Service更新浏览量，返回最新值
        Integer newViewCount = articleService.incrementViewCount(id);
        return Result.success(newViewCount);
    }
}