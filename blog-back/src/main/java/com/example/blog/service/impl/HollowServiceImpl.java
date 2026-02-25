package com.example.blog.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.example.blog.entity.Hollow;
import com.example.blog.mapper.HollowMapper;
import com.example.blog.service.HollowService;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class HollowServiceImpl extends ServiceImpl<HollowMapper, Hollow> implements HollowService {
    @Override
    public List<Hollow> getHollow(Integer userId) {return baseMapper.selectAll();}

    @Override
    public  Integer saveHollow(Hollow hollow) {
        baseMapper.insert(hollow);
        return hollow.getId();
    }
}
