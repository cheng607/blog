package com.example.blog.controller;

import com.example.blog.entity.User;
import com.example.blog.service.UserService;
import com.example.blog.utils.Result;
import com.example.blog.utils.ResultCode;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletRequest;
import java.util.List;

@RestController
@RequestMapping("/user")
public class UserController {

    @Autowired
    private UserService userService;

    // 查询所有用户
    @GetMapping
    public Result listUsers() {
        List<User> users = userService.findAll();
        return Result.success(users);
    }

    // 按id查询用户
    @GetMapping("/{id}")
    public Result getUserById(@PathVariable Integer id) { // id改为Integer
        Result r = new Result();
        if (id == null) {
            r.setResultCode(ResultCode.PARAM_IS_BLANK);
            return r;
        }
        User user = userService.getUserById(id);
        r.setResultCode(ResultCode.SUCCESS);
        r.setData(user);
        return r;
    }

    // 创建用户
    @PostMapping("/create")
    public Result saveUser(@Validated @RequestBody User user) {
        Integer userId = userService.saveUser(user);
        Result r = Result.success();
        r.simple().put("userId", userId);
        r.setData(user);
        return r;
    }

    // 更新用户
    @PutMapping("/update")
    public Result updateUser(@RequestBody User user) {
        Result r = new Result();
        if (user.getId() == null) {
            r.setResultCode(ResultCode.USER_NOT_EXIST);
            return r;
        }
        Integer userId = userService.updateUser(user);
        r.setResultCode(ResultCode.SUCCESS);
        r.simple().put("userId", userId);
        r.setData(user);
        return r;
    }

    // 按id删除用户
    @DeleteMapping("/delete/{id}")
    public Result deleteUserById(@PathVariable Integer id) {
        Result r = new Result();
        if (id == null) {
            r.setResultCode(ResultCode.PARAM_IS_BLANK);
            return r;
        }
        userService.deleteUserById(id);
        r.setResultCode(ResultCode.SUCCESS);
        return r;
    }
}