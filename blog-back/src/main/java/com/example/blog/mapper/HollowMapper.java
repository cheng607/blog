package com.example.blog.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.example.blog.entity.Hollow;
import org.apache.ibatis.annotations.Insert;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import org.apache.ibatis.annotations.Select;

import java.util.List;

@Mapper
public interface HollowMapper extends BaseMapper<Hollow> {
    @Select("select * from hollow")
    List<Hollow> selectAll();

    @Insert("insert into hollow values({#id},{#word},{#date})")
    Integer save(@Param("hollow") Hollow hollow);

}
