package com.example.blog.controller;

import com.example.blog.entity.Hollow;
import com.example.blog.mapper.HollowMapper;
import com.example.blog.utils.Result;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/hollow")
public class HollowController {
    @Autowired
    HollowMapper hollowMapper;

    @GetMapping
    public Result getHollowList(){
        List<Hollow> hollowList = hollowMapper.selectAll();
        return Result.success(hollowList);
    }

    @PostMapping
    public  Result saveHollow(@RequestBody Hollow hollow){
        Integer hollowId = hollowMapper.insert(hollow);
        Result result = new Result();
        result.simple().put("hollowId",hollowId);
        result.setData(hollow);
        return Result.success(hollow);
    }
}
