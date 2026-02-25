package com.example.blog.controller;

import com.example.blog.utils.ImageUploadUtil;
import com.example.blog.utils.Result;
import com.example.blog.utils.ResultCode;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;

@RestController
@RequestMapping("/upload")
public class UploadController {

    @Autowired
    private ImageUploadUtil imageUploadUtil;

    /**
     * 图片上传接口
     * 前端请求方式：multipart/form-data，参数名"file"
     */
    @PostMapping("/image")
    public Result uploadImage(@RequestParam("file") MultipartFile file) {
        try {
            String imageUrl = imageUploadUtil.uploadImage(file);
            // 修正：先获取Result对象，再操作simple map，最后返回Result
            Result result = Result.success();
            result.simple().put("imageUrl", imageUrl);
            return result;
        } catch (IOException e) {
            // 修正：使用ResultCode.UPLOAD_ERROR作为错误码，错误信息放入data
            return Result.error(ResultCode.UPLOAD_ERROR, "图片上传失败：" + e.getMessage());
        } catch (RuntimeException e) {
            // 修正：使用ResultCode.UPLOAD_ERROR作为错误码，错误信息放入data
            return Result.error(ResultCode.UPLOAD_ERROR, e.getMessage());
        }
    }
}