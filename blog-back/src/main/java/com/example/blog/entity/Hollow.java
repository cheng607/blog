package com.example.blog.entity;

import com.baomidou.mybatisplus.annotation.*;
import lombok.Data;
import java.util.Date;


@Data
@TableName("hollow")
public class Hollow {
    @TableId(type= IdType.AUTO)
    private Integer id;
    private String word;
    @TableField(value = "date", fill= FieldFill.INSERT_UPDATE)
    private Date date;

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public String getWord() {
        return word;
    }

    public void setWord(String word) {
        this.word = word;
    }

    public Date getDate() {
        return date;
    }

    public void setDate(Date date) {
        this.date = date;
    }
}
