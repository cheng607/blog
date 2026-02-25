package com.example.blog.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.example.blog.entity.User;
import com.example.blog.mapper.UserMapper;
import com.example.blog.service.UserService;
import com.example.blog.utils.ResultCode;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
public class UserServiceImpl implements UserService {

    @Autowired
    private UserMapper userMapper;

    @Override
    @Transactional
    public Integer saveUser(User user) {
        // 校验用户名是否已存在
        QueryWrapper<User> queryWrapper = new QueryWrapper<>();
        queryWrapper.eq("username", user.getUsername());
        if (userMapper.selectCount(queryWrapper) > 0) {
            throw new RuntimeException(ResultCode.USER_HAS_EXISTED.message());
        }
        userMapper.insert(user);
        return user.getId();
    }

    @Override
    public List<User> findAll() {
        return userMapper.selectList(null); // 查询所有用户
    }

    @Override
    public User getUserByUsername(String username) {
        // 按用户名查询（使用条件构造器）
        QueryWrapper<User> queryWrapper = new QueryWrapper<>();
        queryWrapper.eq("username", username);
        return userMapper.selectOne(queryWrapper);
    }

    @Override
    public User getUserById(Integer id) {
        return userMapper.selectById(id); // 按id查询
    }

    @Override
    public Integer updateUser(User user) {
        userMapper.updateById(user); // 更新用户
        return user.getId();
    }

    @Override
    public void deleteUserById(Integer id) {
        userMapper.deleteById(id); // 按id删除
    }
}