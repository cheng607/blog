package com.example.blog.service;

import com.example.blog.entity.User;
import java.util.List;

public interface UserService {
    // 新增用户
    Integer saveUser(User user);

    // 查询所有用户
    List<User> findAll();

    // 按用户名查询（原按account查询，改为按username）
    User getUserByUsername(String username);

    // 按id查询
    User getUserById(Integer id);

    // 更新用户
    Integer updateUser(User user);

    // 按id删除
    void deleteUserById(Integer id);
}