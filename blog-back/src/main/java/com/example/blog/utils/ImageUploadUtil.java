package com.example.blog.utils;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.web.multipart.MultipartFile;
import java.io.File;
import java.io.IOException;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.UUID;

@Component
public class ImageUploadUtil {

    @Value("${upload.image.path}")
    private String uploadPath;

    @Value("${upload.image.url-prefix}")
    private String urlPrefix;

    /**
     * 处理图片上传
     * @param file 前端上传的文件
     * @return 图片访问URL
     */
    public String uploadImage(MultipartFile file) throws IOException {
        // 1. 校验文件是否为空
        if (file.isEmpty()) {
            throw new RuntimeException("上传文件不能为空");
        }

        // 2. 校验文件类型（仅允许jpg、png、jpeg）
        String originalFilename = file.getOriginalFilename();
        String suffix = originalFilename.substring(originalFilename.lastIndexOf("."));
        if (!".jpg".equalsIgnoreCase(suffix) && !".png".equalsIgnoreCase(suffix) && !".jpeg".equalsIgnoreCase(suffix)) {
            throw new RuntimeException("仅支持jpg、png、jpeg格式图片");
        }

        // 3. 生成唯一文件名（UUID+日期目录，避免重名）
        String uuid = UUID.randomUUID().toString().replaceAll("-", "");
        String dateDir = LocalDate.now().format(DateTimeFormatter.ofPattern("yyyyMMdd"));
        String fileName = uuid + suffix;

        // 4. 创建存储目录（如/Users/xxx/upload/images/20240520/）
        File saveDir = new File(uploadPath + "images/" + dateDir);
        if (!saveDir.exists()) {
            saveDir.mkdirs(); // 递归创建目录
        }

        // 5. 保存文件到本地
        File destFile = new File(saveDir, fileName);
        file.transferTo(destFile);

        // 6. 返回可访问的URL（如/images/20240520/xxx.jpg）
        return urlPrefix + dateDir + "/" + fileName;
    }
}