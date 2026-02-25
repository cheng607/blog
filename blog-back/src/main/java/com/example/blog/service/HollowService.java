package com.example.blog.service;

import com.baomidou.mybatisplus.extension.service.IService;
import com.example.blog.entity.Hollow;

import java.util.List;

public interface HollowService extends IService<Hollow> {
    // 获取树洞
    List<Hollow> getHollow(Integer userId);
    // 保存树洞
    Integer saveHollow(Hollow hollow);
}
